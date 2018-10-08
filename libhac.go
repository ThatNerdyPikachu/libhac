package libhac

import (
	"crypto/tls"
	"errors"
	"net/http"
)

type HacClient struct {
	DeviceCert tls.Certificate
	ShopCert   tls.Certificate
	DauthToken string
	EdgeToken  string
}

func NewHacClient(deviceCert, deviceKey, shopCert, shopKey, dauthToken, edgeToken string) (HacClient, error) {
	// lolwut
	err := errors.New("")

	device := tls.Certificate{}
	if deviceCert != "" && deviceKey != "" {
		device, err = tls.LoadX509KeyPair(deviceCert, deviceKey)
		if err != nil {
			return HacClient{}, err
		}
	}

	// maybe hardcode this as it's common? todo: research how loadx509 loads the cert
	shop := tls.Certificate{}
	if shopCert != "" && shopKey != "" {
		shop, err = tls.LoadX509KeyPair(shopCert, shopKey)
		if err != nil {
			return HacClient{}, err
		}
	}

	return HacClient{
		device,
		shop,
		dauthToken,
		edgeToken,
	}, nil
}

func (c *HacClient) DoRequest(method, url string, certs []tls.Certificate, sendDauthToken, sendEdgeToken bool) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return &http.Response{}, err
	}

	if sendDauthToken {
		req.Header.Set("X-DeviceAuthorization", c.DauthToken)
	}

	if sendEdgeToken {
		req.Header.Set("X-Nintendo-DenebEdgeToken", c.EdgeToken)
	}

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       certs,
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}

	return resp, nil
}
