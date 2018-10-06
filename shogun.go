package libhac

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

func (c *HacClient) doShogunRequest(endpoint string) (response []byte, err error) {
	resp, err := c.DoRequest("GET", fmt.Sprintf("https://bugyo.hac.lp1.eshop.nintendo.net/shogun/v1%s",
		endpoint), true, false)

	bytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return bytes, nil
}

func (c *HacClient) TestDauthToken() error {
	resp, err := c.doShogunRequest("/contents/ids?shop_id=4&lang=en&country=US&type=title&title_ids=999")
	if err != nil || string(resp) != "{\"id_pairs\":[]}" {
		return errors.New("edge token is invalid!")
	}

	return nil
}

func (c *HacClient) GetNSID(tid string) (nsID int, err error) {
	resp, err := c.doShogunRequest(fmt.Sprintf("/contents/ids?shop_id=4&lang=en&country=US&type=title&title_ids=%s"))
	if err != nil {
		return -1, err
	}

	if string(resp) == "{\"id_pairs\":[]}" {
		return -1, errors.New("ns id not avaliable for this title")
	}

	r := idResponse{}

	err = json.Unmarshal(resp, r)
	if err != nil {
		return -1, err
	}

	return r.IDPairs[0].ID, nil
}

func (c *HacClient) GetTitleData(nsID int) (title Title, err error) {
	resp, err := c.doShogunRequest(fmt.Sprintf("/titles/%d?shop_id=4&lang=en&country=US", nsID))
	if err != nil {
		return Title{}, err
	}

	t := Title{}

	err = json.Unmarshal(resp, t)
	if err != nil {
		return Title{}, err
	}

	return t, nil
}
