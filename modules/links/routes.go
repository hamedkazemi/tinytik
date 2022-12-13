package links

import "github.com/gin-gonic/gin"

func ConfigGinLinksRouter(router *gin.Engine) {
	// routes
	router.GET("/", Getlinks)
	router.GET("/:id", Getlink)
	router.PATCH("/", Updatelink)
	router.POST("/", Createlink)
	router.DELETE("/", Deletelink)
}
