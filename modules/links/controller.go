package links

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hamedkazemi/tinytik/common"
)

func CreateLink(c *gin.Context) {
	common.GenerateId()
}

func GetLink(c *gin.Context) {
	t := c.Param("token")

	fmt.Printf("ids: %v \n", t)
}
