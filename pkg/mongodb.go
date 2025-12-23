package channels

import (
	"context"

	database "github.com/uug-ai/hub-pipeline-notification/internal/database"
	message "github.com/uug-ai/hub-pipeline-notification/message"
	"github.com/uug-ai/models/pkg/models"
)

type Mongodb struct{}

func (mongodb Mongodb) SendNotification(msg message.Message) error {

	if msg.UserId != "" {
		client := database.New().Client
		ctx, cancel := context.WithTimeout(context.Background(), database.TIMEOUT)
		defer cancel()

		// Open users collection
		db := client.Database(database.DatabaseName)
		notificationsCollection := db.Collection("notifications")

		// Create event
		_, err := notificationsCollection.InsertOne(ctx, msg)
		return err
	}

	return nil
}

func (mongodb Mongodb) CreateMarker(marker models.Marker) error {
	client := database.New().Client
	ctx, cancel := context.WithTimeout(context.Background(), database.TIMEOUT)
	defer cancel()

	// Open markers collection
	db := client.Database(database.DatabaseName)
	markersCollection := db.Collection("markers")

	// Create marker
	_, err := markersCollection.InsertOne(ctx, marker)
	return err
}
