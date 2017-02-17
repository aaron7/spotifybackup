package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// GetAllItems fetches items from a url which returns Spotify paging objects
func GetAllItems(client *http.Client, url string) ([]interface{}, error) {
	var itemsJSON []interface{}
	nextURL := url

	for nextURL != "" {
		// Make request to url
		log.Printf("getting %v", nextURL)
		resp, err := client.Get(nextURL)
		if err != nil {
			return itemsJSON, err
		}
		defer resp.Body.Close()

		// Check status code of API call
		if resp.StatusCode != 200 {
			return itemsJSON, fmt.Errorf("%v returned status %v instead of 200", nextURL, resp.StatusCode)
		}

		// Decode the pagingObject
		var pagingObject pagingObject
		err = json.NewDecoder(resp.Body).Decode(&pagingObject)
		if err != nil {
			return itemsJSON, err
		}

		// Extend items and set next url
		itemsJSON = append(itemsJSON, pagingObject.Items...)
		nextURL = pagingObject.Next
	}
	return itemsJSON, nil
}

type pagingObject struct {
	Next  string        `json:"next"`
	Items []interface{} `json:"items"`
}
