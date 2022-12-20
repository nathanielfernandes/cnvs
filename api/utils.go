package main

import "regexp"

var spotifyUrlMatch = regexp.MustCompile(`https://?open.spotify.com/track/([a-zA-Z0-9]+)`)
var spotifyUriMatch = regexp.MustCompile(`spotify:track:([a-zA-Z0-9]+)`)
var spotifyIdMatch = regexp.MustCompile(`([a-zA-Z0-9]+)`)

func GetSpotifyUriFromUrl(url string) string {
	return spotifyUrlMatch.FindStringSubmatch(url)[1]
}

func GetSpotifyUriFromUri(uri string) string {
	return spotifyUriMatch.FindStringSubmatch(uri)[1]
}

func GetSpotifyUriFromId(id string) string {
	return spotifyIdMatch.FindStringSubmatch(id)[1]
}

func GetSpotifyIdFromAny(input string) string {
	if spotifyUrlMatch.MatchString(input) {
		return GetSpotifyUriFromUrl(input)
	} else if spotifyUriMatch.MatchString(input) {
		return GetSpotifyUriFromUri(input)
	} else if spotifyIdMatch.MatchString(input) {
		return GetSpotifyUriFromId(input)
	}
	return ""
}

func GetSpotifyUriFromAny(input string) string {
	if spotifyUrlMatch.MatchString(input) {
		return "spotify:track:" + GetSpotifyUriFromUrl(input)
	} else if spotifyUriMatch.MatchString(input) {
		return input
	} else if spotifyIdMatch.MatchString(input) {
		return "spotify:track:" + GetSpotifyUriFromId(input)
	}
	return ""
}
