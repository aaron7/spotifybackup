package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show information about your current backup",
	Run:   infoFunc,
}

func infoFunc(cmd *cobra.Command, args []string) {
	file, err := ioutil.ReadFile("./spotify_backup.json")
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal backup and print basic statistics
	var backup backupFormat
	json.Unmarshal(file, &backup)
	fmt.Printf("Last backup: %v\n", backup.BackupTime)
	fmt.Printf("Saved tracks: %v\n", len(backup.SavedTracks))

	// Count unqiue tracks in savedTracks
	uniqueTracks := make(map[string]bool)
	for _, savedTrack := range backup.SavedTracks {
		uniqueTracks[savedTrack.Track.ID] = true
	}

	// Count total playlist tracks
	var totalPlaylistTracks int
	for _, playlist := range backup.Playlists {
		totalPlaylistTracks = totalPlaylistTracks + len(playlist.FetchedTracks)
		// Count unique tracks
		for _, track := range playlist.FetchedTracks {
			uniqueTracks[track.Track.ID] = true
		}
	}
	fmt.Printf("Playlists: %v with %v tracks\n", len(backup.Playlists), totalPlaylistTracks)
	fmt.Printf("Unique tracks across savedTracks and playlists: %v\n", len(uniqueTracks))
}

func init() {
	RootCmd.AddCommand(infoCmd)
}
