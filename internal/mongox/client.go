package mongox

import (
	"context"
	"github.com/retail-ai-inc/beanq/v3"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
	"sync"
)

var (
	mongoOnce sync.Once
	config    beanq.BeanqConfig
	mongoX    *MongoX
)

type MongoX struct {
	database   *mongo.Database
	collection string
}

func NewMongo() *MongoX {
	mongoOnce.Do(func() {
		if err := viper.Unmarshal(&config); err != nil {
			log.Fatalf("init mongo failed:%+v \n", err)
		}

		history := config.History
		uri := ""
		if history.On {
			uri = strings.Join([]string{"mongodb://", history.Mongo.Host, history.Mongo.Port}, "")
		}
		if uri == "" {
			log.Fatalln("mongo uri is validation")
		}
		opts := options.Client().ApplyURI(uri).
			SetConnectTimeout(history.Mongo.ConnectTimeOut).
			SetMaxPoolSize(history.Mongo.MaxConnectionPoolSize).
			SetMaxConnIdleTime(history.Mongo.MaxConnectionLifeTime)
		if history.Mongo.UserName != "" && history.Mongo.Password != "" {
			auth := options.Credential{
				AuthSource: history.Mongo.Database,
				Username:   history.Mongo.UserName,
				Password:   history.Mongo.Password,
			}
			opts.SetAuth(auth)
		}

		ctx := context.Background()
		client, err := mongo.Connect(ctx, opts)
		if err != nil {
			log.Fatal(err)
		}
		if err := client.Ping(ctx, nil); err != nil {
			log.Fatal(err)
		}
		mongoX = &MongoX{
			database:   client.Database(history.Mongo.Database),
			collection: history.Mongo.Collection,
		}
	})
	return mongoX
}

func (t *MongoX) EventLogs(ctx context.Context, filter bson.M, page, pageSize int64) ([]bson.M, int64, error) {
	skip := (page - 1) * pageSize
	if skip < 0 {
		skip = 0
	}
	opts := options.Find()
	opts.SetSkip(skip)
	opts.SetLimit(pageSize)
	opts.SetSort(bson.D{{"addTime", 1}})

	cursor, err := t.database.Collection(t.collection).Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	var data []bson.M
	if err := cursor.All(ctx, &data); err != nil {
		return nil, 0, err
	}
	total, err := t.database.Collection(t.collection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func (t *MongoX) DetailEventLog(ctx context.Context, id string) (bson.M, error) {

	filter := bson.M{}
	if id != "" {
		filter["id"] = id
	}

	single := t.database.Collection(t.collection).FindOne(ctx, filter)
	if err := single.Err(); err != nil {
		return nil, err
	}
	var data bson.M
	if err := single.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func (t *MongoX) Delete(ctx context.Context, id string) (int64, error) {
	filter := bson.M{}
	if id != "" {
		nid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, err
		}
		filter["_id"] = nid
	}
	result, err := t.database.Collection(t.collection).DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (t *MongoX) Edit(ctx context.Context, id string, payload any) (int64, error) {
	filter := bson.M{}
	if id != "" {
		nid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, err
		}
		filter["_id"] = nid
	}
	update := bson.D{
		{"$set", bson.D{{"payload", payload}}},
	}
	result, err := t.database.Collection(t.collection).UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}
