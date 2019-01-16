/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush delivery plugin]
 */

package plugins
import (
	"os"
	"path"
	"strings"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/2637309949/bulrush"
)

// index html
const index = "index.html"

// LocalFileSystem -
type LocalFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
}

// LocalFile -
func LocalFile(root string, indexes bool) *LocalFileSystem {
	return &LocalFileSystem{
		FileSystem: gin.Dir(root, indexes),
		root:       root,
		indexes:    indexes,
	}
}

// Exists detect the presence of files
func (local *LocalFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		name := path.Join(local.root, p)
		stats, err := os.Stat(name)
		if err != nil {
			return false
		}
		if stats.IsDir() {
			if !local.indexes {
				index := path.Join(name, index)
				_, err := os.Stat(index)
				if err != nil {
					return false
				}
			}
		}
		return true
	}
	return false
}

// Delivery -
type Delivery struct {
	bulrush.PNBase
	Path string
	URLPrefix string
}

// Plugin for gin
func (delivery *Delivery) Plugin() bulrush.PNRet {
	return func(httpProxy *gin.Engine) {
		lf := LocalFile(delivery.Path, false)
		fileserver := http.FileServer(lf)
		if delivery.URLPrefix != "" {
			fileserver = http.StripPrefix(delivery.URLPrefix, fileserver)
		}
		httpProxy.GET(delivery.URLPrefix + "/*any", func(c *gin.Context) {
			if lf.Exists(delivery.URLPrefix, c.Request.URL.Path) {
				fileserver.ServeHTTP(c.Writer, c.Request)
				c.Abort()
			}
		})
	}
}
