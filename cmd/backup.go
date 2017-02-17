package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	oauth2webflow "github.com/aaron7/go-oauth2webflow"
	"github.com/aaron7/spotifybackup/utils"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

// backupCmd represents the backup command
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
	savedTracks, err := utils.GetAllItems(client, "https://api.spotify.com/v1/me/tracks")
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

func init() {
	RootCmd.AddCommand(backupCmd)
}
