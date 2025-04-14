package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	
)

type StockItem struct {
	Ticker string `json:"ticker"`
	Company string `json:"company"`
	Brokerage string `json:"brokerage"`
	Action string `json:"action"`
	RatingFrom string `json:"rating_from"`
	RatingTo string `json:"rating_to"`
	TargetFrom string `json:"target_from"`
	TargetTo string `json:"target_to"`
	Time string `json:"time"`
}

type APIResponse struct {
	Items []StockItem `json:"items"`
	NextPage string `json:"next_page"`
}

func FetchAllStockData() ([]StockItem, error) {
	apiKey := os.Getenv("API_KEY")
	baseUrl := "https://8j5baasof2.execute-api.us-west-2.amazonaws.com/production/swechallenge/list"

	var allItems []StockItem
	nextPage:= ""

	for {

		url := baseUrl
		if nextPage != "" {
			url += "?next_page=" + nextPage
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println("Error creating request:", err)
			return nil, err
		}

		req.Header.Add("Authorization", "Bearer "+apiKey)
		req.Header.Add("Content-Type", "application/json")
		
		client := &http.Client{}
		httpResp, err := client.Do(req)
		if err != nil {
			log.Println("Error making request:", err)
			return nil, err
		}

		defer httpResp.Body.Close()

		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			log.Println("Error reading httpResponse body:", err)
			return nil, err
		}

		if httpResp.StatusCode != 200 {
			log.Println("Error: httpResponse status code", httpResp.StatusCode)
			return nil, fmt.Errorf("error: response status code %d", httpResp.StatusCode)
		}

		// Unmarshal the response body into the APIResponse struct
		var apiResp APIResponse
		err = json.Unmarshal(body, &apiResp)
		if err != nil {
			log.Println("Error unmarshalling response body:", err)
			return nil, err
		}

		allItems = append(allItems, apiResp.Items...)

		if apiResp.NextPage == "" {
			break
		}

		nextPage = apiResp.NextPage

	}
	fmt.Println("Total items fetched:", len(allItems))
	return allItems, nil
}