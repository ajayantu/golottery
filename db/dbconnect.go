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

func ConnectDB() *mongo.Collection {
	var collection *mongo.Collection
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
	return collection
}

func getAllResults(collection *mongo.Collection) []primitive.M {
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

func getLatestResult(collection *mongo.Collection) (primitive.M, error) {
	var myresult bson.M
	opts := options.FindOne().SetSort(bson.D{{"lotterydate", -1}})
	result := collection.FindOne(context.TODO(), bson.D{}, opts)
	err := result.Decode(&myresult)
	if err != nil {
		return nil, err
	}
	return myresult, nil
}
func getByLotteryName(collection *mongo.Collection, lotteryName string) (primitive.M, error) {
	var myresult bson.M
	filter := bson.M{"lotteryname": lotteryName}
	result := collection.FindOne(context.Background(), filter)
	err := result.Decode(&myresult)
	if err != nil {
		return nil, err
	}
	return myresult, nil
}
func insertOneResult(result domain.GetLotteryResultRespose, collection *mongo.Collection) {
	_, err := collection.InsertOne(context.Background(), result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted 1 result in db with id:")
}
func insertManyResult(result []domain.GetLotteryResultRespose, collection *mongo.Collection) {
	var result1 []any
	for _, item := range result {
		result1 = append(result1, item)
	}
	_, err := collection.InsertMany(context.Background(), result1)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted many result in db with id:")
}

func GetMyAllResults(collection *mongo.Collection) []domain.GetLotteryResultRespose {
	allresults := getAllResults(collection)
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
func GetLatestResult(collection *mongo.Collection) domain.GetLotteryResultRespose {
	result := domain.GetLotteryResultRespose{}
	value, err := getLatestResult(collection)
	if err != nil {
		return domain.GetLotteryResultRespose{}
	}
	resultJson, err := bson.Marshal(value)
	if err != nil {
		log.Fatal(err)
	}
	bson.Unmarshal(resultJson, &result)
	return result
}
func GetByLotteryName(collection *mongo.Collection, lotteryName string) domain.GetLotteryResultRespose {
	result := domain.GetLotteryResultRespose{}
	value, err := getByLotteryName(collection, lotteryName)
	if err != nil {
		return domain.GetLotteryResultRespose{}
	}
	resultJson, err := bson.Marshal(value)
	if err != nil {
		log.Fatal(err)
	}
	bson.Unmarshal(resultJson, &result)
	return result
}
func InsertOneResult(collection *mongo.Collection, result domain.GetLotteryResultRespose) {
	insertOneResult(result, collection)
}
func InsertManyResults(collection *mongo.Collection, results []domain.GetLotteryResultRespose) {
	fmt.Println("results to inser are", results)
	insertManyResult(results, collection)
}
