package pkg

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

type Mongodb struct {
	Client *mongo.Client
}

var TIMEOUT = 10 * time.Second
var _init_ctx sync.Once
var _instance *Mongodb

var DatabaseName = os.Getenv("MONGODB_DATABASE_CLOUD")

func New() *Mongodb {

	mongodbURI := os.Getenv("MONGODB_URI")
	host := os.Getenv("MONGODB_HOST")
	databaseCredentials := os.Getenv("MONGODB_DATABASE_CREDENTIALS")
	replicaset := os.Getenv("MONGODB_REPLICASET")
	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")
	authenticationMechanism := os.Getenv("MONGODB_AUTHENTICATION_MECHANISM")
	retryWrites := os.Getenv("MONGODB_RETRY_WRITES") != "false" // Default to true unless explicitly set to "false"

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	_init_ctx.Do(func() {
		_instance = new(Mongodb)

		// We can also apply the complete URI
		// e.g. "mongodb+srv://<username>:<password>@kerberos-hub.shhng.mongodb.net/?retryWrites=true&w=majority&appName=kerberos-hub"
		if mongodbURI != "" {
			serverAPI := options.ServerAPI(options.ServerAPIVersion1)
			opts := options.Client().
				ApplyURI(mongodbURI).
				SetServerAPIOptions(serverAPI).
				SetRetryWrites(retryWrites).
				SetMonitor(otelmongo.NewMonitor(otelmongo.WithCommandAttributeDisabled(false)))

			// Create a new client and connect to the server
			client, err := mongo.Connect(ctx, opts)
			if err != nil {
				fmt.Printf("Error setting up mongodb connection: %+v\n", err)
				os.Exit(1)
			}
			_instance.Client = client

		} else {

			// New MongoDB driver
			mongodbURI := fmt.Sprintf("mongodb://%s:%s@%s", username, password, host)
			if replicaset != "" {
				mongodbURI = fmt.Sprintf("%s/?replicaSet=%s", mongodbURI, replicaset)
			}
			if authenticationMechanism == "" {
				authenticationMechanism = "SCRAM-SHA-256"
			}
			client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbURI).SetRetryWrites(retryWrites).SetAuth(options.Credential{
				AuthMechanism: authenticationMechanism,
				AuthSource:    databaseCredentials,
				Username:      username,
				Password:      password,
			}))
			if err != nil {
				fmt.Printf("Error setting up mongodb connection: %+v\n", err)
				os.Exit(1)
			}
			_instance.Client = client
		}
	})

	return _instance
}

func (mongodb *Mongodb) SendNotification(msg Message) error {

	if msg.UserId != "" {
		client := New().Client
		ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
		defer cancel()

		// Open users collection
		db := client.Database(DatabaseName)
		notificationsCollection := db.Collection("notifications")

		// Create event
		_, err := notificationsCollection.InsertOne(ctx, msg)
		return err
	}

	return nil
}

func (mongodb *Mongodb) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	err := New().Client.Ping(ctx, nil)
	return err
}
