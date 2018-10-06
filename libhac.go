package libhac

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

type HacClient struct {
	HTTP       *http.Client
	DauthToken string
	EdgeToken  string
}

func NewHacClient(deviceCert, deviceKey,
	shopCert, shopKey, dauthToken, edgeToken string) (HacClient, error) {

	certs := []tls.Certificate{}

	if deviceCert != "" && deviceKey != "" {
		device, err := tls.LoadX509KeyPair(deviceCert, deviceKey)
		if err != nil {
			return HacClient{}, err
		}
		certs = append(certs, device)
	}

	// maybe hardcode this as it's common? todo: research how loadx509 loads the cert
	if shopCert != "" && shopKey != "" {
		shop, err := tls.LoadX509KeyPair(shopCert, shopKey)
		if err != nil {
			return HacClient{}, err
		}
		certs = append(certs, shop)
	}

	return HacClient{
		&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Certificates:       certs,
					InsecureSkipVerify: true,
				},
			},
		},
		dauthToken,
		edgeToken,
	}, nil
}

func (c *HacClient) DoRequest(method, url string, sendDauthToken, sendEdgeToken bool) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return &http.Response{}, err
	}

	if sendDauthToken {
		req.Header.Set("X-DeviceAuthorization", fmt.Sprintf("Bearer %s", c.DauthToken))
	}

	if sendEdgeToken {
		req.Header.Set("X-Nintendo-DenebEdgeToken", c.EdgeToken)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return &http.Response{}, err
	}

	return resp, nil
}
