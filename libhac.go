package libhac

import (
	"crypto/tls"
	"errors"
	"net/http"
)

type HacClient struct {
	DeviceCert *tls.Certificate
	ShopCert   *tls.Certificate
	DauthToken string
	EdgeToken  string
}

func NewHacClient(deviceCert, deviceKey []byte, dauthToken, edgeToken string) (*HacClient, error) {
	// lolwut
	err := errors.New("")

	device := tls.Certificate{}
	if deviceCert != nil && deviceKey != nil {
		device, err = tls.X509KeyPair(deviceCert, deviceKey)
		if err != nil {
			return nil, err
		}
	}

	shop, err := tls.X509KeyPair([]byte(`-----BEGIN CERTIFICATE-----
MIIEgjCCA2qgAwIBAgICAZwwDQYJKoZIhvcNAQELBQAwbTELMAkGA1UEBhMCVVMx
EzARBgNVBAgTCldhc2hpbmd0b24xITAfBgNVBAoTGE5pbnRlbmRvIG9mIEFtZXJp
Y2EgSW5jLjELMAkGA1UECxMCSVMxGTAXBgNVBAMTEE5pbnRlbmRvIENBIC0gRzMw
HhcNMTYxMTEyMDExNzIxWhcNNDIwNDMwMDExNzIxWjCBkTELMAkGA1UEBhMCSlAx
DjAMBgNVBAgTBUt5b3RvMRIwEAYDVQQHEwlNaW5hbWkta3UxGzAZBgNVBAoTEk5p
bnRlbmRvIENvLiwgTHRkLjExMC8GA1UECxQoTmV0d29yayAmIEluZm9ybWF0aW9u
IFN5c3RlbXMgRGVwYXJ0bWVudDEOMAwGA1UEAxMFU2hvcE4wggEiMA0GCSqGSIb3
DQEBAQUAA4IBDwAwggEKAoIBAQDQfWbpeyuPMvT1cn1cCg/wbnbcKVztWf4U87rQ
QJq1VVHhZ/aD72BitbbsY89gvhHPP53VQFCt0vxVTpbxWkeHaqTpAmfpCVEjAdaj
aOueHvziJhWvpOc20jescoVCXdn20kUk1+PLsOBaf0GqjGZ+C/piQ5Ti+8jQMRPZ
WIKdNP03iKPZRw4QvrLs9NFXWmCyq26iQkfDvI2Pm5VrrF9QLWHmZ4Ra0SkQyOnN
LJQ9QerrIONK+rPH1Xxxk3SC0TUZT7d9i8h5x1cj+Kyl4QT8pkjNQFdKBTea2xXF
0gQPakoUcUeJiG7fuf9p1/addvkyPQQdy2gVCEh63D4R5TsDAgMBAAGjggEFMIIB
ATAJBgNVHRMEAjAAMB0GA1UdDgQWBBS2iQi+cWi+PjRiWUSfnKNBdIaC4DCBlwYD
VR0jBIGPMIGMgBQE097T/fDI68JZkof7H9c+cvjt+aFxpG8wbTELMAkGA1UEBhMC
VVMxEzARBgNVBAgTCldhc2hpbmd0b24xITAfBgNVBAoTGE5pbnRlbmRvIG9mIEFt
ZXJpY2EgSW5jLjELMAkGA1UECxMCSVMxGTAXBgNVBAMTEE5pbnRlbmRvIENBIC0g
RzOCAQEwOwYDVR0fBDQwMjAwoC6gLIYqaHR0cDovL2NybC5uaW50ZW5kby5jb20v
bmludGVuZG8tY2EtZzMuY3JsMA0GCSqGSIb3DQEBCwUAA4IBAQAwqs9jrGljEhJY
sASBeuDrsJ5Gck68Htp1TjV8OHqDRAxr+D3EEQpDvKGNCez2i9sGpENhFnNR8K6g
3vBWqihOLkgKmsO6U30Bbk9iLUoWcbH/KEo/XmQqE3aXls11KU4Cz/xNunmZzuAP
aDGg6FhBQSpDkN20OorehN+pN1ZM7d6vi7bk1L5kFij2KLX0tvEALKQAdAK4HuIj
MWWAC9kGGHP1Y6lkziF6LghXun6iIiR1RZN6aYerWO6FZxHiMvP6pOUjgmIBWYWu
oD47EW5BDppHnLJCdhzEk+knpWfhqxzulDSCVW2Qx6V02SjQ3eAVTWyeedyICRWK
h0+lMNh+
-----END CERTIFICATE-----`), []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA0H1m6XsrjzL09XJ9XAoP8G523Clc7Vn+FPO60ECatVVR4Wf2
g+9gYrW27GPPYL4Rzz+d1UBQrdL8VU6W8VpHh2qk6QJn6QlRIwHWo2jrnh784iYV
r6TnNtI3rHKFQl3Z9tJFJNfjy7DgWn9Bqoxmfgv6YkOU4vvI0DET2ViCnTT9N4ij
2UcOEL6y7PTRV1pgsqtuokJHw7yNj5uVa6xfUC1h5meEWtEpEMjpzSyUPUHq6yDj
Svqzx9V8cZN0gtE1GU+3fYvIecdXI/ispeEE/KZIzUBXSgU3mtsVxdIED2pKFHFH
iYhu37n/adf2nXb5Mj0EHctoFQhIetw+EeU7AwIDAQABAoIBAFOwctX4FjUmLEQ2
T/HZLCrD/LxFckLoY+B/MZcUx8VQWUzU0ZSGSzd0X9gl/IGF6lo53B5U8c3EqnuH
z3lUVvAZs9bAm3tkvQgDXeg8XpAbOkGBLiVSWWmdGrIDSlCooay9HPo8GvmRp41M
FczjDOU946T8IdC3ZxWDDceqCswVVqx6t//Nfbw9ns9k/M82cCbvQ3PIeSPtMjtF
Bw28NjVtYnF135NDAMPHuPww/4SUUKH/ukkhOqsleHsYwCWOFgowTy08sIDDUcFf
ypFhH/UF7Tof0myK9EZp+qF0r+Gc2Dr5ey16qVjnjSse9GM4Lw9d+ek82y1udpQS
OU6FUKECgYEA6qckJBZThT4tbKanC4Th/E73sV735yZUoWcsPE6UfHRrRtZHxzvD
SeD5UuF8zguzTTHANBCzdeUF+xBBETNo2UsJwB6dtc6wveDBq4jfdfob8HOD2hO+
fIxpN/A80gMS3OL2YieOkA1RYvpZSMW6HLG7IEkhtUPUcshDz/Ka5TMCgYEA43Ty
Vk0j2EpgCyq3uYAcv9mQHoG8mAVNoIqmYokqfiw13rOYnrg1bdlYEBhS81xE6Nio
/qBVqljKd94BjsvRxX6DOJCb+41cbNKq6fEIl08qWAnkEgxFxOGjv51SxsSa+36/
vj534c3MCNhHnV2zygHa1LTOFt8BQGpBV25ysvECgYBv1QyybjmL7SA1Ag/mu4QA
fZoq3HoFng0Dh2pbq4+06L8f2YwH9gGSujfZOkIC6SyhA674oWUGMQpyp0J7cibs
OIkDGRRMy58MP4vNtsPZNEQQit0yF+rW7QxQbpxt8q0QvNJz4SpemRpKs8nNlv6m
IkZukqwP+a3MrokzGAqeZQKBgQDfg59UQS+VLtkcgBjU19joihH+eHzsYmG6/iSs
tZ+lJyq0soWqitguHyNThaOV87jfm3DN7p7f66riia28Nfvod/7YyfoOAwatBYPg
qfRIRAoXQ3j0Z0z0QMsZ065xFug5dX9UoqxJn2L92hLdyCORwarZ1OakQPZI52FY
Wnep8QKBgQCLz/TeN8DapOZJ0ylBzo2F2Oa2RkgGYEtMevBqeIWSVwGG4d89cIig
1MaEXMrWn1rmp/NfCw7AzsFSpnxzLEGrd0yNUTYm32NaziJlCMRqDZn7pz/saUU5
9wzSlH1LFZ4aRiB+UPiQj0g0h8ivE2ewtfDpbEYnbb9jN7vIBroZAw==
-----END RSA PRIVATE KEY-----`))
	if err != nil {
		return nil, err
	}

	return &HacClient{
		&device,
		&shop,
		dauthToken,
		edgeToken,
	}, nil
}

func (c *HacClient) DoRequest(method, url string, certs []tls.Certificate, sendDauthToken, sendEdgeToken bool) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	if sendDauthToken {
		req.Header.Set("X-DeviceAuthorization", c.DauthToken)
	}

	if sendEdgeToken {
		req.Header.Set("X-Nintendo-DenebEdgeToken", c.EdgeToken)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       certs,
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
