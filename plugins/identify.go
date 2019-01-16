/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush identify plugin]
 */

package plugins

import (
	"time"
	"regexp"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
	"github.com/2637309949/bulrush"
	"github.com/2637309949/bulrush/utils"
)

// RoutesGroup iden routes
type RoutesGroup struct {
	ObtainTokenRoute string
	RevokeTokenRoute string
	RefleshTokenRoute string
}

// TokensGroup token store
type TokensGroup struct {
	Save 	func(map[string]interface{})
	Find	func(accessToken string, refreshToken string) map[string]interface{}
	Revoke  func(accessToken string) bool
}

// Identify authentication interface
type Identify struct {
	bulrush.PNBase
	Auth 		func(c *gin.Context) (interface{}, error)
	ExpiresIn	int
	Routes  	RoutesGroup
	Tokens  	TokensGroup
	FakeURLs 	[]interface{}
}

// obtainToken token
func (iden *Identify) obtainToken(authData interface{}) interface{} {
	if authData != nil {
		data := map[string]interface{} {
			"accessToken": 		utils.RandString(32),
			"refreshToken":   	utils.RandString(32),
			"expiresIn": 		utils.Some(iden.ExpiresIn, 86400),
			"created": 			now.New(time.Now()).Unix(),
			"updated": 			now.New(time.Now()).Unix(),
			"extra": 			authData,
		}
		iden.Tokens.Save(data)
		return data
	}
	return nil
}

// revokeToken token
func (iden *Identify) revokeToken(token string) bool {
	return iden.Tokens.Revoke(token)
}

// refleshToken token
func (iden *Identify) refleshToken(refreshToken string) interface{} {
	originToken := iden.Tokens.Find("", refreshToken)
	if originToken != nil {
		accessToken, _ := originToken["accessToken"]
		iden.Tokens.Revoke(accessToken.(string))
		originToken["created"] 		= now.New(time.Now()).Unix()
		originToken["updated"] 		= now.New(time.Now()).Unix()
		originToken["accessToken"] 	= utils.RandString(32)
		iden.Tokens.Save(originToken)
		return originToken
	}
	return nil
}

// verifyToken token
func (iden *Identify) verifyToken(token string) bool {
	verifyToken := iden.Tokens.Find(token, "")
	now 		:= time.Now().Unix()
	if verifyToken == nil {
		return false
	}
	expiresIn, _ := verifyToken["expiresIn"]
	created, _ 	 := verifyToken["created"]
	if (expiresIn.(float64) + created.(float64)) < float64(now){
		return false
	}
	return true
}

// Plugin -
func (iden *Identify) Plugin() bulrush.PNRet {
	return func(router *gin.RouterGroup) {
		obtainTokenRoute 	:= utils.Some(iden.Routes.ObtainTokenRoute,   "/obtainToken").(string)
		revokeTokenRoute 	:= utils.Some(iden.Routes.RevokeTokenRoute,   "/revokeToken").(string)
		refleshTokenRoute 	:= utils.Some(iden.Routes.RefleshTokenRoute, "/refleshToken").(string)
		FakeURLs 			:= iden.FakeURLs
		router.Use(func (c *gin.Context) {
			var accessToken string
			queryToken  := c.Query("accessToken")
			formToken   := c.PostForm("accessToken")
			headerToken := c.Request.Header.Get("Authorization")
			if queryToken != "" {
				accessToken = queryToken
			} else if formToken != "" {
				accessToken = formToken
			} else if headerToken != "" {
				accessToken = headerToken
			}
			c.Set("accessToken", accessToken)
			c.Next()
		})
		router.POST(obtainTokenRoute, func(c *gin.Context) {
			authData, err := iden.Auth(c)
			if authData != nil {
				data := iden.obtainToken(authData)
				c.JSON(http.StatusOK, gin.H{
					"data": 	data, 
					"errcode": 	nil,
					"errmsg":	nil,
				})
				c.Abort()
			} else {
				c.JSON(http.StatusOK, gin.H{
					"data": 	nil,
					"errcode": 	500,
					"errmsg": 	err.Error(),
				})
				c.Abort()
			}
		})
		router.POST(revokeTokenRoute, func(c *gin.Context) {
			accessToken, _ := c.Get("accessToken")
			if accessToken.(string) != "" {
				result := iden.revokeToken(accessToken.(string))
				if result {
					c.JSON(http.StatusOK, gin.H{
						"data": 	nil,
						"errcode": 	nil,
						"errmsg": 	nil,
					})
				} else {
					c.JSON(http.StatusOK, gin.H{
						"data": 	nil,
						"errcode": 	500,
						"errmsg": 	"revoke token failed, token may not exist",
					})
				}
			}
			c.Abort()
		})
		router.POST(refleshTokenRoute, func(c *gin.Context) {
			refreshToken, _ := c.Get("accessToken")
			if refreshToken.(string) != "" {
				originToken := iden.refleshToken(refreshToken.(string))
				// reflesh and save
				if originToken != nil {
					c.JSON(http.StatusOK, gin.H{
						"data": 	originToken,
						"errcode": 	nil,
						"errmsg": 	nil,
					})
				} else {
					c.JSON(http.StatusOK, gin.H{
						"data": 	nil, 
						"errcode": 	500, 
						"errmsg": 	"reflesh token failed, token may not exist",
					})
				}
			}
			c.Abort()
		})
		router.Use(func(c *gin.Context) {
			reqPath   := c.Request.URL.Path
			fakeURL := utils.Find(FakeURLs, func(regex interface{}) bool {
				r, _ := regexp.Compile(regex.(string))
				return r.MatchString(reqPath)
			})
			if fakeURL != nil {
				c.Next()
			} else {
				accessToken, _ := c.Get("accessToken")
				ret := iden.verifyToken(accessToken.(string))
				if ret {
					rawToken := iden.Tokens.Find(accessToken.(string), "")
					c.Set("accessData", rawToken["extra"])
					c.Next()
				} else {
					c.JSON(http.StatusOK, gin.H{
						"data": 	nil,
						"errcode": 	500,
						"errmsg":   "invalid token",
					})
					c.Abort()
				}
			}
		})
	}
}
