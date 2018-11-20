package middles

import (
	"regexp"
	"net/http"
	"github.com/gin-gonic/gin"
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
	Save 	func(interface{}) interface{}
	Revoke  func(interface{}) interface{}
	Verify	func(interface{}) interface{}
}

// AuthenHandle auth user info
type AuthenHandle func(c *gin.Context) (interface{}, error)

// Iden authentication interface
type Iden struct {
	ExpiresIn	int
	Routes  	RoutesGroup
	Auth 		AuthenHandle
	Tokens  	TokensGroup
	IgnoreURLs 	[]interface{}
}

func (iden *Iden) obtainToken(c *gin.Context) {
	authData, err := iden.Auth(c)
	if authData != nil {
		data := map[string]interface{}{
			"accessToken": 		utils.RandString(32),
			"refreshToken":   	utils.RandString(32),
			"expiresIn": 		iden.ExpiresIn,
			"created": 			1542634194,
			"updated": 			1542634194,
			"extra": 			authData,
		}
		c.JSON(http.StatusOK, gin.H{
			"data": 	data, 
			"errcode": 	nil, 
			"errmsg":	 nil,
		})
		iden.Tokens.Save(data)
		c.Abort()
	} else {
		c.JSON(http.StatusOK, gin.H{"data": nil, "errcode": 500, "errmsg": err.Error() })
		c.Abort()
	}
}


func (iden *Iden) revokeToken(c *gin.Context) {

}

func (iden *Iden) refleshToken(c *gin.Context) {

}


func (iden *Iden) authToken(c *gin.Context) {
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
	verifyToken := iden.Tokens.Verify(accessToken)
	if verifyToken != nil {
		c.Set("identify", verifyToken)
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"data": nil, "errcode": 500, "errmsg": "authToken failed" })
	}
}


// Middle for gin
func (iden *Iden) Middle(app *bulrush.Bulrush) gin.HandlerFunc {
	var obtainTokenRoute = iden.Routes.ObtainTokenRoute
	var revokeTokenRoute = iden.Routes.RevokeTokenRoute
	var refleshTokenRoute = iden.Routes.RefleshTokenRoute
	var ignoreUrls = iden.IgnoreURLs
	return func(c *gin.Context) {
		reqPath   := c.Request.URL.Path
		reqMethod := c.Request.Method
		igTarget := utils.Find(ignoreUrls, func(regex interface{}) bool{
			r, _ := regexp.Compile(regex.(string))
			return r.MatchString(reqPath)
		})
		otTarget := reqMethod == "POST" && reqPath == obtainTokenRoute
		rtTarget := reqMethod == "POST" && reqPath == revokeTokenRoute
		ftTarget := reqMethod == "POST" && reqPath == refleshTokenRoute
		// obtainToken
		if otTarget {
			iden.obtainToken(c)
		// refleshToken
		} else if ftTarget {
			iden.refleshToken(c)
		// revokeToken
		} else if rtTarget {
			iden.revokeToken(c)
		// ignore spec req
		} else if igTarget != nil {
			c.Next()
		// authToken
		} else {
			iden.authToken(c)
		}
	}
}
