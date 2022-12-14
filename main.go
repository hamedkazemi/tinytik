package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hamedkazemi/tinytik/common"
	"github.com/hamedkazemi/tinytik/docs"
	"github.com/hamedkazemi/tinytik/modules/links"
	"github.com/sirupsen/logrus"
	_ "github.com/swaggo/gin-swagger"
)

func main() {
	logrus.Info("API system loading ...")

	// swagger documentation information, to change and see
	// other configuration see doc.go
	docs.SwaggerInfo.Title = common.Config.App.Name
	docs.SwaggerInfo.Description = common.Config.App.Description
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = common.Config.App.Host
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// set proxy for outgoing requests ( if needed )
	common.SetProxy()

	// use kafka client
	//common.KafkaCon()

	// init database and Migrate if needed
	db := common.GetDB()
	logrus.Info("Migrating db if needed...")
	//Migrate()
	defer db.Close()

	// init gin router
	r := gin.Default()
	r.Use(common.CORS())

	//v1 := r.Group("/api/v1")
	//{
	links.ConfigGinLinksRouter(r)
	//}

	// health check
	healthCheck := r.Group("/api/ping")

	healthCheck.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// swagger docs files , TODO only for staging and development environments, add condition
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(ginSwagger.Handler))

	logrus.Info("APP will be served at: " + common.Config.App.Host)
	_ = r.Run(common.Config.App.Host)
}
