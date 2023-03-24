package client

type loginReqBody struct {
	Login string `json:"login"`
	Pass string `json:"password"`
}

type playerInfoResponse struct {
	Transfer struct {
		Deadline struct {
			Date string
		}
		Price struct {
			MinBid bidState
		}
	}
}

type bidState struct {
	Value uint
}
