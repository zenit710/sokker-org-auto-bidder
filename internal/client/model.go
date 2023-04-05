package client

// Time layout used by sokker.org API
const TimeLayout = "2006-01-02 15:04:05"

// loginReqBody represents request body sent to auth endpoint
type loginReqBody struct {
	Login string `json:"login"`
	Pass string `json:"password"`
}

// bidReqBody represents request body sent to bid endpoint
type bidReqBody struct {
	Value uint `json:"value"`
}

// bidState represens response body part describing player bid status
type bidState struct {
	Value uint
}

// playerInfoResponse represents response body from player info endpoint
type playerInfoResponse struct {
	Transfer struct {
		Deadline struct {
			Date string
			Timezone string
		}
		Price struct {
			MinBid bidState
		}
		BuyerId uint
	}
}

// transferInfoResponse represents response body from player bid endpoint
type transferInfoResponse struct {
	Deadline struct {
		Date struct {
			Date string
			Timezone string
		}
	}
	Price struct {
		MinBid bidState
	}
	Buyer struct {
		Id uint
	}
}

// clubInfoResponse represents response from club info endpoint
type clubInfoResponse struct {
	Team struct {
		Id uint
	}
}
