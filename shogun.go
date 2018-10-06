// Package libshogun provides various utilities for working with the Nintendo Switch's title metadata server, Shogun
package libshogun

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"net/http"
)

// NewShogunClient creates a new ShogunClient
func NewShogunClient(shopnCert, shopnKey, dauthToken string) (client *ShogunClient, err error) {
	shopn, err := tls.LoadX509KeyPair(shopnCert, shopnKey)
	if err != nil {
		return nil, err
	}

	return &ShogunClient{
		&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Certificates:       []tls.Certificate{shopn},
					InsecureSkipVerify: true,
				},
			},
		},
		dauthToken,
	}, nil
}

// DoRequest sends a request to the specified URL with the proper authentication
func (c *ShogunClient) DoRequest(url string) (response []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-DeviceAuthorization", c.DauthToken)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return bytes, nil
}

// DoShogunRequest sends a request to the specified Shogun endpoint with the proper authentication
func (c *ShogunClient) DoShogunRequest(endpoint string) (response []byte, err error) {
	req, err := http.NewRequest("GET", "https://bugyo.hac.lp1.eshop.nintendo.net/shogun/v1"+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-DeviceAuthorization", c.DauthToken)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return bytes, nil
}

// GetNsID returns the NS ID for the given Title ID
func (c *ShogunClient) GetNsID(tid string) (nsID int64, err error) {
	resp, err := c.DoShogunRequest("/contents/ids?shop_id=4&lang=en&country=US&type=title&title_ids=" + tid)
	if err != nil {
		return 0, err
	}

	if string(resp) == "{\"id_pairs\":[]}" {
		return 0, errors.New("NS ID not avaliable for this title")
	}

	id, err := jsonparser.GetInt(resp, "id_pairs", "[0]", "id")
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetTitleData returns a Title instance for the given NS ID
func (c *ShogunClient) GetTitleData(nsID int64) (title *Title, err error) {
	resp, err := c.DoShogunRequest(fmt.Sprintf("/titles/%d?shop_id=4&lang=en&country=US", nsID))
	if err != nil {
		return &Title{}, err
	}

	id, err := jsonparser.GetInt(resp, "id")
	if err != nil {
		return &Title{}, err
	}

	name, err := jsonparser.GetString(resp, "formal_name")
	if err != nil {
		return &Title{}, err
	}

	bannerURL, err := jsonparser.GetString(resp, "hero_banner_url")
	if err != nil {
		return &Title{}, err
	}
	bannerURL = "https://bugyo.hac.lp1.eshop.nintendo.net" + bannerURL

	releaseDate, err := jsonparser.GetString(resp, "release_date_on_eshop")
	if err != nil {
		return &Title{}, err
	}

	isNew, err := jsonparser.GetBoolean(resp, "is_new")
	if err != nil {
		return &Title{}, err
	}

	description, err := jsonparser.GetString(resp, "description")
	if err != nil {
		return &Title{}, err
	}

	genre, err := jsonparser.GetString(resp, "genre")
	if err != nil {
		return &Title{}, err
	}

	size, err := jsonparser.GetInt(resp, "total_rom_size")
	if err != nil {
		return &Title{}, err
	}

	screenshots := []string{}
	jsonparser.ArrayEach(resp, func(value []byte, value_type jsonparser.ValueType, offset int, err error) {
		// todo: add error checking
		url, _ := jsonparser.GetString(value, "images", "[0]", "url")
		screenshots = append(screenshots, "https://bugyo.hac.lp1.eshop.nintendo.net"+string(url))
	}, "screenshots")

	movies := []*Movie{}
	jsonparser.ArrayEach(resp, func(value []byte, value_type jsonparser.ValueType, offset int, err error) {
		// todo: add error checking
		url, _ := jsonparser.GetString(value, "movie_url")
		thumbnailURL, _ := jsonparser.GetString(value, "thumbnail_url")

		movies = append(movies, &Movie{
			"https://bugyo.hac.lp1.eshop.nintendo.net" + url,
			"https://bugyo.hac.lp1.eshop.nintendo.net" + thumbnailURL,
		})
	}, "movies")

	pubID, err := jsonparser.GetInt(resp, "publisher", "id")
	if err != nil {
		return &Title{}, err
	}

	pubName, err := jsonparser.GetString(resp, "publisher", "name")
	if err != nil {
		return &Title{}, err
	}

	titleID, err := jsonparser.GetString(resp, "applications", "[0]", "id")
	if err != nil {
		return &Title{}, err
	}

	iconURL, err := jsonparser.GetString(resp, "applications", "[0]", "image_url")
	if err != nil {
		return &Title{}, err
	}
	iconURL = "https://bugyo.hac.lp1.eshop.nintendo.net" + iconURL

	return &Title{
		id,
		name,
		bannerURL,
		releaseDate,
		isNew,
		description,
		genre,
		size,
		screenshots,
		movies,
		&Publisher{
			pubID,
			pubName,
		},
		titleID,
		iconURL,
	}, nil
}
