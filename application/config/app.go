package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type AppConfig struct {
	AppEnv          string
	AppName         string
	AppPort         string
	AppDebug        bool
	AppAllowOrigins []string
	AppWebUrl       string
}

type Database struct {
	DatabaseName      string
	DriverName        string
	ConnectionString  string
	MaxConnectionOpen int
	MaxConnectionIdle int
	Timezone          string
}

var (
	App = AppConfig{}
	Db  = Database{}
)

func Init() {
	if err := godotenv.Load(".env"); err != nil {
		logrus.WithFields(logrus.Fields{
			"environment": os.Getenv("APP_ENV"),
			"error":       err.Error(),
		}).Fatalln(".env is not loaded properly")
		os.Exit(1)
	}

	App.AppEnv = os.Getenv("APP_ENV")
	App.AppName = os.Getenv("APP_NAME")
	App.AppPort = os.Getenv("APP_PORT")
	App.AppDebug, _ = strconv.ParseBool(os.Getenv("APP_DEBUG"))
	App.AppAllowOrigins = strings.Split(os.Getenv("APP_ALLOW_ORIGINS"), ";")
	App.AppWebUrl = os.Getenv("APP_WEB_URL")

	Db.DatabaseName = os.Getenv("DB_NAME")
	Db.DriverName = os.Getenv("DB_DRIVER_NAME")
	Db.ConnectionString = os.Getenv("DB_CONNECTION_STRING")
	Db.MaxConnectionOpen, _ = strconv.Atoi(os.Getenv("DB_MAX_CONNECTION_OPEN"))
	Db.MaxConnectionIdle, _ = strconv.Atoi(os.Getenv("DB_MAX_CONNECTION_IDLE"))
	Db.Timezone = os.Getenv("DB_TIMEZONE")
}
