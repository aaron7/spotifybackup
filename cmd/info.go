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

	var backup backupFormat
	json.Unmarshal(file, &backup)
	fmt.Printf("Last backup: %v\n", backup.BackupTime)
	fmt.Printf("Number of saved tracks: %v\n", len(backup.SavedTracks))

	var totalPlaylistTracks int
	for _, playlist := range backup.Playlists {
		totalPlaylistTracks = totalPlaylistTracks + len(playlist.FetchedTracks)
	}
	fmt.Printf("Number of playlists: %v with %v tracks\n", len(backup.Playlists), totalPlaylistTracks)
}

func init() {
	RootCmd.AddCommand(infoCmd)
}
