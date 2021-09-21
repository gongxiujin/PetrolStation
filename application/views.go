package application

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var SAVEPATH = "/opt/static/file/"
var STATICFILE = "/static/img/"

//var STATICFILE = "/Users/gongxiujin/code/PetrolStation/"

var Logger *logrus.Logger
var ctx = context.Background()

// ----------------小程序api-------------
// 获取用户信息
// @Summary 获取用户信息
// @Description 获取用户信息
// @Accept  json
// @Produce  json
// @Param Token header string true "+Q7xeBtwHmvmwhcMU0ZnQZ6N2jboP8wa5z1MIsrfLck="
// @Param longitude query number true "经度"
// @Param latitude query number true "纬读"
// @Success 200 {object} ResponseJson{code=int,msg=string,data=UserRecordRes} "desc"
// @Router /user/profile [get]
func GetUserProfile(o *gin.Context) {
	user, err := getCurrentUser(o)
	var records []PetrolRecordRes
	if err != nil {
		PackJSONRESP(o, 5001, err.Error())
		return
	}
	var count int64
	var pq ProfileQuery
	var pr LastQtrip
	var mileage float64
	DB.Model(&PetrolRecord{}).Where("user_id = ?", user.ID).Count(&count)
	if count > 1 {
		DB.Raw("select user_id, min(mileage) min_mileage, max(mileage) max_mileage, sum(volume) sum_volume, min(create_time) min_create_time, max(create_time) max_create_time from petrol_records where user_id=? group by user_id;", user.ID).Scan(&pq)
		DB.Raw("select max(a.volume)/(max(a.mileage)-min(a.mileage))*100 last_qtrip from (select * from petrol_records where user_id=? order by create_time desc  limit 2) a group by a.user_id;", user.ID).Scan(&pr)
		mileage = pq.MaxMileage - pq.MinMileage
		if pq.MaxMileage == pq.MinMileage && pq.MaxMileage == 0 {
			mileage = 1.0
		}
	} else {
		pq = ProfileQuery{
			UserId:        0,
			MinMileage:    0,
			MaxMileage:    0,
			SumVolume:     0,
			MinCreateTime: 0,
			MaxCreateTime: 0,
		}
		pr = LastQtrip{LastQtrip: 0}
		mileage = 1.0
	}
	DB.Raw("select pr.id, pr.version, pr.volume, pr.price, pr.mileage, s.name, pr.create_time from petrol_records pr left join stations s on pr.station_id=s.id where pr.user_id = ? order by pr.id desc;", user.ID).Scan(&records)
	res := UserRecordRes{
		LastQtrip:        pr.LastQtrip,
		AvgQtrip:         Decimal(pq.SumVolume / mileage * 100),
		AvgTrip:          Decimal((pq.MaxMileage - pq.MinMileage) / float64(GetDiffDaysBySecond(pq.MaxCreateTime, pq.MinCreateTime))),
		RealMileage:      pq.MaxMileage,
		LastMileage:      Decimal(pq.MaxMileage - pq.MinMileage),
		CumulativeDosage: pq.SumVolume,
		Records:          records,
	}
	res.NickName = user.NickName
	res.HeadImage = user.Avator
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "查询成功",
		Data: res,
	})

}

// 更新用户信息
// @Summary 更新用户信息
// @Description 更新用户信息
// @Accept  json
// @Produce  json
// @Param Token header string true "+Q7xeBtwHmvmwhcMU0ZnQZ6N2jboP8wa5z1MIsrfLck="
// @Param longitude query number true "经度"
// @Param latitude query number true "纬读"
// @Param data body UserProfile true "body data"
// @Success 200 {object} ResponseJson{code=int,msg=string} "desc"
// @Router /user/profile [post]
func UpdateUserProfile(o *gin.Context) {
	user, err := getCurrentUser(o)
	if err != nil {
		PackJSONRESP(o, 5001, err.Error())
		return
	}
	var body UserProfile
	if err = o.ShouldBindJSON(&body); err != nil {
		PackJSONRESP(o, 5001, err.Error())
		return
	}
	user.NickName = body.NickName
	user.Avator = body.HeadImage
	DB.Save(user)
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "更新成功",
	})
}


// 小程序登录
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
		longitude := o.DefaultQuery("longitude", "")
		latitude := o.DefaultQuery("latitude", "")
		err = logUserTrack(longitude, latitude, &user)
	}

	token, err := getUserToken(user.ID)
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "登录成功",
		Data: map[string]string{
			"token": token,
		},
	})
}

// 获取每日油价
// @Summary 每日油价
// @Description 获取每日油价
// @Accept  json
// @Produce  json
// @Param Token header string true "+Q7xeBtwHmvmwhcMU0ZnQZ6N2jboP8wa5z1MIsrfLck="
// @Param longitude query number true "经度"
// @Param latitude query number true "未读"
// @Success 200 {object} ResponseJson{code=int,msg=string,data=[]PetrolDaily} "desc"
// @Router /home/daily_petrol [get]
func DailyPetrol(o *gin.Context) {
	//var p PetrolPrice
	var daily []PetrolDaily
	result := Redis.Get(ctx, "dailyPetrol")
	s, _ := result.Result()
	if err := json.Unmarshal([]byte(s), &daily); err != nil {
		fmt.Println(err)
	}
	// DB.Find(&p, "day", time.Now().Format("2016-01-02"))
	// result = Redis.Get(ctx, "DailyInsert")
	//if result.Val() == "1" && p.ID == 0 {
	//	DB.Create(&daily)
	//	Redis.Del(ctx, "DailyInsert")
	//	Logger.Info("delete DailyInsert")
	//}
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "查询成功",
		Data: daily,
	})
}

// 根据坐标获取当前位置
// @Summary 获取当前位置
// @Description 上传当前位置坐标获取位置信息
// @Accept  json
// @Produce  json
// @Param Token header string true "+Q7xeBtwHmvmwhcMU0ZnQZ6N2jboP8wa5z1MIsrfLck="
// @Param longitude query number true "经度"
// @Param latitude query number true "纬读"
// @Success 200 {object} ResponseJson{code=int,msg=string,data=LocalArea} "desc"
// @Router /user/location [get]
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
		Data: LocalArea{
			Country:  country,
			Province: province,
			City:     city,
			District: district,
		},
	})
}

// 附近加油站
// @Summary 附近加油站
// @Description 根据当前位置坐标获取附近加油站
// @Accept  json
// @Produce  json
// @Param Token header string true "+Q7xeBtwHmvmwhcMU0ZnQZ6N2jboP8wa5z1MIsrfLck="
// @Param longitude query number true "经度"
// @Param latitude query number true "纬度"
// @Param data body NearbyReq true "body data"
// @Success 200 {object} ResponseJson{code=int,msg=string,data=[]NearbyStationRes} "desc"
// @Router /discover/nearby [post]
func NeighborStation(o *gin.Context) {
	var body NearbyReq
	var response []NearbyStationRes
	if err := o.ShouldBindJSON(&body); err != nil {
		PackJSONRESP(o, 4001, err.Error())
		return
	}
	if body.Num == "" || body.Area == "" || body.OrderBy == "" {
		PackJSONRESP(o, 4001, "params error")
		return
	}
	user, err := getCurrentUser(o)
	if err != nil {
		PackJSONRESP(o, 5001, "Access denied")
		return
	}
	var stations []Station

	fmt.Println(fmt.Sprintf("Longitude is: %f, Latitude is: %f", user.Longitude, user.Latitude))
	// redis中根据远近取出所有的加油站
	result := Redis.GeoRadius(ctx, "STATION_NEARBY", user.Longitude, user.Latitude, &redis.GeoRadiusQuery{
		WithDist: true,
		Radius:   100,
		Sort:     "ASC",
	})
	locations, err := result.Result()
	if err != nil {
		PackJSONRESP(o, 4001, "query error: "+err.Error())
		return
	}
	var stationIds []string
	stationIdToDis := make(map[string]float64)

	for _, location := range locations {
		_, err := strconv.Atoi(location.Name)
		if err != nil {
			PackJSONRESP(o, 4001, "parse int error")
			return
		}
		stationIds = append(stationIds, location.Name)
		stationIdToDis[location.Name] = location.Dist
	}
	// 根据油号、范围内的加油站id、区县找出所有加油站
	DB.Debug().Model(&Station{}).Where(fmt.Sprintf("id in (%s) and country = '%s' and publish = true", strings.Join(stationIds, ","), body.Area)).Preload("Petrol", "version = ? ", body.Num).Unscoped().Find(&stations)
	for _, station := range stations {
		response = append(response, NearbyStationRes{
			station,
			stationIdToDis[fmt.Sprintf("%d", station.ID)],
		})
	}
	// 排序
	if body.OrderBy == "price" {
		sort.Sort(stationOrderByPrice(response))
	} else if body.OrderBy == "distance" {
		sort.Sort(stationOrderByDistance(response))
	} else if body.OrderBy == "smart" {
		sort.Sort(stationOrderBySmart(response))
	}

	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "查询成功",
		Data: response,
	})
}

// 获取广告列表
// @Summary 获取广告列表
// @Description 获取广告列表
// @Accept  json
// @Produce  json
// @Param Token header string true "+Q7xeBtwHmvmwhcMU0ZnQZ6N2jboP8wa5z1MIsrfLck="
// @Param longitude query number true "经度"
// @Param latitude query number true "纬度"
// @Param data body AdvertisingReq true "body data"
// @Success 200 {object} ResponseJson{code=int,msg=string,data=[]Advertising} "desc"
// @Router /home/advertising [post]
func AdvertisingRequest(o *gin.Context) {
	var body AdvertisingReq
	if err := o.ShouldBind(&body); err != nil {
		PackJSONRESP(o, 4001, "参数错误: "+err.Error())
		return
	}
	var advertising []Advertising
	DB.Where("location = ? and Publish = true and type = ?", body.Location, body.Type).Find(&advertising)

	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "查询成功",
		Data: advertising,
	})
}

// 添加加油记录
// @Summary 添加加油记录
// @Description 添加加油记录
// @Accept  json
// @Produce  json
// @Param Token header string true "+Q7xeBtwHmvmwhcMU0ZnQZ6N2jboP8wa5z1MIsrfLck="
// @Param longitude query number true "经度"
// @Param latitude query number true "纬度"
// @Param data body PetrolRecord true "body data"
// @Success 200 {object} ResponseJson{code=int,msg=string} "desc"
// @Router /user/record [post]
func AddPetrolRecord(o *gin.Context) {
	var addRecord PetrolRecord
	if err := o.ShouldBind(&addRecord); err != nil {
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

// 删除加油记录
// @Summary 删除加油记录
// @Description 删除加油记录
// @Accept  json
// @Produce  json
// @Param Token header string true "+Q7xeBtwHmvmwhcMU0ZnQZ6N2jboP8wa5z1MIsrfLck="
// @Param longitude query number true "经度"
// @Param latitude query number true "纬度"
// @Param data body DeleteRecordReq true "body data"
// @Success 200 {object} ResponseJson{code=int,msg=string} "desc"
// @Router /user/record [delete]
func DeletePetrolRecord(o *gin.Context) {
	var delRecord DeleteRecordReq
	if err := o.ShouldBind(&delRecord); err != nil {
		PackJSONRESP(o, 4004, err.Error())
		return
	}
	_, err := getCurrentUser(o)
	if err != nil {
		PackJSONRESP(o, 5001, "Access denied")
		return
	}
	DB.Where("id = ?", delRecord.Id).Delete(&PetrolRecord{})
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "删除成功",
	})
}

// 获取当前位置
// @Summary 获取所有可展示的区县
// @Description 获取所有可展示的区县列表
// @Accept  json
// @Produce  json
// @Param Token header string true "+Q7xeBtwHmvmwhcMU0ZnQZ6N2jboP8wa5z1MIsrfLck="
// @Param longitude query number true "经度"
// @Param latitude query number true "纬读"
// @Success 200 {object} ResponseJson{code=int,msg=string,data=[]Area} "desc"
// @Router /discover/area [get]
func GetAreaList(o *gin.Context) {
	var area []Area
	DB.Raw("select country from stations group by country;").Scan(&area)
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "查询成功",
		Data: area,
	})
}

// ----------------web页面api-------------
// 用户登录
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

// 加油站列表
// @Summary 获取所有油站列表
// @Description 获取所有油站列表
// @Accept  json
// @Produce  json
// @Param Token header string true "+Q7xeBtwHmvmwhcMU0ZnQZ6N2jboP8wa5z1MIsrfLck="
// @Param longitude query number true "经度"
// @Param latitude query number true "纬读"
// @Param page query number false "页数" default(1)
// @Param limit query number false "每页条数" default(10)
// @Success 200 {object} ResponseJson{code=int,msg=string,data=[]Station} "desc"
// @Router /station [get]
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
	DB.Limit(perPageInt).Offset((pageInt-1)*perPageInt).Preload("Petrol").Find(&stations)
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

func UpdateStation(o *gin.Context) {
	var station Station
	if err := o.ShouldBind(&station); err != nil {
		PackJSONRESP(o, 4002, err.Error())
		return
	}
	if station.ID == 0 {
		DB.Create(&station)
	} else {
		DB.Save(&station)
	}
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "更新成功",
		Data: map[string]int{
			"id": int(station.ID),
		},
	})
}

// 删除加油站
func DeleteStation(o *gin.Context) {
	ID := o.Param("stationId")
	stationId, err := strconv.Atoi(ID)
	if err != nil {
		PackJSONRESP(o, 4002, err.Error())
		return
	}
	DB.Where("id = ?", stationId).Delete(&Station{})
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "删除成功",
	})
}

// 分享详情
// @Summary 获取分享详情
// @Description 获取分享的文案和图片
// @Accept  json
// @Produce  json
// @Param Token header string true "+Q7xeBtwHmvmwhcMU0ZnQZ6N2jboP8wa5z1MIsrfLck="
// @Param longitude query number true "经度"
// @Param latitude query number true "纬读"
// @Success 200 {object} ResponseJson{code=int,msg=string,data=[]ShareInfoRes} "desc"
// @Router /discover/share_info [get]
func ShareInfo(o *gin.Context) {
	msgs := []string{"今日油价", "油站优惠信息", "附近油价信息", "众车主共分享", "油价共享，就近加油", "最新油价信息"}
	images := []string{"/static/img/20210810113859_image.png", "/static/img/20210810113859_image.png", "/static/img/20210810113859_image.png"}
	rand.Seed(time.Now().Unix())
	msg := msgs[rand.Intn(len(msgs))]
	img := images[rand.Intn(len(images))]
	info := ShareInfoRes{
		Msg: msg,
		Img: img,
	}
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "查询成功",
		Data: info,
	})
}

// 添加/修改加油站油价
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

// 删除加油站油价
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

// 获取广告列表
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

// 删除广告
func DeleteAdvertising(o *gin.Context) {
	var advert Advertising
	ID := o.Param("adverId")
	adverId, err := strconv.Atoi(ID)
	if err != nil {
		PackJSONRESP(o, 4002, err.Error())
		return
	}
	DB.Where("id = ?", adverId).First(&advert)
	if advert.Url != "" {
		file_path := strings.Replace(advert.Url, STATICFILE, SAVEPATH, -1)
		if err := os.Remove(file_path); err != nil {
			PackJSONRESP(o, 4002, err.Error())
			return
		}
	}
	DB.Where("id = ?", adverId).Delete(&Advertising{})
	o.JSON(200, ResponseJson{
		Code: 0,
		Msg:  "删除成功",
	})
}

// 更新广告
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

// 上传广告图片
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

// 创建加油站
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
