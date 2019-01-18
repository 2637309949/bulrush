/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush Upload plugin]
 */

package plugins
import (
	"path"
	"net/http"
	"path/filepath"
	"github.com/gin-gonic/gin"
	"github.com/2637309949/bulrush"
	"github.com/2637309949/bulrush/utils"
)
 
// Upload -
type Upload struct {
	bulrush.PNBase
	AssetPath string
	URLPrefix string
}

 // Plugin -
func (upload *Upload) Plugin () bulrush.PNRet {
	return func(router *gin.RouterGroup) {
		router.POST("/upload", func(c *gin.Context) {
			form, err := c.MultipartForm()
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"data": 	nil,
					"errcode": 	500,
					"errmsg": 	err.Error(),
				})
				c.Abort()
			}
			rets := make([] map[string] interface{}, 0)
			for _, files := range form.File {
				for _, file := range files {
					filename := filepath.Base(file.Filename)
					uuid := utils.RandString(32)
					uuidFileName := utils.RandString(32) + string(filename[len(filename)-len(filepath.Ext(filename)):])
					if err := c.SaveUploadedFile(file, path.Join(upload.AssetPath, uuidFileName)); err != nil {
						c.JSON(http.StatusOK, gin.H{
							"data": 	nil,
							"errcode": 	500,
							"errmsg": 	err.Error(),
						})
						return
					}
					ret := map[string] interface{} {
						"uid": 	  uuid,
						"status": "done",
						"name":   filename,
						"url":    upload.URLPrefix + "/" + uuidFileName,
					}
					rets = append(rets, ret)
				}
			}
			c.JSON(http.StatusOK, gin.H{
				"data": 	nil,
				"errcode": 	nil,
				"errmsg": 	rets,
			})
		})
	}
}
