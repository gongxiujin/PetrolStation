package application


type ResponseJson struct {
	Code int `json:"code"`
	Msg	string `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}