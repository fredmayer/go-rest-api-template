package mongodb

import (
	"context"
	"github.com/fredmayer/go-rest-api-template/pkg/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

const KeepAlivePollPeriod = 3

var (
	Client *mongo.Client
	once   bool

	mongoDsn    string
	mongoUser   string
	mongoPasswd string
)

func NewClient(ctx context.Context, dsn string, user string, password string) *mongo.Client {
	mongoDsn = dsn
	mongoUser = user
	mongoPasswd = password
	Client, _ = dial(ctx)

	go keepAlive(ctx)

	return Client
}

func dial(ctx context.Context) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoDsn).SetAuth(options.Credential{
		Username: mongoUser,
		Password: mongoPasswd,
	})
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		if !once {
			logging.Log().Panicf("Cann`t connect to mongo %v \n", err)
		} else {
			logging.Log().Errorf("Cann`t connect to mongo %v \n", err)
		}
		return nil, err
	}

	once = true
	return client, nil
}

func keepAlive(ctx context.Context) {
	for {
		time.Sleep(time.Second * KeepAlivePollPeriod)

		select {
		case <-ctx.Done():
			if err := Client.Disconnect(context.TODO()); err != nil {
				log.Fatal(err)
			}
			return
		default:
			ctxPing, _ := context.WithTimeout(ctx, time.Second*3)

			err := Client.Ping(ctxPing, readpref.Primary())
			if err == nil {
				continue
			}

			logging.Log().Error("Lost mongo connection. Restoring...")

			Client, err = dial(ctx)
			if err != nil {
				logging.Log().Error(err)
				continue
			}
			logging.Log().Info("Mongo reconected")
		}
	}
}
