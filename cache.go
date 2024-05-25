package main

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CacheAd struct {
	ad       Ad
	gender   map[string]bool
	country  map[string]bool
	platform map[string]bool
}

type Cache struct {
	genderIndex   map[string](map[*CacheAd]bool)
	countryIndex  map[string](map[*CacheAd]bool)
	platformIndex map[string](map[*CacheAd]bool)
	ageIndex      []([]*CacheAd)
	ads           []*CacheAd
	lock          sync.RWMutex
}

// Continuously updates the cache from the collection of ads at the specified time interval.
// It takes a TTL duration as a parameter determines how often the cache should be updated.
// The updater function runs indefinitely.
func (cache *Cache) Updater(ttl time.Duration) {
	for {
		cache.UpdateFromCollection(ads)
		time.Sleep(ttl)
	}
}

// Updates the cache with active ads from given MongoDB collection.
// Return true if the update is successful false otherwise
func (cache *Cache) UpdateFromCollection(collection *mongo.Collection) bool {
	now := time.Now()

	filter := bson.M{
		"start_at": bson.M{"$lte": now},
		"end_at":   bson.M{"$gte": now},
	}
	opts := options.Find().SetSort(bson.D{{Key: "start_at", Value: 1}})

	cur, err := coll.Find(context.Background(), filter, opts)
	if err != nil {
		log.Println("Error while fetching ads from MongoDB: ", err)
		return false
	}

	var results []Ad
	for cur.Next(context.Background()) {
		var ad Ad
		if err := cur.Decode(&ad); err != nil {
			log.Println("Error while decoding ad: ", err)
			continue
		}
		results = append(results, ad)
	}

	if err := cur.Err(); err != nil {
		log.Println("Error while fetching ads from MongoDB: ", err)
	}
	cur.Close(context.Background())

	cache.Update(results)

	return true
}

// update updates the cache with the given ads.
// It populates the genderIndex, countryIndex, platformIndex, ageIndex, and ads fields of the Cache.
// The ads parameter is a slice of Ad structs representing the ads to be added to the cache.
// Each ad is processed to update the corresponding indexes and maps in the cache.
// The function is thread-safe and uses a lock to ensure concurrent access to the cache is synchronized.
func (cache *Cache) Update(ads []Ad) {
	log.Println("Updating cache")

	cachedAds := make([]*CacheAd, 0, len(ads))
	genderIndex := make(map[string](map[*CachedAd]bool))
	countryIndex := make(map[string](map[*CachedAd]bool))
	platformIndex := make(map[string](map[*CachedAd]bool))
	ageIndex := make([]([]*CachedAd), 101)

	for _, ad := range ads {
		CachedAd := &CachedAd{
			ad:       ad,
			gender:   make(map[string]bool),
			country:  make(map[string]bool),
			platform: make(map[string]bool),
		}

		addToIndex(genderIndex, ad.Conditions.Gender, cachedAd, cachedAd.gender)
		addToIndex(countryIndex, ad.Conditions.Country, cachedAd, cachedAd.country)
		addToIndex(platformIndex, ad.Conditions.Platform, cachedAd, cachedAd.platform)
		addAgeToIndex(ageIndex, ad.Conditions.AgeStart, ad.Conditions.AgeEnd, cachedAd)

		cachedAds = append(cachedAds, cachedAd)
	}

	cache.lock.Lock()
	defer cache.lock.Unlock()
	
	cache.genderIndex = genderIndex
	cache.countryIndex = countryIndex
	cache.platformIndex = platformIndex
	cache.ageIndex = ageIndex
	cache.ads = cachedAds

	log.Println("Cache updated")
}

func addToIndex[T any](index map[T]map[*CachedAd]bool, items *[]T, cachedAd *CachedAd, cachedMap map[T]bool) {
	if items == nil {
		return
	}

	for _, item := range *items {
		if index[item] == nil {
			index[item] = make(map[*CachedAd]bool)
		}
		index[item][cachedAd] = true
		cachedMap[item] = true
	}
}

func addAgeToIndex(index map[int][]*CachedAd, ageStart, ageEnd *int, cachedAd *CachedAd) {
	if ageStart != nil && ageEnd != nil {
		for age := *ageStart; age <= *ageEnd; age++ {
			if index[age] == nil {
				index[age] = make([]*CachedAd, 0)
			}
			index[age] = append(index[age], cachedAd)
		}
	}
}


func (cache *Cache) Filter(query AdQuery) []Ad{
	cache.lock.RLock()
	defer cache.lock.RUnlock()

	matchingAds := make([]*CachedAd, 0, len(cache.ads))
	if query.Age > 0{
		if cache.ageIndex[query.Age] != nil{
			matchingAds = append(matchingAds, cache.ageIndex[query.Age]...)
		}
	}else{
		matchingAds = append(matchingAds, cache.ads...)
	}

	if query.Country != "" {
		if cache.countryIndex[query.Country] != nil {
			for i, ad := range matchingAds {
				if !cache.countryIndex[query.Country][ad] {
					matchingAds[i] = nil
				}
			}
		} else {
			matchingAds = make([]*CachedAd, 0)
		}
	}

	if query.Platform != "" {
		if cache.platformIndex[query.Platform] != nil {
			for i, ad := range matchingAds {
				if !cache.platformIndex[query.Platform][ad] {
					matchingAds[i] = nil
				}
			}
		} else {
			matchingAds = make([]*CachedAd, 0)
		}
	}

	if query.Gender != "" {
		if cache.genderIndex[query.Gender] != nil {
			for i, ad := range matchingAds {
				if !cache.genderIndex[query.Gender][ad] {
					matchingAds[i] = nil
				}
			}
		} else {
			matchingAds = make([]*CachedAd, 0)
		}
	}

	var skiped int64 = 0
	results := []Ad{}
	for _, ad := range matchingAds {
		if ad == nil{
			continue
		}

		if skiped < query.Offset{
			skiped++
			continue
		}

		results = append(results, ad.ad)

		if int64(len(results)) >= query.Limit{
			break
		}
	}

	return results
}