package cache

import (
	"time"

	"github.com/RLungWu/Dcard-Backend-v2/pkg/ad"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// updateFromCollection updates the cache with active ads from the given mongoDB collection.
// Return true if the update is successful, false otherwise
func (c *ad.Cachee) updateFromCollection(coll *mongo.Collection) bool {
	now := time.Now()

	filter := bson.M{
		"startAt": bson.M{"$lte": now},
		"endAt":   bson.M{"$gte": now},
	}
	opts := options.Find().SetSort(bson.D{{Key: "startAt",Value: 1}})
}
