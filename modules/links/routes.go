package links

import "github.com/gin-gonic/gin"

func ConfigGinLinksRouter(router *gin.Engine) {
	// routes
	//router.GET("/", GetLinks)
	router.GET("/:token", GetLink)
	//router.PATCH("/", UpdateLink)
	router.POST("/shorten/:url", CreateLink)
	//router.DELETE("/", DeleteLink)
}
