package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/simonc2123/WEB_TEST/backend/services"
)

func ConnectDB() (*pgx.Conn, context.Context, error) {
	// Connect to the database
	dsn := os.Getenv("DB_URL") // Url to connect to the database
	if dsn == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	return conn, ctx, nil // Return the connection, context, and nil error
}

func InsertData(conn *pgx.Conn, ctx context.Context, stocks []services.StockItem) {
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
}
