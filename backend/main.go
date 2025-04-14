package main

import (
    "context"
    "log"
    "os"

    "github.com/simonc2123/WEB_TRUORA_TEST/backend/services"
    "github.com/joho/godotenv"
    "github.com/jackc/pgx/v4"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Connect to the database
	dsn := os.Getenv("DB_URL")// Url to connect to the database
	if dsn == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer conn.Close(ctx)

	// Data from the API
	stocks, err := services.FetchAllStockData()
	if err != nil {
		log.Fatal("Error fetching stock data:", err)
	}

	for _, stock := range stocks {
		// Insert data into the database
		_, err := conn.Exec(ctx, `
		INSERT INTO stock_items (ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, time)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, stock.Ticker, stock.Company, stock.Brokerage, stock.Action, stock.RatingFrom, stock.RatingTo, stock.TargetFrom, stock.TargetTo, stock.Time)
		if err != nil {
			log.Println("Error inserting stock item:", err)
		}
	}
	log.Println("Data inserted successfully")
	os.Exit(0)// Exit with success status code
}