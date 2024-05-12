package cache

import (
	"context"
	"sync"
	"time"
	"log"

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

func (c *Cache) Update(ads []ad.Ad) {
}

func (c *Cache) Updater(ttl time.Duration) {
}

func (c *Cache) Filter(query ad.AdQuery) []ad.Ad {
	return nil
}