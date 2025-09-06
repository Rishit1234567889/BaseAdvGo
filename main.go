package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Rishit1234567889/baseToAdvGo/internal/handlers"
	"github.com/Rishit1234567889/baseToAdvGo/internal/routes"
	"github.com/Rishit1234567889/baseToAdvGo/serverconfig"
)

func main() {

	config, err := serverconfig.LoadConfig() // 1.6 Load config
	if err != nil {
		log.Fatalf("Failed to load config %v", err)
	}
	fmt.Println("App is running i guess ")
	// 1.7 () Create a new Handler
	handler := handlers.NewHandlers()

	mux := http.NewServeMux() // 1.0 setUp the HTTP server first

	routes.SetupRoutes(mux, handler) // 1.7 (D) setup Routes

	serverAddr := fmt.Sprintf(":%s", config.ServerPort) // 1.1server instance
	server := &http.Server{
		Addr:    serverAddr,
		Handler: nil,
	}
	fmt.Printf("Server is up and running on PORT %s\n", serverAddr)
	if err := server.ListenAndServe(); err != nil { // 1.2 Listen
		log.Fatalf("Server failed %v", err)
	}

}
