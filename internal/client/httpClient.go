package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sokker-org-auto-bidder/tools"
)

var _ Client = &httpClient{}

const (
	AUTH_URL = "https://sokker.org/api/auth/login"
)

type httpClient struct {
	user string
	pass string
	sessId string
}

func NewHttpClient(user, pass string) *httpClient {
	return &httpClient{user: user, pass: pass, sessId: tools.String(26)}
}

func (s *httpClient) Auth() error {
	body := &loginReqBody{Login: s.user, Pass: s.pass}
	jsonBody, err := json.Marshal(body)
    if err != nil {
        return err
    }
	bodyReader := bytes.NewReader(jsonBody)

	res, err := http.Post(AUTH_URL, "application/json", bodyReader)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return ErrBadCredentials
	}

	return nil
}
