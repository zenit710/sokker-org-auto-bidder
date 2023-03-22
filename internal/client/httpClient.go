package client

import (
	"bytes"
	"fmt"
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
	jsonBody := []byte(fmt.Sprintf(`{"login":"%s", "password":"%s"}`, s.user, s.pass))
	bodyReader := bytes.NewReader(jsonBody)

	res, err := http.Post(AUTH_URL, "application/json", bodyReader)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("could not authorize")
	}

	return nil
}
