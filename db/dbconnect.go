package db

import (
	"context"
	"fmt"
	"log"
	"lotteryapi/domain"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func init() {
	uri := os.Getenv("MONGO_URI")
	dbname := os.Getenv("MONGO_DB_NAME")
	colname := os.Getenv("MONGO_COL_NAME")

	clientOption := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Mongodb connected successfully")
	collection = client.Database(dbname).Collection(colname)

	fmt.Println("collection reference is ready", collection)
}

func getAllResults() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.M

	for cur.Next(context.Background()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}
	defer cur.Close(context.Background())

	return results
}

func getLatestResult() primitive.M {
	var myresult bson.M
	opts := options.FindOne().SetSort(bson.D{{"lotterydate", -1}})
	result := collection.FindOne(context.TODO(), bson.D{}, opts)
	err := result.Decode(&myresult)
	if err != nil {
		log.Fatal(err)
	}
	return myresult
}
func getByLotteryName(lotteryName string) primitive.M {
	var myresult bson.M
	filter := bson.M{"lotteryname": lotteryName}
	result := collection.FindOne(context.Background(), filter)
	err := result.Decode(&myresult)
	if err != nil {
		log.Fatal(err)
	}
	return myresult
}
func insertOneResult(result domain.GetLotteryResultRespose, collection *mongo.Collection) {
	_, err := collection.InsertOne(context.Background(), result)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted 1 result in db with id:")
}

func GetMyAllResults() []domain.GetLotteryResultRespose {
	allresults := getAllResults()
	var results []domain.GetLotteryResultRespose
	for _, value := range allresults {
		result := domain.GetLotteryResultRespose{}
		resultJson, err := bson.Marshal(value)
		if err != nil {
			log.Fatal(err)
		}
		bson.Unmarshal(resultJson, &result)
		results = append(results, result)

	}
	return results
}
func GetLatestResult() domain.GetLotteryResultRespose {
	result := domain.GetLotteryResultRespose{}
	value := getLatestResult()
	resultJson, err := bson.Marshal(value)
	if err != nil {
		log.Fatal(err)
	}
	bson.Unmarshal(resultJson, &result)
	return result
}
func GetByLotteryName(lotteryName string) domain.GetLotteryResultRespose {
	result := domain.GetLotteryResultRespose{}
	value := getByLotteryName(lotteryName)
	resultJson, err := bson.Marshal(value)
	if err != nil {
		log.Fatal(err)
	}
	bson.Unmarshal(resultJson, &result)
	return result
}
func CreateResult(result domain.GetLotteryResultRespose) {
	insertOneResult(result, collection)
}
