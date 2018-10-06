package libshogun

import (
	"net/http"
)

// ShogunClient repersents a client to interact with Shogun
type ShogunClient struct {
	HTTP       *http.Client
	DauthToken string
}

// Title repersents a title in Shogun
type Title struct {
	ID          int64
	Name        string
	BannerURL   string
	ReleaseDate string
	IsNew       bool
	Description string
	Genre       string
	Size        int64
	Screenshots []string
	Movies      []*Movie
	Publisher   *Publisher
	TitleID     string
	IconURL     string
}

// Movie repersents a movie
type Movie struct {
	URL       string
	Thumbnail string
}

// Publisher repersents a title's publisher
type Publisher struct {
	ID   int64
	Name string
}
