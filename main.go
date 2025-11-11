package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"triplea-backend-assignment/config"
	"triplea-backend-assignment/database"
	"triplea-backend-assignment/handlers"
	"triplea-backend-assignment/middleware"
	"triplea-backend-assignment/repository"
	"triplea-backend-assignment/service"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	accountRepo := repository.NewAccountRepository()
	transactionRepo := repository.NewTransactionRepository()

	accountService := service.NewAccountService(accountRepo)
	transactionService := service.NewTransactionService(transactionRepo, accountRepo)

	accountHandler := handlers.NewAccountHandler(accountService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware)

	router.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
	router.HandleFunc("/accounts/{account_id}", accountHandler.GetAccount).Methods("GET")
	router.HandleFunc("/transactions", transactionHandler.CreateTransaction).Methods("POST")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	serverAddr := cfg.GetServerAddress()
	log.Printf("Server starting on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

