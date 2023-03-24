package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sokker-org-auto-bidder/tools"
)

var _ Client = &httpClient{}

const (
	urlAuth = "https://sokker.org/api/auth/login"
	urlPlayerInfoFormat = "https://sokker.org/api/player/%d"
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

	res, err := http.Post(urlAuth, "application/json", bodyReader)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return ErrBadCredentials
	}

	return nil
}

func (s *httpClient) FetchPlayerInfo(id uint) (*playerInfoResponse ,error) {
	res, err := http.Get(fmt.Sprintf(urlPlayerInfoFormat, id))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	p := &playerInfoResponse{}
	if err = json.Unmarshal(body, p); err != nil {
		return nil, err
	}

	return p, nil
}
