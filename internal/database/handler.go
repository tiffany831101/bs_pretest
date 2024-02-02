package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
	db     *mongo.Database
}

var MongoDB *DB

type Task struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Name   string             `bson:"name,omitempty"`
	Status int                `bson:"status"`
}

const taskCollection = "tasks"

func NewDB(dbURL, dbName string) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbURL))
	if err != nil {
		log.Fatal("Error initializeing MongoDB: ", err)
		return
	}

	// Try connecting to mongodb
	if err = client.Ping(context.TODO(), nil); err != nil {
		log.Fatal("Error Unable to Connect to MongoDB: ", err)
	} else {
		log.Println("Successfully Connected To MongDB.")
	}

	db := client.Database(dbName)

	MongoDB = &DB{
		client: client,
		db:     db,
	}
}

func (db *DB) CloseConnection() {
	if err := db.client.Disconnect(context.TODO()); err != nil {
		log.Fatal("Error disconnect to mongodb.")
	} else {
		log.Println("Successfully disconnect to mongdb.")
	}
}

func (db *DB) InsertSingleTask(task Task) error {
	collection := db.db.Collection(taskCollection)

	_, err := collection.InsertOne(context.TODO(), task)

	if err != nil {
		log.Println("Error Insert Single Task: ", err)
		return err
	}

	return nil
}

func (db *DB) GetTaskByID(taskID string) (Task, error) {

	collection := db.db.Collection(taskCollection)

	objectId, err := primitive.ObjectIDFromHex(taskID)

	if err != nil {
		return Task{}, nil
	}

	filter := bson.M{"_id": objectId}

	result := collection.FindOne(context.TODO(), filter)

	var task Task
	err = result.Decode(&task)

	return task, err
}

func (db *DB) GetTasks() ([]Task, error) {
	collection := db.db.Collection(taskCollection)

	var results []Task

	cursor, err := collection.Find(context.TODO(), bson.M{})
	err = cursor.All(context.TODO(), &results)

	return results, err
}

func (db *DB) DeleteTaskByID(taskID string) (int64, error) {
	collection := db.db.Collection(taskCollection)
	idPrimitive, err := primitive.ObjectIDFromHex(taskID)

	if err != nil {
		return -1, err
	}

	deletedResult, err := collection.DeleteOne(context.TODO(), bson.M{"_id": idPrimitive})

	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	return deletedResult.DeletedCount, nil

}

func (db *DB) UpdateTaskID(taskID string, task Task) error {
	collection := db.db.Collection(taskCollection)

	id, _ := primitive.ObjectIDFromHex(taskID)

	filter := bson.M{"_id": id}

	update := bson.M{"$set": bson.M{"name": task.Name, "status": task.Status}}
	_, err := collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	return nil
}
