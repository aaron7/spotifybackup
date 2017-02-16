package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aaron7/go-oauth2webflow"
	"golang.org/x/oauth2"
)

func main() {
	// Setup oauth2 config
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		Scopes:       []string{"user-library-read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
	}

	// Get token via the OAuth2 Authorization Flow
	token, err := oauth2webflow.BrowserAuthCodeFlow(ctx, conf)
	if err != nil {
		log.Fatal(err)
	}
	client := conf.Client(ctx, token)

	// Get saved tracks
	savedTracks, err := getAllItems(client, "https://api.spotify.com/v1/me/tracks")
	if err != nil {
		log.Fatal(err)
	}

	savedTracksJSON, err := json.Marshal(savedTracks)
	if err != nil {
		log.Fatal(err)
	}

	// Save the savedTracks as json
	err = ioutil.WriteFile("saved_tracks.json", savedTracksJSON, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getAllItems(client *http.Client, url string) ([]interface{}, error) {
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

		// Add items to list
		itemsJSON = append(itemsJSON, pagingObject.Items...)
		nextURL = pagingObject.Next
	}

	return itemsJSON, nil
}

type pagingObject struct {
	Next  string        `json:"next"`
	Items []interface{} `json:"items"`
}
