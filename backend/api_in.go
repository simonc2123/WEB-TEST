package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/simonc2123/WEB_TEST/backend/db"
	"github.com/simonc2123/WEB_TEST/backend/services"
)

// cleaning the data
func parseMonetaryValue(value string) (float64, error) {
	cleanedValue := strings.ReplaceAll(value, "$", "")
	cleanedValue = strings.ReplaceAll(cleanedValue, ",", "")
	parsedValue, err := strconv.ParseFloat(cleanedValue, 64)
	if err != nil {
		log.Println("Error parsing monetary value '%s': %v", value, err)
		return 0, err
	}
	return parsedValue, nil
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})
	// Endpoint to fetch all stock data
	r.GET("/api/stocks", func(c *gin.Context) {
		// Connect to the database
		conn, ctx, err := db.ConnectDB()
		if err != nil {
			log.Fatal("Error connecting to the database:", err)
		}
		defer conn.Close(ctx) // Close the connection when done

		// Fetch all stock data from the database
		rows, err := conn.Query(ctx, "SELECT ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, time FROM stock_items")
		if err != nil {
			log.Fatal("Error fetching stock data:", err)
		}
		defer rows.Close()

		var stocks []services.StockItem
		for rows.Next() {
			var stock services.StockItem
			err := rows.Scan(
				&stock.Ticker,
				&stock.Company,
				&stock.Brokerage,
				&stock.Action,
				&stock.RatingFrom,
				&stock.RatingTo,
				&stock.TargetFrom,
				&stock.TargetTo,
				&stock.Time,
			)
			if err != nil {
				log.Println("Error scanning row:", err)
				continue
			}

			//Data cleaning
			targetFrom, err := parseMonetaryValue(stock.TargetFrom)
			targetTo, err := parseMonetaryValue(stock.TargetTo)

			stock.TargetFrom = fmt.Sprintf("%.2f", targetFrom)
			stock.TargetTo = fmt.Sprintf("%.2f", targetTo)

			stocks = append(stocks, stock)
		}

		c.JSON(http.StatusOK, stocks) // Return the stock data as JSON

	})

	r.GET("/api/recommendations", func(c *gin.Context) {
		// Connect to the database
		conn, ctx, err := db.ConnectDB()
		if err != nil {
			log.Fatal("Error connecting to the database:", err)
		}
		defer conn.Close(ctx) // Close the connection when done

		// Fetch recommendations from the database
		rows, err := conn.Query(ctx, `
			SELECT ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, time
        	FROM stock_items
		`)
		if err != nil {
			log.Fatal("Error fetching recommendations:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recommendations"})
			return
		}

		defer rows.Close()

		var stocks []services.StockItem
		for rows.Next() {
			var stock services.StockItem
			err := rows.Scan(
				&stock.Ticker,
				&stock.Company,
				&stock.Brokerage,
				&stock.Action,
				&stock.RatingFrom,
				&stock.RatingTo,
				&stock.TargetFrom,
				&stock.TargetTo,
				&stock.Time,
			)
			if err != nil {
				log.Println("Error scanning row:", err)
				continue
			}

			//Data cleaning
			targetFrom, err := parseMonetaryValue(stock.TargetFrom)
			targetTo, err := parseMonetaryValue(stock.TargetTo)

			stock.TargetFrom = fmt.Sprintf("%.2f", targetFrom)
			stock.TargetTo = fmt.Sprintf("%.2f", targetTo)

			stocks = append(stocks, stock)
		}

		//Filter stocks based on the criteria
		var recommendations []services.StockItem
		for _, stock := range stocks {
			targetFrom, err := parseMonetaryValue(stock.TargetFrom)
			if err != nil {
				log.Println("Error parsing target_from:", err)
				continue
			}

			targetTo, err := parseMonetaryValue(stock.TargetTo)
			if err != nil {
				log.Println("Error parsing target_to:", err)
				continue
			}

			if (stock.Action == "target raised by" || stock.Action == "upgraded by") &&
				(stock.RatingTo == "Buy" || stock.RatingTo == "Outperform" || stock.RatingTo == "Overweight") &&
				targetTo > targetFrom {
				stock.TargetIncrease = targetTo - targetFrom
				recommendations = append(recommendations, stock)
			}
		}

		sort.Slice(recommendations, func(i, j int) bool {
			return recommendations[i].TargetIncrease > recommendations[j].TargetIncrease
		})

		if len(recommendations) > 5 {
			recommendations = recommendations[:5] // Limit to top 5 recommendations
		}

		c.JSON(http.StatusOK, recommendations) // Return the recomendations as JSON
	})
	r.Run(":8080") // Start the server on port 8080
	log.Println("Server started on :8080")
}
