package application

import "gorm.io/gorm"

var DB *gorm.DB

type PetrolDaily struct {
	ID         uint    `gorm:"<-;column:id;primaryKey"`
	Province   string  `json:"province" gorm:"<-;column:province"`
	Version    string  `json:"version" gorm:"column:version"`
	Price      float64 `json:"price" gorm:"column:price"`
	Day        string  `json:"day" gorm:"column:day"`
	CreateTime int64   `gorm:"<-;autoCreateTime:milli"`
}

type Station struct {
	ID         uint    `json:"id,omitempty" gorm:"column:id;primaryKey"`
	Province   string  `json:"province,omitempty" gorm:"column:province"`
	City       string  `json:"city,omitempty"  gorm:"column:city"`
	Country    string  `json:"country,omitempty"  gorm:"column:country"`
	Address    string  `json:"address,omitempty"  gorm:"column:address"`
	Phone      string  `json:"phone,omitempty"  gorm:"column:phone"`
	Logo       string  `json:"logo,omitempty"  gorm:"column:logo"`
	Name       string  `json:"name,omitempty"  gorm:"column:name;type:varchar;size:50;index"`
	Longitude  float64 `json:"longitude,omitempty"  gorm:"column:longitude"`
	Latitude   float64 `json:"latitude,omitempty"  gorm:"column:latitude"`
	CreateTime int64   `json:"create_time,omitempty" gorm:"autoCreateTime:milli"`
}

type PetrolPrice struct {
	ID         uint `gorm:"column:id;primaryKey"`
	StationID  uint `gorm:"index"`
	Station    Station
	Version    string
	Price      float64
	Day        string
	CreateTime int64 `gorm:"autoCreateTime:milli"`
}

type Advertising struct {
	ID         uint   `json:"id,omitempty" gorm:"<-:create;column:id;primaryKey"`
	Location   string `json:"location" gorm:"<-:create;column:location"`
	Type       string `json:"type" gorm:"<-:create;column:type"`
	Publish    bool   `json:"publish" gorm:"<-;column:publish"`
	Content    string `json:"content" gorm:"column:content"`
	Url        string `json:"url" gorm:"column:url"`
	CreateTime int64  `json:"create_time,omitempty" gorm:"<-;autoCreateTime:milli"`
}

type Users struct {
	ID         uint   `gorm:"column:id;primaryKey"`
	UserName   string `gorm:"column:user_name"`
	PassWord   string `gorm:"column:pass_word"`
	Avator     string
	Phone      string
	NickName   string
	Active     bool
	Longitude  float64 `json:"longitude,omitempty"  gorm:"column:longitude"`
	Latitude   float64 `json:"latitude,omitempty"  gorm:"column:latitude"`
	CreateTime int64   `gorm:"autoCreateTime:milli"`
	UpdateTime int64   `gorm:"autoUpdateTime:milli"`
}

type PetrolRecord struct {
	ID         uint    `gorm:"column:id;primaryKey"`
	UserId     uint    `gorm:"index"`
	Volume     float64 `json:"volume" gorm:"column:volume"`
	Price      float64 `json:"price" gorm:"column:price"`
	StationId  uint    `json:"station_id" gorm:"column:station_id"`
	Version    string  `json:"version" gorm:"column:version"`
	Mileage    float64 `json:"mileage" gorm:"column:mileage"`
	CreateTime int     `json:"create_time" gorm:"column:create_time;autoCreateTime"`
}
