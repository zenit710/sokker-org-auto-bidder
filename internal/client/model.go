package client

// Time layout used by sokker.org API
const TimeLayout = "2006-01-02 15:04:05"

// loginReqBody represents request body sent to auth endpoint
type loginReqBody struct {
	Login string `json:"login"`
	Pass  string `json:"password"`
}

// bidReqBody represents request body sent to bid endpoint
type bidReqBody struct {
	Value uint `json:"value"`
}

// playerInfoResponse represents response body from player info endpoint
type playerInfoResponse struct {
	Transfer transfer
}

// transferInfoResponse represents response body from player bid endpoint
type transferInfoResponse struct {
	Deadline transferInfoDeadline
	Price    price
	Buyer    buyer
}

// clubInfoResponse represents response from club info endpoint
type clubInfoResponse struct {
	Team team
}

// transfer stores data about current player transfer state
type transfer struct {
	Deadline playerInfoDeadline
	Price    price
	BuyerId  uint
}

/*
playerInfoDeadline stores data about player transfer end date when asked about player info
you may be want to use transferInfoDeadline instead
*/
type playerInfoDeadline struct {
	Date     string
	Timezone string
}

/*
transferInfoDeadline stores data about player transfer end date when asked about transfer info
you may be want to use playerInfoDeadline instead
*/
type transferInfoDeadline struct {
	Date date
}

// date stores data about date itself and the timezone
type date struct {
	Date     string
	Timezone string
}

// price stores data about player price
type price struct {
	MinBid bidState
}

// bidState represens response body part describing player bid status
type bidState struct {
	Value uint
}

// buyer stores data about current bid leader
type buyer struct {
	Id uint
}

// team struct store data about sokker team
type team struct {
	Id uint
}
