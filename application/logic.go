package application

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
	"io"
	"net/http"
	"strconv"
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

func validateAuth(token string) (user Users, err error) {
	return getUserByToken(token)
}

func AuthenticationToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		Token := context.Request.Header.Get("Token")
		if Token == "" {
			context.AbortWithStatusJSON(http.StatusOK, map[string]interface{}{
				"errno":   5001,
				"result":  struct{}{},
				"message": "Access denied",
			})
			return
		}
		user, err := validateAuth(Token)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusOK, map[string]interface{}{
				"errno":   5001,
				"result":  struct{}{},
				"message": "Access denied",
			})
			return
		}
		if user.ID == 0 {
			context.AbortWithStatusJSON(http.StatusOK, map[string]interface{}{
				"errno":   5001,
				"result":  struct{}{},
				"message": "Access denied",
			})
			return
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
	DB.Table("stations").Find(&stationCache)
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
	userId, err := strconv.Atoi(string(origData))
	if err != nil {
		return
	}
	DB.Model(Users{}).Where("id=?", userId).First(&user)
	return user, nil
}

//func getUserByToken(token string) (user Users, err error) {
//	key := []byte("561d1c4ce26e42189fdecafa15295c73")
//
//	ciphertext, _ := hex.DecodeString(token)
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		return
//	}
//
//	if len(ciphertext) < aes.BlockSize {
//		return
//	}
//	iv := ciphertext[:aes.BlockSize]
//	ciphertext = ciphertext[aes.BlockSize:]
//
//	if len(ciphertext)%aes.BlockSize != 0 {
//		return
//	}
//
//	mode := cipher.NewCBCDecrypter(block, iv)
//
//	mode.CryptBlocks(ciphertext, ciphertext)
//	userId, err := strconv.Atoi(string(ciphertext))
//	if err != nil {
//		return
//	}
//	DB.Model(Users{}).Where("id=?", userId).First(&user)
//	return user, nil
//}
