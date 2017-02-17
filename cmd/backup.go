package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	oauth2webflow "github.com/aaron7/go-oauth2webflow"
	"github.com/aaron7/spotifybackup/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup your Spotify saved tracks and playlists",
	Run:   backupFunc,
}

func backupFunc(cmd *cobra.Command, args []string) {
	// Setup oauth2 config
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		Scopes: []string{
			"user-library-read",
			"playlist-read-private",
			"playlist-read-collaborative",
		},
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
	savedTracks, err := utils.GetAllItems(client, "https://api.spotify.com/v1/me/tracks")
	if err != nil {
		log.Fatal(err)
	}

	// Get playlists
	playlists, err := utils.GetAllItems(client, "https://api.spotify.com/v1/me/playlists")
	if err != nil {
		log.Fatal(err)
	}

	// Fetch tracks for each playlist
	var playlistsWithTracks []playlistObject
	for _, playlist := range playlists {
		// Decode interface{} as playlistObject
		var playlistObject playlistObject
		err := mapstructure.Decode(playlist, &playlistObject)
		if err != nil {
			log.Fatal("there was an error decoding the playlist to a playlistObject")
		}

		// Fetch tracks for playlist
		tracksURL := playlistObject.Tracks.Href
		tracks, err := utils.GetAllItems(client, tracksURL)
		if err != nil {
			log.Fatal(err)
		}

		// Set FetchedTracks and append playlistObject to playlistsWithTracks
		playlistObject.FetchedTracks = tracks
		playlistsWithTracks = append(playlistsWithTracks, playlistObject)
	}

	// Save the savedTracks and playlists in a json file
	backup := backupFormat{
		BackupTime:  time.Now().Format("2006-01-02T15:04:05-0700"),
		SavedTracks: savedTracks,
		Playlists:   playlistsWithTracks,
	}
	backupJSON, err := json.Marshal(backup)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("spotify_backup.json", backupJSON, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

type playlistWithoutTracksObject struct {
	Tracks struct {
		href string
	} `json:"tracks"`
}

type playlistObject struct {
	Collaborative bool          `json:"collaborative"`
	ExternalUrls  interface{}   `json:"external_urls"`
	Href          string        `json:"href"`
	ID            string        `json:"id"`
	Images        []interface{} `json:"images"`
	Name          string        `json:"name"`
	Owner         interface{}   `json:"owner"`
	Public        bool          `json:"public"`
	SnapshotID    string        `json:"snapshot_id"`
	Tracks        struct {
		Href string `json:"href"`
	} `json:"tracks"`
	FetchedTracks []interface{} `json:"fetchedTracks"`
	Type          string        `json:"type"`
	URI           string        `json:"uri"`
}

type backupFormat struct {
	BackupTime  string           `json:"backupTime"`
	SavedTracks []interface{}    `json:"savedTracks"`
	Playlists   []playlistObject `json:"playlists"`
}

func init() {
	RootCmd.AddCommand(backupCmd)
}
