package libhac

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type SuperflyTitle struct {
	ID      string `json:"title_id"`
	Version int    `json:"version"`
	Type    string `json:"title_type"`
}

func (c *HacClient) GetSuperflyResponse(tid string) ([]SuperflyTitle, error) {
	resp, err := c.DoRequest("GET", fmt.Sprintf("https://superfly.hac.lp1.d4c.nintendo.net/v1/a/%s/dv", tid),
		false, false)
	if err != nil {
		return []SuperflyTitle{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []SuperflyTitle{}, err
	}
	resp.Body.Close()

	t := []SuperflyTitle{}

	err = json.Unmarshal(body, &t)
	if err != nil {
		return []SuperflyTitle{}, err
	}

	return t, nil
}
