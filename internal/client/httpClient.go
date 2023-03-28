package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sokker-org-auto-bidder/tools"
)

var _ Client = &httpClient{}

const (
	urlAuth = "https://sokker.org/api/auth/login"
	urlClubInfo = "https://sokker.org/api/current"
	urlPlayerInfoFormat = "https://sokker.org/api/player/%d"
	urlPlayerBidFormat = "https://sokker.org/api/transfer/%d/bid"
)

type httpClient struct {
	user string
	pass string
	sessId string
	auth bool
}

func NewHttpClient(user, pass string) *httpClient {
	return &httpClient{user: user, pass: pass, auth: false, sessId: tools.String(26)}
}

func (s *httpClient) Auth() (*clubInfoResponse, error) {
	// prepare auth request body
	body := &loginReqBody{Login: s.user, Pass: s.pass}

	// make http request
	res, err := s.makeRequest(urlAuth, http.MethodPost, body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, ErrBadCredentials
	}

	s.auth = true

	return s.ClubInfo()
}

func (s*httpClient) ClubInfo() (*clubInfoResponse, error) {
	// make http request
	res, err := s.makeRequest(urlClubInfo, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("unauthorized")
	}

	c := &clubInfoResponse{}
	err = extractResponseObject(res, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *httpClient) FetchPlayerInfo(id uint) (*playerInfoResponse ,error) {
	res, err := http.Get(fmt.Sprintf(urlPlayerInfoFormat, id))
	if err != nil {
		return nil, err
	}

	p := &playerInfoResponse{}
	err = extractResponseObject(res, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *httpClient) Bid(id, price uint) (*playerInfoResponse, error) {
	// prepare req params
	body := &bidReqBody{Value: price}
	bidUrl := fmt.Sprintf(urlPlayerBidFormat, id)

	// make http request
	res, err := s.makeRequest(bidUrl, http.MethodPut, body)
	if err != nil {
		fmt.Println("request sent but error")
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code: %d", res.StatusCode)
	}

	p := &playerInfoResponse{}
	err = extractResponseObject(res, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *httpClient) makeRequest(url string, method string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader = nil
	
	if body != nil {
		// prepare auth request body
		jsonBody, err := json.Marshal(body)
		if err != nil {
			fmt.Println("marshal error, body to json")
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	// prepare bid request
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		fmt.Println("request error")
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", fmt.Sprintf("PHPSESSID=%s", s.sessId))

	// // make http request
	return http.DefaultClient.Do(req)
}

func extractResponseObject(res *http.Response, obj interface{}) (error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("cannot read body")
		return err
	}
	res.Body.Close()

	if err = json.Unmarshal(body, obj); err != nil {
		fmt.Println("unmarshal error")
		fmt.Println(string(body))
		return err
	}

	fmt.Println(obj)

	return nil
}
