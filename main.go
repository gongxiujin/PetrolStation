package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gongxiujin/PetrolStation/application"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/gin-contrib/static"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path"
	"runtime"
	"time"
)

func InitDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("orm.db"), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err.Error())
	}
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	_ = db.AutoMigrate(&application.Users{}, &application.Advertising{}, &application.Station{},
		&application.PetrolDaily{}, &application.PetrolPrice{}, &application.PetrolRecord{}, &application.UserTrack{})
	// Display SQL queries
	db.Logger.LogMode(1)

	return db
}

func initLogger() *logrus.Logger {
	fileName := path.Join("./log", "access.log")
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.Out = src

	logWriter, err := rotatelogs.New(
		fileName+".%Y%m%d.log",
		rotatelogs.WithLinkName(fileName),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	logger.AddHook(lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}))
	return logger
}

func loggerMiddleware() gin.HandlerFunc {
	logger := initLogger()
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUrl := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUrl,
		}).Info()

	}
}

func main() {
	router := gin.New()
	router.Use(loggerMiddleware())
	var defaultCPUNum = 2
	if runtime.NumCPU() > defaultCPUNum {
		runtime.GOMAXPROCS(3)
	} else {
		if runtime.NumCPU() > 1 {
			runtime.GOMAXPROCS(runtime.NumCPU() - 1)
		} else {
			runtime.GOMAXPROCS(1)
		}

	}
	application.Logger = initLogger()
	application.DB = InitDb()
	application.Redis = application.InitRedisClient()
	//router.Static("/", "./dist/static")
	router.Static(application.STATICFILE, application.SAVEPATH)
	router.Use(static.Serve("/", static.LocalFile("./dist/", false)))
	router.GET("/api/auth", application.AuthLogin)
	router.POST("/api/login", application.WebLogin)
	authRoute := router.Group("/api/v1/")
	authRoute.Use(application.AuthenticationToken())
	{
		authRoute.POST("/discover/nearby", application.NeighborStation)
		authRoute.GET("/station", application.StationList)
		authRoute.POST("/discover/upgrade/nearby", application.CreateStation)
		authRoute.PUT("/discover/petrol", application.AddPetrolPrice)
		authRoute.DELETE("/discover/petrol/:priceId", application.DeletePetrolPrice)
		authRoute.GET("/home/daily_petrol", application.DailyPetrol)
		authRoute.GET("/user/profile", application.GetUserProfile)
		authRoute.POST("/user/record", application.AddPetrolRecord)
		authRoute.GET("/advertising", application.GetAdvertising)
		authRoute.DELETE("/advertising/:adverId", application.DeleteAdvertising)
		authRoute.POST("/advertising", application.UpdateAdvertising)
		authRoute.PUT("/upload/advertising", application.UploadAdvertisingPic)
		authRoute.GET("/user/location", application.GetLocation)
	}

	_ = router.Run()
}
