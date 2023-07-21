package client_test

import (
	"errors"
	"net/http"
	"regexp"
	"sokker-org-auto-bidder/internal/client"
	"sokker-org-auto-bidder/internal/repository/session"
	"testing"

	"github.com/jarcoal/httpmock"
)

const (
	urlAuth     = "https://sokker.org/api/auth/login"
	urlClubInfo = "https://sokker.org/api/current"
	// urlPlayerBidFormat = "https://sokker.org/api/transfer/%d/bid"
)

var urlPlayerInfoRegex = regexp.MustCompile(`https://sokker.org/api/player/\d+`)

func createHttpClient() (*session.MockSessionRepository, client.Client) {
	r := &session.MockSessionRepository{}
	r.On("Get").Return("key", nil)
	c := client.NewHttpClient("user", "pass", r)
	return r, c
}

func TestAuthFailureWhenRequestFailed(t *testing.T) {
	var expectedError *client.ErrRequestFailed
	_, c := createHttpClient()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	if _, err := c.Auth(); err == nil || !errors.As(err, &expectedError) {
		t.Errorf("expected '%T' error but '%T' returned", expectedError, err)
	}
}

func TestAuthFailureWhenBadCredentials(t *testing.T) {
	expectedErr := client.ErrBadCredentials
	_, c := createHttpClient()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		http.MethodPost,
		urlAuth,
		httpmock.NewStringResponder(http.StatusBadRequest, ""),
	)

	if _, err := c.Auth(); err == nil || !errors.Is(err, expectedErr) {
		t.Errorf("expected '%v' error but '%v' returned", expectedErr, err)
	}
}

func TestAuthSuccess(t *testing.T) {
	_, c := createHttpClient()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		http.MethodPost,
		urlAuth,
		httpmock.NewStringResponder(http.StatusOK, ""),
	)
	// test depends on club info method
	httpmock.RegisterResponder(
		http.MethodGet,
		urlClubInfo,
		httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("./mocks/clubInfo.json")),
	)

	info, err := c.Auth()
	if err != nil {
		t.Errorf("expected <nil> error but '%v' returned", err)
	}
	if info == nil {
		t.Error("expected team info, <nil> returned")
	}
}

func TestClubInfoFailureWhenRequestFailed(t *testing.T) {
	var expectedError *client.ErrRequestFailed
	_, c := createHttpClient()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	if _, err := c.ClubInfo(); err == nil || !errors.As(err, &expectedError) {
		t.Errorf("expected '%T' error but '%T' returned", expectedError, err)
	}
}

func TestClubInfoFailureWhenUnauthorized(t *testing.T) {
	expectedErr := client.ErrUnauthorized
	_, c := createHttpClient()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		http.MethodGet,
		urlClubInfo,
		httpmock.NewStringResponder(http.StatusBadRequest, ""),
	)

	if _, err := c.ClubInfo(); err == nil || !errors.Is(err, expectedErr) {
		t.Errorf("expected '%v' error but '%v' returned", expectedErr, err)
	}
}

func TestClubInfoFailureWhenCanNotParseResponse(t *testing.T) {
	var expectedErr *client.ErrResponseParseFailed
	_, c := createHttpClient()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		http.MethodGet,
		urlClubInfo,
		httpmock.NewStringResponder(http.StatusOK, ""),
	)

	if _, err := c.ClubInfo(); err == nil || !errors.As(err, &expectedErr) {
		t.Errorf("expected '%T' error but '%T' returned", expectedErr, err)
	}
}

func TestClubInfoSuccess(t *testing.T) {
	_, c := createHttpClient()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		http.MethodGet,
		urlClubInfo,
		httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("./mocks/clubInfo.json")),
	)

	info, err := c.ClubInfo()
	if err != nil {
		t.Errorf("expected <nil> error but '%v' returned", err)
	}
	if info == nil {
		t.Error("expected team info, <nil> returned")
	}
}

func TestFetchPlayerInfoFailureWhenRequestFailed(t *testing.T) {
	var expectedError *client.ErrRequestFailed
	_, c := createHttpClient()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	if _, err := c.FetchPlayerInfo(0); err == nil || !errors.As(err, &expectedError) {
		t.Errorf("expected '%T' error but '%T' returned", expectedError, err)
	}
}

func TestFetchPlayerInfoFailureWhenUnavailble(t *testing.T) {
	var expectedErr *client.ErrResourceUnavailable
	_, c := createHttpClient()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterRegexpResponder(
		http.MethodGet,
		urlPlayerInfoRegex,
		httpmock.NewStringResponder(http.StatusNotFound, ""),
	)

	if _, err := c.FetchPlayerInfo(0); err == nil || !errors.As(err, &expectedErr) {
		t.Errorf("expected '%T' error but '%T' returned", expectedErr, err)
	}
}

func TestFetchPlayerInfoFailureWhenCanNotParseResponse(t *testing.T) {
	var expectedErr *client.ErrResponseParseFailed
	_, c := createHttpClient()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterRegexpResponder(
		http.MethodGet,
		urlPlayerInfoRegex,
		httpmock.NewStringResponder(http.StatusOK, ""),
	)

	if _, err := c.FetchPlayerInfo(0); err == nil || !errors.As(err, &expectedErr) {
		t.Errorf("expected '%T' error but '%T' returned", expectedErr, err)
	}
}

func TestFetchPlayerInfoSuccess(t *testing.T) {
	_, c := createHttpClient()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterRegexpResponder(
		http.MethodGet,
		urlPlayerInfoRegex,
		httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("./mocks/playerinfo.json")),
	)

	info, err := c.FetchPlayerInfo(0)
	if err != nil {
		t.Errorf("expected <nil> error but '%v' returned", err)
	}
	if info == nil {
		t.Error("expected team info, <nil> returned")
	}
}
