package application

type ResponseJson struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type Area struct {
	Country string `json:"country"` // 区域
}

type UserProfile struct {
	NickName  string `json:"nick_name"`  // 昵称
	HeadImage string `json:"head_image"` // 头像
}

type PetrolRecordRes struct {
	Id         int     `json:"id"`
	Volume     float64 `json:"volume"`      // 体积
	Version    string  `json:"version"`     // 油号
	Price      float64 `json:"price"`       // 价格
	Mileage    float64 `json:"mileage"`     // 里程数
	Name       string  `json:"name"`        // 加油站名称
	CreateTime int     `json:"create_time"` // 加油时间
}

type ShareInfoRes struct {
	Msg string `json:"msg"` // 内容
	Img string `json:"img"` // 图片
}

type DeleteRecordReq struct {
	Id int `json:"id"` // 加油记录
}

type UserRecordRes struct {
	UserProfile
	LastQtrip        float64           `json:"last_qtrip"`        // 最近油耗
	AvgQtrip         float64           `json:"avg_qtrip"`         // 平均油耗
	AvgTrip          float64           `json:"avg_trip"`          // 平均行程
	RealMileage      float64           `json:"real_mileage"`      // 表显里程
	LastMileage      float64           `json:"last_mileage"`      // 统计里程数
	CumulativeDosage float64           `json:"cumulative_dosage"` // 累计加油量
	Records          []PetrolRecordRes `json:"records"`           // 加油记录
}

type AdvertisingReq struct {
	Location string `json:"location" enums:"home,display"`            // 位置，home 首页 display 发现页
	Type     string `json:"type" enums:"advertising,notice,carousel"` // 类型 advertising图片广告,notice 文字通知,carousel 轮播图
}

type ProfileQuery struct {
	UserId        int     `json:"user_id"`     // 用户ID
	MinMileage    float64 `json:"min_mileage"` // 最小里程
	MaxMileage    float64 `json:"max_mileage"`
	SumVolume     float64 `json:"sum_volume"`
	MinCreateTime int     `json:"min_create_time"`
	MaxCreateTime int     `json:"max_create_time"`
}
type LastQtrip struct {
	LastQtrip float64 `json:"last_qtrip"`
}

type NearbyStationRes struct {
	Station
	Distance float64 `json:"distance"` // 距离/ km
}

type StationSearch struct {
	Name string `json:"name"`
}

type NearbyReq struct {
	// 汽油号:
	Num string `json:"num" swaggertype:"string" enums:"92,95,98,0"`
	// 排序方式:
	// * price - 价格
	// * distance - 距离
	// * smart - 智能
	OrderBy string `json:"order_by" swaggertype:"string" enums:"price,distance,smart"`
	// 县，接口获取:
	Area string `json:"area"`
}

type LocationType struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type LocalArea struct {
	Country  string `json:"country"`  // 国家
	Province string `json:"province"` // 省份
	City     string `json:"city"`     // 市
	District string `json:"district"` // 区县
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
