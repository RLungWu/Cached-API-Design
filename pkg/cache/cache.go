package cache

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/RLungWu/Dcard-Backend-v2/pkg/ad"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CachedAd struct {
	ad       ad.Ad
	gender   map[string]bool
	country  map[string]bool
	platform map[string]bool
}

type Cache struct {
	genderIndex   map[string](map[*CachedAd]bool)
	countryIndex  map[string](map[*CachedAd]bool)
	platformIndex map[string](map[*CachedAd]bool)
	ageIndex      []([]*CachedAd) // 0-100
	ads           []*CachedAd
	lock          sync.RWMutex
}

// updateFromCollection updates the cache with active ads from the given mongoDB collection.
// Return true if the update is successful, false otherwise
func (c *Cache) updateFromCollection(coll *mongo.Collection) bool {
	now := time.Now()

	filter := bson.M{
		"startAt": bson.M{"$lte": now},
		"endAt":   bson.M{"$gte": now},
	}
	opts := options.Find().SetSort(bson.D{{Key: "startAt", Value: 1}})

	cur, err := coll.Find(context.Background(), filter, opts)
	if err != nil {
		log.Println("Error fetching active ads from DB: ", err)
		return false
	}

	var results []ad.Ad
	for cur.Next(context.Background()) {
		var ad ad.Ad
		if err := cur.Decode(&ad); err != nil {
			log.Println("Error decoding ad from DB: ", err)
			return false
		}
		results = append(results, ad)
	}

	if err := cur.Err(); err != nil {
		log.Println("Cursor error: ", err)
	}
	cur.Close(context.Background())

	c.Update(results)

	return true
}

// Update updates the cache with the given ads.
// It populates the genderIndex, countryIndex, platformIndex, ageIndex and ads fields of the Cache.
// The ads parameter is a slice of Ad structs that are currently active.
// Each ad is processed to update the corresponding indexes and maps in the cache.
// The function is thread-safe and uses a lock to ensure concurrent access to the cache is synchronized.
func (c *Cache) Update(ads []ad.Ad) {
	log.Println("Updating Cache...")

	cacheAds := make([]*CachedAd, 0, len(ads))
	genderIndex := make(map[string](map[*CachedAd]bool))
	countryIndex := make(map[string](map[*CachedAd]bool))
	platformIndex := make(map[string](map[*CachedAd]bool))
	ageIndex := make([]([]*CachedAd), 101)

	for _, ad := range ads {
		cachedAd := &CachedAd{
			ad:       ad,
			gender:   make(map[string]bool),
			country:  make(map[string]bool),
			platform: make(map[string]bool),
		}

		if ad.Conditions.Gender != nil {
			for _, g := range *ad.Conditions.Gender {
				if genderIndex[g] == nil {
					genderIndex[g] = make(map[*CachedAd]bool)
				}
				genderIndex[g][cachedAd] = true
				cachedAd.gender[g] = true
			}
		}

		if ad.Conditions.Country != nil {
			for _, c := range *ad.Conditions.Country {
				if countryIndex[c] == nil {
					countryIndex[c] = make(map[*CachedAd]bool)
				}
				countryIndex[c][cachedAd] = true
				cachedAd.country[c] = true
			}
		}

		if ad.Conditions.Platform != nil {
			for _, p := range *ad.Conditions.Platform {
				if platformIndex[p] == nil {
					platformIndex[p] = make(map[*CachedAd]bool)
				}
				platformIndex[p][cachedAd] = true
				cachedAd.platform[p] = true
			}
		}

		if ad.Conditions.AgeStart != nil && ad.Conditions.AgeEnd != nil {
			for age := *ad.Conditions.AgeStart; age <= *ad.Conditions.AgeEnd; age++ {
				if ageIndex[age] == nil {
					ageIndex[age] = make([]*CachedAd, 0)
				}
				ageIndex[age] = append(ageIndex[age], cachedAd)
			}
		}

		cacheAds = append(cacheAds, cachedAd)
	}

	c.lock.Lock()
	c.genderIndex = genderIndex
	c.countryIndex = countryIndex
	c.platformIndex = platformIndex
	c.ageIndex = ageIndex
	c.ads = cacheAds
	c.lock.Unlock()

	log.Println("Cache updated.")
}

func (c *Cache) Updater(ttl time.Duration) {
	for {
		time.Sleep(ttl)
		c.updateFromCollection(nil)
	}
}

func (c *Cache) Filter(query ad.AdQuery) []ad.Ad {
	c.lock.RLock()
	defer c.lock.RUnlock()

	matchingAds := make([]*CachedAd, 0, len(c.ads))
	if query.Age > 0 {
		if c.ageIndex[query.Age] != nil {
			matchingAds = append(matchingAds, c.ageIndex[query.Age]...)
		}
	}else{
		matchingAds = append(matchingAds, c.ads...)
	}

	if query.Country != "" {
		if c.countryIndex[query.Country] != nil {
			for i, ad := range matchingAds {
				if !c.countryIndex[query.Country][ad] {
					matchingAds[i] = nil
				}
			}
		} else {
			matchingAds = make([]*CachedAd, 0)
		}
	}

	if query.Platform != "" {
		if c.platformIndex[query.Platform] != nil {
			for i, ad := range matchingAds {
				if !c.platformIndex[query.Platform][ad] {
					matchingAds[i] = nil
				}
			}
		} else {
			matchingAds = make([]*CachedAd, 0)
		}
	}

	if query.Gender != "" {
		if c.genderIndex[query.Gender] != nil {
			for i, ad := range matchingAds {
				if !c.genderIndex[query.Gender][ad] {
					matchingAds[i] = nil
				}
			}
		} else {
			matchingAds = make([]*CachedAd, 0)
		}
	}

	var skiped int64 = 0
	results := []ad.Ad{}
	for _, ad := range matchingAds {
		if ad == nil {
			continue
		}
		if skiped < query.Offset {
			skiped++
			continue
		}
		results = append(results, ad.ad)
		if int64(len(results)) >= query.Limit {
			break
		}
	}

	return results
}
