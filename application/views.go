package application

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

var SAVEPATH = "/opt/static/file/"
var STATICFILE = "/static/img/"

//var STATICFILE = "/Users/gongxiujin/code/PetrolStation/"

var Logger *logrus.Logger
var ctx = context.Background()

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
	Num string `json:"num"`
	OrderBy	string `json:"order_by"`
	Area	string `json:"area"`
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

func GetUserProfile(o *gin.Context) {
	var pq struct {
		UserId        int     `json:"user_id"`
		MinMileage    float64 `json:"min_mileage"`
		MaxMileage    float64 `json:"max_mileage"`
		SumVolume     float64 `json:"sum_volume"`
		MinCreateTime int     `json:"min_create_time"`
		MaxCreateTime int     `json:"max_create_time"`
	}
	var pr struct {
		LastQtrip float64 `json:"last_qtrip"`
	}
	var user Users
	userStr := o.GetString("User")
	if err := json.Unmarshal([]byte(userStr), &user); err != nil {
		PackJSONRESP(o, 5001, "Access denied")
		return
	}
	DB.Raw("select user_id, min(mileage) min_mileage, max(mileage) max_mileage, sum(volume) sum_volume, min(create_time) min_create_time, max(create_time) max_create_time from petrol_records where user_id=? group by user_id;", user.ID).Scan(&pq)
	DB.Raw("select max(a.volume)/(max(a.mileage)-min(a.mileage))*100 last_qtrip from (select * from petrol_records where user_id=? order by create_time desc  limit 2) a group by a.user_id;", user.ID).Scan(&pr)
	res := UserProfileRes{
		NickName:         user.NickName,
		HeadImage:        user.Avator,
		LastQtrip:        pr.LastQtrip,
		AvgQtrip:         Decimal(pq.SumVolume / (pq.MaxMileage - pq.MinMileage) * 100),
		AvgTrip:          Decimal((pq.MaxMileage - pq.MinMileage) / float64(GetDiffDaysBySecond(pq.MaxCreateTime, pq.MinCreateTime))),
		RealMileage:      pq.MaxMileage,
		LastMileage:      pq.MaxMileage,
		CumulativeDosage: pq.SumVolume,
	}
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "查询成功",
		Data: res,
	})

}

func AuthLogin(o *gin.Context) {
	var authRes WXAuthRes
	code := o.DefaultQuery("code", "")
	if code == "" {
		PackJSONRESP(o, 5001, "Empty code!")
		return
	}
	res, err := SendRequest(fmt.Sprintf(`https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code`, XCXAppId, XCXAppSecret, code))
	if err != nil {
		PackJSONRESP(o, 5001, err.Error())
		return
	}
	if err = json.Unmarshal(res, &authRes); err != nil {
		PackJSONRESP(o, 5002, err.Error())
		return
	}
	if authRes.Errcode != 0 {
		PackJSONRESP(o, 5003, authRes.Errmsg)
		return
	}
	var user Users
	DB.Model(&Users{}).First(&user, &Users{UserName: authRes.Openid})
	if user == (Users{}) || user.ID == 0 {
		user = Users{
			UserName:  authRes.Openid,
			PassWord:  authRes.SessionKey,
			Avator:    "",
			Phone:     "",
			NickName:  "",
			Active:    true,
			Longitude: 0,
			Latitude:  0,
		}
		DB.Create(&user)
	}

	token, err := getUserToken(int(user.ID))
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "登录成功",
		Data: map[string]string{
			"token": token,
		},
	})
}

func WebLogin(o *gin.Context) {
	var body LoginReq
	if err := o.ShouldBindJSON(&body); err != nil {
		PackJSONRESP(o, 4000, err.Error())
		return
	}
	var user Users
	DB.Model(&Users{}).Where(body).First(&user)
	var userByte []byte
	_ = json.Unmarshal(userByte, &user)
	if user == (Users{}) || user.ID == 0 {
		PackJSONRESP(o, 4001, "用户不存在")
		return
	}
	token, err := getUserToken(int(user.ID))
	if err != nil {
		PackJSONRESP(o, 1, fmt.Sprintf("创建token失败: %s", err.Error()))
		return
	}
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "插入成功",
		Data: map[string]string{
			"token": token,
		},
	})
}

func CreateStation(o *gin.Context) {
	file, err := o.FormFile("file")
	if err != nil {
		PackJSONRESP(o, 4001, fmt.Sprintf("read request file error: %s", err.Error()))
		return
	}
	f, err := file.Open()
	if err != nil {
		PackJSONRESP(o, 4001, fmt.Sprintf("open error: %s", err.Error()))
		return
	}
	defer f.Close()
	err = writeToDb(f)
	if err != nil {
		PackJSONRESP(o, 4001, fmt.Sprintf("Write db error: %s", err.Error()))
		return
	}
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "插入成功",
	})
}

func GetLocation(o *gin.Context) {
	user, userErr := getCurrentUser(o)
	if userErr != nil {
		PackJSONRESP(o, 5001, "Access denied")
		return
	}
	country, province, city, district, err := GetLocationByCoord(user.Longitude, user.Latitude)
	if err != nil {
		PackJSONRESP(o, 4001, err.Error())
		return
	}
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "查询成功",
		Data: map[string]interface{}{
			"country": country,
			"province": province,
			"city":    city,
			"district": district,
		},
	})
}

func NeighborStation(o *gin.Context) {
	var body NearbyReq
	var response []NearbyStationRes
	if err := o.ShouldBindJSON(&body); err!=nil{
		PackJSONRESP(o, 4001, err.Error())
		return
	}
	user, err := getCurrentUser(o)
	if err != nil {
		PackJSONRESP(o, 5001, "Access denied")
		return
	}
	locations, err := Redis.GeoRadius(ctx, "station", user.Longitude, user.Latitude, &redis.GeoRadiusQuery{
		WithDist: true,
		Sort:     "ASC",
	}).Result()
	if err != nil {
		PackJSONRESP(o, 4001, "query error: "+err.Error())
		return
	}
	for _, location := range locations {
		var tmpStation Station
		DB.Model(&Station{}).Where("id = ?", strconv.Atoi(location.Name)).First(&tmpStation)
		response =  append(response, NearbyStationRes{
			tmpStation,
			location.Dist,
		})
	}
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "查询成功",
		Data: response,
	})
}

func StationList(o *gin.Context) {
	var stations []Station
	page := o.DefaultQuery("page", "1")
	perPage := o.DefaultQuery("limit", "10")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		PackJSONRESP(o, 4001, "page error")
		return
	}
	perPageInt, err := strconv.Atoi(perPage)
	if err != nil {
		PackJSONRESP(o, 4001, "page error")
		return
	}
	DB.Limit(perPageInt).Offset(pageInt - 1).Preload("Petrol").Find(&stations)
	var count int64
	DB.Model(&Station{}).Count(&count)
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "查询成功",
		Data: map[string]interface{}{
			"stations": stations,
			"total":    count,
		},
	})
}

func AddPetrolPrice(o *gin.Context) {
	var price PetrolPrice
	if err := o.Bind(&price); err != nil {
		PackJSONRESP(o, 4004, err.Error())
		return
	}
	if price.ID == 0 {
		DB.Create(&price)
	} else {
		DB.Save(&price)
	}
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "插入成功",
		Data: map[string]int{
			"id": int(price.ID),
		},
	})
}

func DeletePetrolPrice(o *gin.Context) {
	ID := o.Param("priceId")
	priceId, err := strconv.Atoi(ID)
	if err != nil {
		PackJSONRESP(o, 4002, err.Error())
		return
	}
	DB.Where("id = ?", priceId).Delete(&PetrolPrice{})
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "删除成功",
	})
}

func DailyPetrol(o *gin.Context) {
	var p PetrolPrice
	var daily []PetrolDaily
	result := Redis.Get(ctx, "dailyPetrol")
	s, _ := result.Result()
	if err := json.Unmarshal([]byte(s), &daily); err != nil {
		fmt.Println(err)
	}
	DB.Find(&p, "day", time.Now().Format("2016-01-02"))
	result = Redis.Get(ctx, "DailyInsert")
	if result.Val() == "1" && p.ID == 0 {
		DB.Create(&daily)
		Redis.Del(ctx, "DailyInsert")
		Logger.Info("delete DailyInsert")
	}
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "查询成功",
		Data: daily,
	})
}

func AddPetrolRecord(o *gin.Context) {
	var addRecord PetrolRecord
	if err := o.Bind(&addRecord); err != nil {
		PackJSONRESP(o, 4004, err.Error())
		return
	}
	user, err := getCurrentUser(o)
	if err != nil {
		PackJSONRESP(o, 5001, "Access denied")
		return
	}
	addRecord.UserId = user.ID
	DB.Create(&addRecord)
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "插入成功",
	})
}

func GetAdvertising(o *gin.Context) {
	var advertising []Advertising
	var total int64
	location := o.DefaultQuery("location", "")
	advType := o.DefaultQuery("type", "")
	page := o.DefaultQuery("page", "1")
	perPage := o.DefaultQuery("limit", "10")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		PackJSONRESP(o, 4001, "page error")
		return
	}
	perPageInt, err := strconv.Atoi(perPage)
	if err != nil {
		PackJSONRESP(o, 4001, "page error")
		return
	}
	table := DB.Model(&Advertising{})
	if location != "" || advType != "" {
		table = table.Where("location = ? and Type = ? and publish = 1", location, advType).Order("create_time desc")
		//DB.Find(&advertising, Advertising{
		//	Location: location,
		//	Type:     advType,
		//	Publish:  true,
		//})
	} else {
		table = table.Order("publish, create_time desc")
	}
	table.Limit(perPageInt).Offset(pageInt - 1).Find(&advertising)
	table.Count(&total)

	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "查询成功",
		Data: map[string]interface{}{
			"advertising": advertising,
			"total":       total,
		},
	})
}

func DeleteAdvertising(o *gin.Context) {
	ID := o.Param("adverId")
	adverId, err := strconv.Atoi(ID)
	if err != nil {
		PackJSONRESP(o, 4002, err.Error())
		return
	}
	DB.Where("id = ?", adverId).Delete(&Advertising{})
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "删除成功",
	})
}

func UpdateAdvertising(o *gin.Context) {
	var advert Advertising
	if err := o.Bind(&advert); err != nil {
		PackJSONRESP(o, 4004, err.Error())
		return
	}
	if advert.ID == 0 {
		DB.Create(&advert)
	} else {
		DB.Save(&advert)
	}
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "更新成功",
		Data: map[string]int{
			"id": int(advert.ID),
		},
	})
}

func UploadAdvertisingPic(o *gin.Context) {
	file, err := o.FormFile("image")
	if err != nil {
		PackJSONRESP(o, 4001, err.Error())
		return
	}
	tm := time.Unix(time.Now().Unix(), 0)
	fileName := fmt.Sprintf("%s_%s", tm.Format("20060102030405"), file.Filename)
	if err = o.SaveUploadedFile(file, fmt.Sprintf("%s%s", SAVEPATH, fileName)); err != nil {
		PackJSONRESP(o, 4001, err.Error())
		return
	}
	saveImg := STATICFILE + fileName
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "插入成功",
		Data: map[string]string{
			"file_name": saveImg,
		},
	})
}

func GetDiffDaysBySecond(t1, t2 int) int {
	time1 := time.Unix(int64(t1), 0)
	time2 := time.Unix(int64(t2), 0)

	return GetDiffDays(time1, time2)
}

func GetDiffDays(t1, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)
	diff := int(t1.Sub(t2).Hours() / 24)
	if diff < 1 {
		return 1
	}
	return diff
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
