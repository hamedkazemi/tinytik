package common

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"math/rand"
	"sync"
	"time"
)

var instance *gorm.DB
var onceDb sync.Once

// GormLogger is a custom logger for Gorm, making it use logrus.
type GormLogger struct{}

// Print handles log events from Gorm for the custom logrus.
func (*GormLogger) Print(v ...interface{}) {
	switch v[0] {
	case "sql":
		logrus.WithFields(
			logrus.Fields{
				"module":  "gorm",
				"type":    "sql",
				"rows":    v[5],
				"src_ref": v[1],
				"values":  v[4],
			},
		).Debug(v[3])
	case "log":
		logrus.WithFields(logrus.Fields{"module": "gorm", "type": "log"}).Print(v[2])
	}
}

func GetDB() *gorm.DB {
	onceDb.Do(func() {
		// random seed
		rand.Seed(time.Now().UnixNano())
		dbConfig := Config.Database.User + ":" + Config.Database.Password + "@tcp(" + Config.Database.Server + ":" + Config.Database.Port + ")/" + Config.Database.Database
		var dbConnection *gorm.DB
		var connectionError error
		err := retry(3, time.Second*20, func() error {
			dbConnection, connectionError = gorm.Open("mysql", dbConfig+"?charset=utf8&parseTime=True&loc=Local")
			if connectionError != nil {
				logrus.WithFields(logrus.Fields{
					"Connection Error": connectionError,
				}).Warn(dbConfig + " Failed to connect, trying again in 15 sec ...")
				return connectionError
			}
			return nil
		})

		//dbConnection.SetLogger(&GormLogger{}) // setting our custom logger
		dbConnection.LogMode(Config.Database.Debug)
		dbConnection.Set("gorm:auto_preload", true)
		if err != nil {
			logrus.Fatal(err)
		}
		instance = dbConnection
		//defer dbConnection.Close() // closed in main.app
	})
	return instance
}
