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
	_savedTracks, err := utils.GetAllItems(client, "https://api.spotify.com/v1/me/tracks?limit=50")
	if err != nil {
		log.Fatal(err)
	}
	// Decode savedTracks []interface{} as []trackObject
	var savedTracks []trackObject
	err = mapstructure.Decode(_savedTracks, &savedTracks)
	if err != nil {
		log.Fatal("there was an error decoding _savedTracks to []trackObject")
	}

	// Get playlists
	_playlists, err := utils.GetAllItems(client, "https://api.spotify.com/v1/me/playlists")
	if err != nil {
		log.Fatal(err)
	}
	// Decode playlists []interface{} as []playlistObject
	var playlists []playlistObject
	err = mapstructure.Decode(_playlists, &playlists)
	if err != nil {
		log.Fatal("there was an error decoding _playlists to []playlistObject")
	}

	// Fetch tracks for each playlist
	var playlistsWithTracks []playlistObject
	for _, playlist := range playlists {
		// Fetch tracks for playlist
		tracksURL := playlist.Tracks.Href
		_tracks, err := utils.GetAllItems(client, tracksURL)
		if err != nil {
			log.Fatal(err)
		}
		// Decode tracks []interface{} as []trackObject
		var tracks []trackObject
		err = mapstructure.Decode(_tracks, &tracks)
		if err != nil {
			log.Fatal("there was an error decoding _tracks to []trackObject")
		}

		// Set FetchedTracks and append playlistObject to playlistsWithTracks
		playlist.FetchedTracks = tracks
		playlistsWithTracks = append(playlistsWithTracks, playlist)
	}

	// Encode backup
	backup := backupFormat{
		BackupTime:  time.Now().Format("2006-01-02T15:04:05-0700"),
		SavedTracks: savedTracks,
		Playlists:   playlistsWithTracks,
	}
	backupJSON, err := json.Marshal(backup)
	if err != nil {
		log.Fatal(err)
	}

	// Save the backup to a JSON file
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
	FetchedTracks []trackObject `json:"fetchedTracks"`
	Type          string        `json:"type"`
	URI           string        `json:"uri"`
}

type fullTrackObject struct {
	Album            interface{}   `json:"album"`
	Artists          []interface{} `json:"artists"`
	AvailableMarkets []string      `json:"available_markets"`
	DiscNumber       int           `json:"disc_number"`
	DurationMs       int           `json:"duration_ms"`
	Explicit         bool          `json:"explicit"`
	ExternalIds      interface{}   `json:"external_ids"`
	ExternalUrls     interface{}   `json:"external_urls"`
	Href             string        `json:"href"`
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	Popularity       int           `json:"popularity"`
	PreviewURL       string        `json:"preview_url"`
	TrackNumber      int           `json:"track_number"`
	Type             string        `json:"type"`
	URI              string        `json:"uri"`
}

type trackObject struct {
	AddedAt string          `json:"added_at"`
	AddedBy interface{}     `json:"added_by"`
	IsLocal bool            `json:"is_local"`
	Track   fullTrackObject `json:"track"`
}

type backupFormat struct {
	BackupTime  string           `json:"backupTime"`
	SavedTracks []trackObject    `json:"savedTracks"`
	Playlists   []playlistObject `json:"playlists"`
}

func init() {
	RootCmd.AddCommand(backupCmd)
}
