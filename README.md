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
6. run `spotifybackup info` to get information about the backup

### Example
The following took 49 seconds to backup 213 [playlist objects](https://developer.spotify.com/web-api/object-model/#playlist-object-full)
and 28,096 [full track objects](https://developer.spotify.com/web-api/object-model/#track-object-full):
```bash
> spotifybackup
(progress printed here)
> spotifybackup info
Last backup: 2017-02-17T23:45:59+0000
Saved tracks: 1908
Playlists: 213 with 26188 tracks
Unique tracks across savedTracks and playlists: 23138
```

### Backup format (JSON)

- `backupTime` - ISO 8601 e.g. `2006-01-02T15:04:05-0700`
- `savedTracks` - list of [saved track objects](https://developer.spotify.com/web-api/object-model/#saved-track-object)
(with [full track objects](https://developer.spotify.com/web-api/object-model/#track-object-full))
- `playlists` - list of [simplified playlist objects](https://developer.spotify.com/web-api/object-model/#playlist-object-simplified)
containing the extra field `fetchedTracks` of [playlist track objects](https://developer.spotify.com/web-api/object-model/#playlist-track-object)
(with [full track objects](https://developer.spotify.com/web-api/object-model/#track-object-full))

### Todo

- Tests

Note: created this project as part of learning go
