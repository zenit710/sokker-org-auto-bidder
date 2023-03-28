package client

// Time layout used by sokker.org API
const TimeLayout = "2006-01-02 15:04:05"

// request models

type loginReqBody struct {
	Login string `json:"login"`
	Pass string `json:"password"`
}

type bidReqBody struct {
	Value uint `json:"value"`
}


// response models

type playerInfoResponse struct {
	Transfer struct {
		Deadline struct {
			Date string
			Timezone string
		}
		Price struct {
			MinBid bidState
		}
	}
}

type bidState struct {
	Value uint
}
