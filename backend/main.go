package main

import (
    "log"
    "os"

    "github.com/simonc2123/WEB_TEST/backend/db"
    "github.com/joho/godotenv"
	"github.com/simonc2123/WEB_TEST/backend/services"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the database
	conn, ctx, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer conn.Close(ctx) // Close the connection when done

	// Data from the API
	stocks, err := services.FetchAllStockData()
	if err != nil {
		log.Fatal("Error fetching stock data:", err)
	}

	db.InsertData(conn, ctx, stocks) // Insert data into the database

	log.Println("Data inserted successfully")
	os.Exit(0)// Exit with success status code
}