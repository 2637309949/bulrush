package plugins

import (
	"errors"
	"time"
	"regexp"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
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

// AuthenHandle auth user info
type AuthenHandle func(c *gin.Context) (interface{}, error)

// Identify authentication interface
type Identify struct {
	ExpiresIn	int
	Routes  	RoutesGroup
	Auth 		AuthenHandle
	Tokens  	TokensGroup
	IgnoreURLs 	[]interface{}
}

// Authorization userinfo
func Authorization (iden *Identify, authData interface{}) map[string]interface{} {
	if authData != nil {
		data := map[string]interface{} {
			"accessToken": 		utils.RandString(32),
			"refreshToken":   	utils.RandString(32),
			"expiresIn": 		iden.ExpiresIn,
			"created": 			now.New(time.Now()).Unix(),
			"updated": 			now.New(time.Now()).Unix(),
			"extra": 			authData,
		}
		iden.Tokens.Save(data)
		return data
	}
	return nil
}

// Authentication userinfo
func Authentication (iden *Identify, accessToken string) (map[string]interface{}, error) {
	verifyToken := iden.Tokens.Find(accessToken, "")
	now 		:= time.Now().Unix()
	if verifyToken == nil {
		return nil, errors.New("auth token failed, token may not exist")
	}
	expiresIn, _ := verifyToken["expiresIn"]
	created, _ 	 := verifyToken["created"]
	if (expiresIn.(float64) + created.(float64)) < float64(now){
		return nil, errors.New("auth token failed, token may be overdue")
	}
	return verifyToken, nil
}

// obtainToken token
func (iden *Identify) obtainToken(c *gin.Context) {
	authData, err := iden.Auth(c)
	if authData != nil {
		data := Authorization(iden, authData)
		c.JSON(http.StatusOK, gin.H{
			"data": 	data, 
			"errcode": 	nil,
			"errmsg":	nil,
		})
		iden.Tokens.Save(data)
		c.Abort()
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": 	nil,
			"errcode": 	500,
			"errmsg": 	err.Error(),
		})
		c.Abort()
	}
}

// revokeToken token
func (iden *Identify) revokeToken(c *gin.Context) {
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
	if accessToken != "" {
		result := iden.Tokens.Revoke(accessToken)
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
}

// refleshToken token
func (iden *Identify) refleshToken(c *gin.Context) {
	var refreshToken string
	queryToken  := c.Query("accessToken")
	formToken   := c.PostForm("accessToken")
	headerToken := c.Request.Header.Get("Authorization")
	if queryToken != "" {
		refreshToken = queryToken
	} else if formToken != "" {
		refreshToken = formToken
	} else if headerToken != "" {
		refreshToken = headerToken
	}
	if refreshToken != "" {
		originToken := iden.Tokens.Find("", refreshToken)
		// reflesh and save
		if originToken != nil {
			accessToken, _ := originToken["accessToken"]
			iden.Tokens.Revoke(accessToken.(string))
			// reflesh info
			originToken["created"] 		= now.New(time.Now()).Unix()
			originToken["updated"] 		= now.New(time.Now()).Unix()
			originToken["accessToken"] 	= utils.RandString(32)
			iden.Tokens.Save(originToken)
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
}

// verifyToken token
func (iden *Identify) verifyToken(c *gin.Context) {
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
	verifyToken, err := Authentication(iden, accessToken)
	if verifyToken != nil {
		c.Set("identify", verifyToken)
		c.Next()
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": 	nil,
			"errcode": 	500,
			"errmsg": 	err.Error(),
		})
		c.Abort()
	}
}

// ResolveToken -
func ResolveToken(c *gin.Context) (string, bool) {
	var token string
	identify, ok :=  c.Get("identify")
	if ok {
		token = identify.(map[string]interface{})["accessToken"].(string)
		return token, ok
	}
	return "", false
}

// ResolveUser -
func ResolveUser(c *gin.Context) (interface{}, bool) {
	var token interface{}
	identify, ok :=  c.Get("identify")
	if ok {
		token = identify.(map[string]interface{})["extra"]
		return token, ok
	}
	return nil, false
}

// Inject for gin
func (iden *Identify) Inject(router *gin.RouterGroup) {
	obtainTokenRoute 	:= iden.Routes.ObtainTokenRoute
	revokeTokenRoute 	:= iden.Routes.RevokeTokenRoute
	refleshTokenRoute 	:= iden.Routes.RefleshTokenRoute
	ignoreUrls 			:= iden.IgnoreURLs
	router.POST(obtainTokenRoute,  iden.obtainToken)
	router.POST(revokeTokenRoute,  iden.revokeToken)
	router.POST(refleshTokenRoute, iden.refleshToken)
	router.Use(func(c *gin.Context) {
		reqPath   := c.Request.URL.Path
		igTarget := utils.Find(ignoreUrls, func(regex interface{}) bool {
			r, _ := regexp.Compile(regex.(string))
			return r.MatchString(reqPath)
		})
		if igTarget != nil {
			c.Next()
		} else {
			iden.verifyToken(c)
		}
	})
}
