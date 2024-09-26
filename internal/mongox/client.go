package mongox

import (
	"context"
	"github.com/retail-ai-inc/beanq"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
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
