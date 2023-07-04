// Package common tools and helper functions
package common

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/gorm"
	"github.com/jxskiss/base62"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/big"
	mrand "math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/julienschmidt/httprouter"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandString A helper function to generate random string
func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mrand.Intn(len(letters))]
	}
	return string(b)
}

// NBSecretPassword Keep this two config private, it should not expose to open source
const NBSecretPassword = "A String Very Very Very Strong!!@##$!@#$"

// GenToken A Util function to generate jwt_token which can be used in the request header
func GenToken(id uint, rule string) string {
	jwtToken := jwt.New(jwt.GetSigningMethod("HS256"))
	// Set some claims
	jwtToken.Claims = jwt.MapClaims{
		"id":   id,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"rule": rule,
	}
	// Sign and get the complete encoded token as a string
	token, _ := jwtToken.SignedString([]byte(NBSecretPassword))
	return token
}

// Error My own Error type that will help return my customized Error info
//
//	{"database": {"hello":"no such table", error: "not_exists"}}
type Error struct {
	Errors map[string]interface{} `json:"errors"`
}

// NewValidatorError To handle the error returned by c.Bind in gin framework
func NewValidatorError(err error) Error {
	res := Error{}
	res.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	for _, v := range errs {
		// can translate each error one at a time.
		//fmt.Println("gg",v.NameNamespace)
		if v.Param() != "" {
			res.Errors[v.Field()] = fmt.Sprintf("{%v: %v}", v.Tag(), v.Param())
		} else {
			res.Errors[v.Field()] = fmt.Sprintf("{key: %v}", v.Tag())
		}

	}
	return res
}

// NewError Warp the error info in a object
func NewError(key string, err error) Error {
	res := Error{}
	res.Errors = make(map[string]interface{})
	res.Errors[key] = err.Error()
	return res
}

// Bind Changed the c.MustBindWith() ->  c.ShouldBindWith().
// I don't want to auto return 400 when error happened.
// origin function is here: https://github.com/gin-gonic/gin/blob/master/context.go
func Bind(c *gin.Context, obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.ShouldBindWith(obj, b)
}

// PrettyPrint Debugging Tool / Pretty Print of Struct
func PrettyPrint(input interface{}) {
	b, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}

// function for retrying task execution, using this function
// you can easily retry task over failures with specified duration.
func retry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}

		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(mrand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2

			time.Sleep(sleep)
			return retry(attempts, 2*sleep, f)
		}
		return err
	}

	return nil
}

type stop struct {
	error
}

// retry end

// ConverHttprouterToGin Convert Http Router to Gin Routing
func ConverHttprouterToGin(f httprouter.Handle) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params httprouter.Params
		_len := len(c.Params)
		if _len == 0 {
			params = nil
		} else {
			params = ((*[1 << 10]httprouter.Param)(unsafe.Pointer(&c.Params[0])))[:_len]
		}

		f(c.Writer, c.Request, params)
	}
}

// ReadInt Read Integer From String
func ReadInt(r *http.Request, param string, v int64) (int64, error) {
	p := r.FormValue(param)
	if p == "" {
		return v, nil
	}

	return strconv.ParseInt(p, 10, 64)
}

// WriteJSON Write Json
func WriteJSON(w http.ResponseWriter, v interface{}) {
	data, _ := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	_, _ = w.Write(data)
}

// ReadJSON DEPRECATED Read Json
func ReadJSON(r *http.Request, v interface{}) error {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(buf, v)
}

// CORS middleware for Gin
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "OPTIONS" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Cache-Control")
			c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
			c.Next()
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Cache-Control")
			c.Header("Content-Type", "application/json")
			c.AbortWithStatus(http.StatusOK)
		}
	}
}

// IsExistInFilters function to check if key exist in filter types
// return index if exists and return -1 if not found
func IsExistInFilters(filters []Filter, key string) int {
	for i, filter := range filters {
		if filter.Field == key {
			return i
		}
	}
	return -1
}

// BuildSearchByTags by using this function we can build search queries by filter types
func BuildSearchByTags(model interface{}, filters []Filter, db *gorm.DB) *gorm.DB {
	t := reflect.TypeOf(model)
	// Iterate over all available fields and read the tag value

	// first let see if it's search request
	searchReq := -1
	for i, filter := range filters {
		if filter.Field == "allKeys" {
			searchReq = i
		}
	}

	for i := 0; i < t.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := t.Field(i)

		// Get the field tag value
		tag := field.Tag.Get("gorm")
		col := ""
		if tag != "" {
			_, _ = fmt.Sscanf(tag, "column:%s", &col)
			if col != "" {
				if strings.ContainsAny(col, ";") {
					col = col[0:strings.Index(col, ";")]
				}
			}
		}
		isSearchable := field.Tag.Get("searchable")

		if searchReq != -1 {
			if isSearchable != "false" {
				if isSearchable == "string" {
					db = db.Or(col + " LIKE '%" + filters[searchReq].Value + "%'")
				}
				if isSearchable == "int" {
					db = db.Or(col + " = '" + filters[searchReq].Value + "'")
				}

			}
		} else {
			isFound := IsExistInFilters(filters, col)
			if isFound != -1 && isSearchable != "false" {
				f := filters[isFound]
				if f.Operator != "LIKE" {
					db = db.Or(f.Field + " " + f.Operator + " '" + f.Value + "'")
				} else {
					db = db.Or(f.Field + " LIKE '%" + f.Value + "%'")
				}
			}
		}

	}
	return db
}

func GetOptions(input GetAllRequest) ([]Filter, string, string, int, int) {

	var searchQuery string
	var orderBy string
	var orderType string

	limit := input.Limit
	offset := input.Offset

	if input.Query != "" {
		searchQuery = input.Query
	}
	if input.OrderBy != "" {
		orderBy = input.OrderBy
	} else {
		orderBy = "created_at"
	}
	if input.OrderType != "" {
		orderType = input.OrderType
	} else {
		orderType = "DESC"
	}

	filters := input.Filters

	if searchQuery != "" {
		filters = append(filters, Filter{Field: "allKeys", Operator: "LIKE", Value: searchQuery})
	}
	return filters, orderBy, orderType, limit, offset
}

func GetMetadata(modelCount int, limit int, offset int, filters []Filter, orderBy string, orderType string) MetaData {
	var metaData MetaData
	metaData.Pagination.Count = modelCount
	metaData.Pagination.Limit = limit
	metaData.Pagination.Offset = offset
	metaData.Filters = filters
	metaData.Order.OrderBy = orderBy
	metaData.Order.OrderType = orderType
	return metaData
}

func SetProxy() {
	getProxy := ""
	if os.Getenv("HTTP_PROXY") != "" {
		getProxy = os.Getenv("HTTP_PROXY")
	} else {
		getProxy = Config.App.Proxy
	}

	proxyUrl, err := url.Parse(getProxy)
	if err != nil {
		logrus.Error("the URL of Proxy is wrong.")
		logrus.Panic("the URL of Proxy is wrong, please check config.toml file.")
	}

	http.DefaultTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		Proxy:                 http.ProxyURL(proxyUrl),
	}
}

// Substr this isn't multi-Unicode-codepoint aware, like specifying skin tone or
//
//	gender of an emoji: https://unicode.org/emoji/charts/full-emoji-modifiers.html
func Substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

func GenerateId() string {
	p, _ := rand.Prime(rand.Reader, 40)
	s := big.NewInt(time.Now().UnixNano()).Int64()
	l := p.Int64()
	hd := []byte(strconv.FormatInt(s+l, 10))
	d := base62.EncodeToString(hd)
	return strings.ToLower(Substr(Reverse(d), len(d)-10, len(d)))
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
