package links

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func CreateLink(c *gin.Context) {

}

func GetLink(c *gin.Context) {
	t := c.Param("token")

	fmt.Printf("ids: %v \n", t)
}
