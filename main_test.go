package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
)

func Test_getAllItems__single_page(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.spotify.com").
		Get("/v1/me/tracks").
		Reply(200).
		JSON(`{"items": [{"track": "1"}, {"track": "2"}]}`)

	client := &http.Client{}
	items, err := getAllItems(client, "https://api.spotify.com/v1/me/tracks")

	assert.Nil(t, err)
	assert.Equal(t, []interface{}(
		[]interface{}{
			map[string]interface{}{"track": "1"},
			map[string]interface{}{"track": "2"},
		},
	), items)
}

func Test_getAllItems__multiple_pages(t *testing.T) {
	defer gock.Off()

	// Mock first request
	gock.New("https://api.spotify.com").
		Get("/v1/me/tracks").
		Reply(200).
		JSON(`{
			"next": "https://api.spotify.com/v1/me/tracks?offset=2&limit=2",
			"items": [{"track": "1"}, {"track": "2"}]
		}`)

	// Mock second request (second page)
	gock.New("https://api.spotify.com").
		Get("/v1/me/tracks").
		MatchParams(map[string]string{"offset": "2", "limit": "2"}).
		Reply(200).
		JSON(`{"items": [{"track": "3"}, {"track": "4"}]}`)

	client := &http.Client{}
	items, err := getAllItems(client, "https://api.spotify.com/v1/me/tracks")

	assert.Nil(t, err)
	assert.Equal(t, []interface{}(
		[]interface{}{
			map[string]interface{}{"track": "1"},
			map[string]interface{}{"track": "2"},
			map[string]interface{}{"track": "3"},
			map[string]interface{}{"track": "4"},
		},
	), items)
}
