package links

import (
	"github.com/gin-gonic/gin"
	"github.com/hamedkazemi/tinytik/common"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	URL string `json:"url"`
}

func PostLink(c *gin.Context) {
	request := Request{}
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	_, err := url.ParseRequestURI(request.URL)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	md5 := common.GetMD5Hash(request.URL)

	cr := common.ConnectRedis()
	defer cr.Close()

	urlExist, err := cr.Do(c, cr.B().Get().Key("TINYTIK_"+md5).Build()).ToString()
	if len(urlExist) > 1 {
		strb := strings.Split(urlExist, "|||")
		response := CreateResponse{URL: strb[0], Hash: strb[1]}
		c.IndentedJSON(http.StatusOK, response)
		return
	}
	hash := common.GenerateId()
	err = cr.Do(c, cr.B().Set().Key("TINYTIK_"+md5).Value(request.URL+"|||"+hash).Nx().Build()).Error()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	link := Links{Hash: hash, URL: request.URL, MD5: md5}
	result := common.GetDB().Create(&link)
	if result.Error != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response := CreateResponse{URL: link.URL, Hash: link.Hash}
	c.IndentedJSON(http.StatusOK, response)
}

func GetLink(c *gin.Context) {
	t := c.Param("token")
	cr := common.ConnectRedis()
	defer cr.Close()
	urlExist, _ := cr.Do(c, cr.B().Get().Key("TINYTIK_"+t).Build()).ToString()
	if len(urlExist) > 1 {
		c.Redirect(http.StatusFound, urlExist)
		return
	}

	var link Links
	common.GetDB().First(&link, "hash = ?", t)
	_ = cr.Do(c, cr.B().Set().Key("TINYTIK_"+t).Value(link.URL).Nx().Build()).Error()
	c.Redirect(http.StatusFound, link.URL)
}
