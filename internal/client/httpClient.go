package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sokker-org-auto-bidder/internal/repository/session"
	"sokker-org-auto-bidder/tools"

	log "github.com/sirupsen/logrus"
)

var _ Client = &httpClient{}

const (
	// urlAuth is sokker API endpoint for user authorization
	urlAuth = "https://sokker.org/api/auth/login"
	// urlClubInfo is sokker API endpoint with user club info
	urlClubInfo = "https://sokker.org/api/current"
	// urlPlayerInfoFormat is sokker API endpoint (has to be filled with playerId) with player info
	urlPlayerInfoFormat = "https://sokker.org/api/player/%d"
	// urlPlayerBidFormat is sokker API endpoint (has to be filled with playerId) for transfer bid
	urlPlayerBidFormat = "https://sokker.org/api/transfer/%d/bid"
)

// httpClient handles connection with sokker API through http
type httpClient struct {
	user string
	pass string
	sessRepo session.SessionRepository
	sessId string
	auth bool
}

// NewHttpClient returns new HttpClient for sokker.org
func NewHttpClient(user, pass string, sessRepo session.SessionRepository) *httpClient {
	log.Trace("creating new http client for sokker.org")
	c := &httpClient{user: user, pass: pass, auth: false, sessRepo: sessRepo}
	c.resolveSessKey()
	
	return c
}

func (s *httpClient) Auth() (*clubInfoResponse, error) {
	log.Trace("prepare auth request body object")
	body := &loginReqBody{Login: s.user, Pass: s.pass}

	log.Trace("make auth request")
	res, err := s.makeRequest(urlAuth, http.MethodPost, body)
	if err != nil {
		log.Error(err)
		return nil, &ErrRequestFailed{"auth"}
	}

	log.Debugf("auth request http status code: %d", res.StatusCode)
	if res.StatusCode != http.StatusOK {
		return nil, ErrBadCredentials
	}

	s.auth = true
	log.Debug("authenication successful")

	return s.ClubInfo()
}

func (s*httpClient) ClubInfo() (*clubInfoResponse, error) {
	log.Trace("make club info request")
	res, err := s.makeRequest(urlClubInfo, http.MethodGet, nil)
	if err != nil {
		log.Error(err)
		return nil, &ErrRequestFailed{"club info"}
	}

	log.Debugf("club info request http status code: %d", res.StatusCode)
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("unauthorized")
	}

	log.Trace("parse club info response")
	c := &clubInfoResponse{}
	err = extractResponseObject(res, c)
	if err != nil {
		log.Error(err)
		return nil, &ErrResponseParseFailed{"club info"}
	}

	return c, nil
}

func (s *httpClient) FetchPlayerInfo(id uint) (*playerInfoResponse, error) {
	log.Tracef("make player (%d) info request", id)
	res, err := http.Get(fmt.Sprintf(urlPlayerInfoFormat, id))
	if err != nil {
		log.Error(err)
		return nil, &ErrRequestFailed{fmt.Sprintf("player (%d) info", id)}
	}

	log.Debugf("fetch player (%d) info request http status code: %d", id, res.StatusCode)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("player (%d) details can not be fetched", id)
	}

	log.Tracef("parse player (%d) info response", id)
	p := &playerInfoResponse{}
	err = extractResponseObject(res, p)
	if err != nil {
		log.Error(err)
		return nil, &ErrResponseParseFailed{fmt.Sprintf("player (%d) info", id)}
	}

	return p, nil
}

func (s *httpClient) Bid(id, price uint) (*transferInfoResponse, error) {
	log.Trace("prepare bid request body object")
	body := &bidReqBody{Value: price}
	bidUrl := fmt.Sprintf(urlPlayerBidFormat, id)

	log.Tracef("make player (%d) bid request")
	res, err := s.makeRequest(bidUrl, http.MethodPut, body)
	if err != nil {
		log.Error(err)
		return nil, &ErrRequestFailed{fmt.Sprintf("player (%d) bid", id)}
	}

	log.Debugf("player (%d) bid request http status code: %d", id, res.StatusCode)
	if res.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("no funds for player (%d) bid", id)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("player (%d) bid response failed", id)
	}

	log.Tracef("parse player (%d) bid response", id)
	p := &transferInfoResponse{}
	err = extractResponseObject(res, p)
	if err != nil {
		log.Error(err)
		return nil, &ErrResponseParseFailed{fmt.Sprintf("player (%d) bid", id)}
	}

	return p, nil
}

/*
makeRequest sends new JSON http request with body made from interface{}.
PHPSESSID cookie is passed with this request.
*/
func (s *httpClient) makeRequest(url string, method string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader = nil
	
	if body != nil {
		log.Trace("request has body to send, convert to json")
		jsonBody, err := json.Marshal(body)
		if err != nil {
			log.Error(err)
			return nil, fmt.Errorf("%T object with value %v could not be transformed to JSON", body, body)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	log.Trace("prepare request")
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("could not create new http request")
	}
	log.Trace("set request headers (content-type, cookie)")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", fmt.Sprintf("PHPSESSID=%s", s.sessId))

	log.Trace("make http request")
	return http.DefaultClient.Do(req)
}

func (s *httpClient) resolveSessKey() {
	log.Trace("resolve http session key")
	dbKey, err := s.sessRepo.Get()
	if err == nil {
		log.Debug("use http session key from db")
		s.sessId = dbKey
		return
	}

	log.Warnf("could not get sess key from db (%v), new one will be generated", err)
	s.createNewSessKey()
}

func (s* httpClient) createNewSessKey() {
	log.Trace("generate new http sess key")
	s.sessId = tools.String(26)
	
	log.Trace("save new http session key to the db")
	if err := s.sessRepo.Save(s.sessId); err != nil {
		log.Warnf("http session key could not be stored: %v", err)
	}
}

// extractResponseObject parse response to interface{} 
func extractResponseObject(res *http.Response, obj interface{}) (error) {
	log.Trace("read http response body")
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("could not read http response body %v", res.Body)
	}
	res.Body.Close()

	log.Trace("parse json response to object")
	if err = json.Unmarshal(body, obj); err != nil {
		log.Error(err)
		return fmt.Errorf("could not transform http response to %T object", obj)
	}

	return nil
}
