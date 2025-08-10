package main

import (
	"context"
	"fmt"
	"investbot/pkg/repositories"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func main() {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://OrestisStefanou:nAvqKPRTCAdnoUV2@cluster0.7lutc2a.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0").SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	userContextRepo, _ := repositories.NewUserContextMongoRepo(client, "investor_chatbot", "user_context")
	// userContext := domain.UserContext{
	// 	UserID: "kostis_pou_ta_kazia",
	// 	UserProfile: map[string]any{
	// 		"name": "Kostis",
	// 		"age":  300,
	// 	},
	// 	UserPortfolio: []domain.UserPortfolioHolding{
	// 		domain.UserPortfolioHolding{
	// 			AssetClass:          domain.Stock,
	// 			Symbol:              "MSFT",
	// 			Quantity:            10,
	// 			PortfolioPercentage: 1,
	// 		},
	// 	},
	// }
	// err = userContextRepo.UpdateUserContext(userContext)
	// if err != nil {
	// 	panic(err)
	// }
	user, err := userContextRepo.GetUserContext("kostis_pou_ta_kaziaa")
	if err != nil {
		panic(err)
	}
	fmt.Printf("USER: %+v\n", user)
}
