package application


type ResponseJson struct {
	Code int `json:"code"`
	Msg	string `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type UserProfileRes struct {
	NickName         string  `json:"nick_name"`
	HeadImage        string  `json:"head_image"`
	LastQtrip        float64 `json:"last_qtrip"`
	AvgQtrip         float64 `json:"avg_qtrip"`
	AvgTrip          float64 `json:"avg_trip"`
	RealMileage      float64 `json:"real_mileage"`
	LastMileage      float64 `json:"last_mileage"`
	CumulativeDosage float64 `json:"cumulative_dosage"`
}

type NearbyStationRes struct {
	Station
	Distance float64 `json:"distance"`
}

type NearbyReq struct {
	Num     string `json:"num"`
	OrderBy string `json:"order_by"`
	Area    string `json:"area"`
}

type LocationType struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type WXAuthRes struct {
	Errcode    int    `json:"errcode"`
	Unionid    string `json:"unionid"`
	Errmsg     string `json:"errmsg"`
	SessionKey string `json:"session_key"`
	Openid     string `json:"openid"`
}

type LoginReq struct {
	UserName string `json:"user_name"`
	PassWord string `json:"pass_word"`
}