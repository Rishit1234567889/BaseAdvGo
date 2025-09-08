package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Rishit1234567889/baseToAdvGo/config"
	"github.com/Rishit1234567889/baseToAdvGo/internal/handlers"
	"github.com/Rishit1234567889/baseToAdvGo/internal/routes"
	"github.com/Rishit1234567889/baseToAdvGo/internal/store"
	"github.com/redis/go-redis/v9"
)

func main() {

	configLd, err := config.LoadConfig() // 1.6 Load config
	if err != nil {
		log.Fatalf("Failed to load config %v", err)
	}
	fmt.Println("App is running i guess ")

	db := config.ConnectDB(configLd.DataBaseUrl) // 2.1 connect to db
	defer db.Close()

	//connect to redis 7.2
	rdb := config.ConnectRedis()
	defer func(rdb *redis.Client) {
		_ = rdb.Close()
	}(rdb)
	queries := store.New(db) // 3.6

	handler := handlers.NewHandlers(db, queries, rdb) // 1.7 () Create a new Handler

	mux := http.NewServeMux() // 1.0 setUp the HTTP server first

	routes.SetupRoutes(mux, handler) // 1.7 (D) setup Routes

	serverAddr := fmt.Sprintf(":%s", configLd.ServerPort) // 1.1server instance
	server := &http.Server{
		Addr:    serverAddr,
		Handler: mux,
	}
	fmt.Printf("Server is up and running on PORT %s\n", serverAddr)
	if err := server.ListenAndServe(); err != nil { // 1.2 Listen
		log.Fatalf("Server failed %v", err)
	}

}
