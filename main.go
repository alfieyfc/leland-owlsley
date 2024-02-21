package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alfieyfc/nwtp/configs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/bson"
)

var collection = configs.GetCollection("marketprice")

type MarketPrices struct {
	ItemId string
	ItemName string
	Price string
	Availability int
	LastUpdated string
	HighestBuyOrder float32
	Qty int
}

func main() {

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/hello", helloWorld)
	router.Get("/random", getRandomPrice)

	server := &http.Server{
		Addr: ":3000",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err!= nil {
		fmt.Println("failed to listen to server", err)
	}
}

func helloWorld(w http.ResponseWriter, r *http.Request){
	// // Uncomment to use /hello for seeding your database
	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()

	// byteValues, err := os.ReadFile("configs/db/seed.json")
	// if err != nil {
	// 	panic(err)
	// }

	// var docs []MarketPrices
	// err = json.Unmarshal(byteValues, &docs)
	// if err != nil {
	// 	panic(err)
	// }
	// for i := range docs {
	// 	doc := docs[i]
	// 	result, insertErr := collection.InsertOne(ctx, doc)
	// 	if insertErr != nil {
	// 		panic(insertErr)
	// 	} else {
	// 		fmt.Println(result)
	// 	}
	// }

	w.Write([]byte("Hello, world!"))
}

func getRandomPrice(w http.ResponseWriter, r *http.Request){
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	pipeline := bson.A{
		bson.M{"$sample": bson.M{"size": 1}},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			panic(err)
		}
		results = append(results, result)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
