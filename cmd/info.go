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
	file, err := ioutil.ReadFile("./saved_tracks.json")
	if err != nil {
		log.Fatal(err)
	}

	var savedTracks []interface{}
	json.Unmarshal(file, &savedTracks)
	fmt.Printf("Number of saved tracks: %v\n", len(savedTracks))
}

func init() {
	RootCmd.AddCommand(infoCmd)
}
