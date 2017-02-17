spotifybackup
=============
[![Build Status](https://travis-ci.org/aaron7/spotifybackup.svg?branch=master)](https://travis-ci.org/aaron7/spotifybackup)

spotifybackup is a fast and simple command line tool to backup your Spotify
saved tracks and playlists. It saves all track information.

This application uses [go-oauth2webflow](https://github.com/aaron7/go-oauth2webflow)
to authenticate with Spotify automatically through your system browser.

All information about your saved tracks and playlists will be saved to a file
called `spotify_backup.json`. This file can be large due to the amount of track
information stored in each playlist.

### Instructions

1. Create a [Spotify app](https://developer.spotify.com/my-applications/#!/applications) and add `http://localhost:5000` as a Redirect URI
2. `export SPOTIFY_CLIENT_ID=... SPOTIFY_CLIENT_SECRET=...` (can save this in .bashrc/.zshrc)
3. `go get -v github.com/aaron7/spotifybackup`
4. `go install github.com/aaron7/spotifybackup`
5. run `spotifybackup` if `$GOPATH/bin` is in your path


### Todo

- Tests

Note: created this project as part of learning go
