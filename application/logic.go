package application

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const BlockSize = 32
const Key = "561d1c4ce26e42189fdecafa15295c73"

type GdApiRes struct {
	Status    string `json:"status"`
	Regeocode struct {
		AddressComponent struct {
			City     string `json:"city"`
			Province string `json:"province"`
			District string `json:"district"`
			Country  string `json:"country"`
		} `json:"addressComponent"`
	} `json:"regeocode"`
	Infocode string `json:"infocode"`
	Info     string `json:"info"`
}

var STATIONKEY = "STATION_NEARBY"

type GDApiError struct{}

func (m *GDApiError) Error() string {
	return "gd api interface error"
}

type CustomError struct {
	Err string
}

func (c *CustomError) Error() string {
	return c.Err
}

func (c *CustomError) setError(e string) {
	c.Err = e
}

func PackJSONRESP(o *gin.Context, code int, msg string) {
	Logger.Error(fmt.Sprintf("origin ip: %s, url: %s, method: %s, error: %s", o.ClientIP(), o.Request.RequestURI, o.Request.Method, msg))
	o.AbortWithStatusJSON(http.StatusOK, ResponseJson{
		Code: code,
		Msg:  msg,
	})
}

func validateAuth(token string) (user Users, err error) {
	return getUserByToken(token)
}

func logUserTrack(longitudeStr, latitudeStr string, user *Users) error {
	if longitudeStr == "" || latitudeStr == "" {
		return errors.New("longitude or latitude needed!")
	}
	longitude, lngErr := strconv.ParseFloat(longitudeStr, 64)
	latitude, latErr := strconv.ParseFloat(latitudeStr, 64)
	if lngErr != nil || latErr != nil {
		return errors.New("incorrect format： longitude or latitude!")
	}
	user.Longitude = longitude
	user.Latitude = latitude
	DB.Save(user)
	DB.Create(&UserTrack{
		Longitude: longitude,
		Latitude:  latitude,
	})
	return nil
}

func AuthenticationToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		Token := context.Request.Header.Get("Token")
		if Token == "" {
			PackJSONRESP(context, 5001, "Access denied")
			return
		}
		user, err := validateAuth(Token)
		if err != nil {
			PackJSONRESP(context, 5001, "token parse error: "+err.Error())
			return
		}
		if user.ID == 0 {
			PackJSONRESP(context, 5001, "用户不存在")
			return
		}
		if user.Active {
			longitude := context.DefaultQuery("longitude", "")
			latitude := context.DefaultQuery("latitude", "")
			if err = logUserTrack(longitude, latitude, &user); err != nil {
				PackJSONRESP(context, 5001, err.Error())
				return
			}
		}
		userStr, _ := json.Marshal(user)
		context.Set("User", string(userStr))
		context.Next()
	}
}
func writeToDb(r io.Reader) (err error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		Logger.Info("open reader err: ", err.Error())
		return err
	}
	var Stations []Station
	index := 2
	for {
		var tmp []string
		var station Station
		for i := 65; i < 73; i++ {
			cell := f.GetCellValue("Sheet1", fmt.Sprintf("%c%d", i, index))
			if cell == "" {
				goto breakHere
			}
			tmp = append(tmp, cell)
		}
		longitude, err := strconv.ParseFloat(tmp[6], 64)
		if err != nil {
			Logger.Info("longitude turn over float64 err: ", err.Error())
			return err
		}
		latitude, err := strconv.ParseFloat(tmp[7], 64)
		if err != nil {
			Logger.Info("latitude turn over float64 err: ", err.Error())
			return err
		}
		stationStr := fmt.Sprintf(`{"name": "%s", "province": "%s", "city": "%s", "country": "%s", "address": "%s", "phone": "%s", "longitude": %f, "latitude": %f}`, tmp[0], tmp[1], tmp[2], tmp[3], tmp[4], tmp[5], longitude, latitude)
		err = json.Unmarshal([]byte(stationStr), &station)
		if err != nil {
			Logger.Info("json unmarshal error: ", err.Error())
			return err
		}
		Stations = append(Stations, station)
		index += 1
	}
breakHere:
	var stationCache []*redis.GeoLocation
	DB.Create(&Stations)
	Redis.Del(ctx, STATIONKEY)
	//DB.Table("stations").Find(&stationCache)
	DB.Raw("select id name, longitude, latitude from stations;").Scan(&stationCache)
	Redis.GeoAdd(ctx, STATIONKEY, stationCache...)
	return
}

func GetLocationByCoord(longitude, latitude float64) (country, province, city, district string, err error) {
	var locationRes GdApiRes
	res, err := SendRequest(fmt.Sprintf(GDLOCATIONAPI, GDKEY, longitude, latitude))
	if err != nil {
		Logger.Error(err.Error())
		return
	}
	if err = json.Unmarshal(res, &locationRes); err != nil {
		Logger.Error(err.Error())
		return
	}
	if locationRes.Info != "OK" {
		err = &GDApiError{}
		return
	}
	country = locationRes.Regeocode.AddressComponent.Country
	province = locationRes.Regeocode.AddressComponent.Province
	if locationRes.Regeocode.AddressComponent.City == "" {
		city = locationRes.Regeocode.AddressComponent.Province
	} else {
		city = locationRes.Regeocode.AddressComponent.City
	}

	district = locationRes.Regeocode.AddressComponent.District
	return
}

func getUserToken(userId int) (string, error) {
	key := []byte(Key)
	origData := []byte(fmt.Sprintf("%32d", userId))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) //初始向量的长度必须等于块block的长度16字节
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func getUserByToken(token string) (user Users, err error) {
	key := []byte(Key)
	crypted, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	//AES分组长度为128位，所以blockSize=16，单位字节
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) //初始向量的长度必须等于块block的长度16字节
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	userStingId := strings.Replace(string(origData), " ", "", -1)
	userId, err := strconv.Atoi(userStingId)
	if err != nil {
		return
	}
	DB.Model(Users{}).Where("id=?", userId).First(&user)
	return user, nil
}

func getCurrentUser(o *gin.Context) (user *Users, err error) {
	userStr := o.GetString("User")
	if err := json.Unmarshal([]byte(userStr), &user); err != nil {
		return nil, errors.New("access denied")
	}
	return user, nil
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

// sort方法
type stationOrderByPrice []NearbyStationRes

func (a stationOrderByPrice) Len() int      { return len(a) }
func (a stationOrderByPrice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a stationOrderByPrice) Less(i, j int) bool {
	return a[i].Petrol[0].Price < a[j].Petrol[0].Price
}

type stationOrderByDistance []NearbyStationRes

func (a stationOrderByDistance) Len() int      { return len(a) }
func (a stationOrderByDistance) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a stationOrderByDistance) Less(i, j int) bool {
	return a[i].Distance < a[j].Distance
}

type stationOrderBySmart []NearbyStationRes

func (a stationOrderBySmart) Len() int      { return len(a) }
func (a stationOrderBySmart) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a stationOrderBySmart) Less(i, j int) bool {
	return a[i].Petrol[0].Price*0.7+a[i].Distance*0.3 < a[j].Petrol[0].Price*0.7+a[j].Distance*0.3
}
