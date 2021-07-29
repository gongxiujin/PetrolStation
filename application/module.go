package application

import "gorm.io/gorm"

var DB *gorm.DB

type PetrolDaily struct {
	ID         int     `gorm:"<-;column:id;primaryKey"` // id
	Province   string  `json:"province" gorm:"<-;column:province"` // 省份，忽略
	Version    string  `json:"version" gorm:"column:version"` // 油号
	Price      float64 `json:"price" gorm:"column:price"` // 价格
	Day        string  `json:"day" gorm:"column:day"` // 日期，忽略
	CreateTime int64   `gorm:"<-;autoCreateTime"` // 创建时间, 忽略
}

type Station struct {
	ID int `json:"id,omitempty" swaggertype:"integer" gorm:"column:id;primaryKey"` // id:
	Province string `json:"province,omitempty" swaggertype:"string" example:"省份" gorm:"column:province"` // 省份
	City string `json:"city,omitempty" swaggertype:"string" gorm:"column:city"` // 市
	Country string `json:"country,omitempty" swaggertype:"string" gorm:"column:country"` // 县
	Address string `json:"address,omitempty" swaggertype:"string" gorm:"column:address"` // 详细地址
	Phone string `json:"phone,omitempty" swaggertype:"string" gorm:"column:phone"` // 电话
	Logo string `json:"logo,omitempty" swaggertype:"string" gorm:"column:logo"` // 加油站logo
	Name string `json:"name,omitempty" swaggertype:"string" gorm:"column:name;type:varchar;size:50;index"` // 加油站名字
	Longitude float64 `json:"longitude,omitempty" swaggertype:"number" gorm:"column:longitude"` // 经度
	Latitude float64 `json:"latitude,omitempty" swaggertype:"number" gorm:"column:latitude"` // 维度
	CreateTime int64         `json:"create_time,omitempty" swaggertype:"integer" gorm:"autoCreateTime"` // 创建时间
	Petrol     []PetrolPrice `json:"petrol,omitempty" gorm:"ForeignKey:StationID;save_associations:false"` // 油价信息
}

type PetrolPrice struct {
	ID        int  `json:"id" swaggertype:"integer" gorm:"column:id;primaryKey"` // id 可以忽略
	StationID uint `json:"station_id" swaggertype:"integer" gorm:"index"` // 可以忽略
	Version string `json:"version" swaggertype:"string"` // 油号
	Price float64 `json:"price" swaggertype:"number"` // 价格
	Day        string `json:"day" swaggertype:"string"` // 天，忽略
	CreateTime int64  `json:"create_time" swaggertype:"integer" gorm:"autoCreateTime"` // 创建时间
}

type Advertising struct {
	ID         int    `json:"id,omitempty" gorm:"<-:create;column:id;primaryKey"` // ID
	Location   string `json:"location" gorm:"<-:create;column:location" enums:"home,display"` // 位置，home 首页 display 发现页
	Type       string `json:"type" gorm:"<-:create;column:type" enums:"advertising,notice,carousel"` // 类型 advertising图片广告,notice 文字通知,carousel 轮播图
	Publish    bool   `json:"publish" gorm:"<-;column:publish"` // 是否展示，可忽略
	Content    string `json:"content" gorm:"column:content"` // 内容，图片为跳转连接，文字为内容
	Url        string `json:"url" gorm:"column:url"` // 图片的真实地址
	CreateTime int64  `json:"create_time,omitempty" gorm:"<-;autoCreateTime"` // 可忽略
}

type Users struct {
	ID         int    `gorm:"column:id;primaryKey"`
	UserName   string `gorm:"column:user_name"`
	PassWord   string `gorm:"column:pass_word"` //
	Avator     string `gorm:"column:avator"` // 头像
	Phone      string
	NickName   string // 昵称
	Active     bool
	Longitude  float64 `json:"longitude,omitempty"  gorm:"column:longitude"`
	Latitude   float64 `json:"latitude,omitempty"  gorm:"column:latitude"`
	CreateTime int64   `gorm:"autoCreateTime"`
	UpdateTime int64   `gorm:"autoUpdateTime"`
}

type PetrolRecord struct {
	ID         int     `gorm:"column:id;primaryKey"` // id
	UserId     int     `gorm:"index"` // 用户id
	Volume     float64 `json:"volume" gorm:"column:volume"` // 加油的体积
	Price      float64 `json:"price" gorm:"column:price"` // 单价
	StationId  uint    `json:"station_id" gorm:"column:station_id"` // 加油站id
	Version    string  `json:"version" gorm:"column:version"` // 油号
	Mileage    float64 `json:"mileage" gorm:"column:mileage"` // 当前公里数
	CreateTime int     `json:"create_time" gorm:"column:create_time;autoCreateTime"` // 创建时间
}

type UserTrack struct {
	ID         int     `gorm:"column:id;primaryKey"`
	UserId     int     `json:"user_id" gorm:"column:user_id"` // 用户id
	Longitude  float64 `json:"longitude,omitempty"  gorm:"column:longitude"` // 用户的经度
	Latitude   float64 `json:"latitude,omitempty"  gorm:"column:latitude"` // 用户的纬度
	CreateTime int64   `gorm:"autoCreateTime"` // 创建时间
}
