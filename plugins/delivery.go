package plugins

import (
	"os"
	"path"
	"strings"
	"net/http"
	"github.com/gin-gonic/gin"
)

// index html
const index = "index.html"

// DeliveryFileSystem -
type DeliveryFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}

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
	URLPrefix string
	Fs DeliveryFileSystem
}

// Inject for gin
func (delivery *Delivery) Inject(injects map[string]interface{}) {
	engine, _  := injects["Engine"].(*gin.Engine);
	fileserver := http.FileServer(delivery.Fs)
	if delivery.URLPrefix != "" {
		fileserver = http.StripPrefix(delivery.URLPrefix, fileserver)
	}
	engine.GET(delivery.URLPrefix + "/*any", func(c *gin.Context) {
		if delivery.Fs.Exists(delivery.URLPrefix, c.Request.URL.Path) {
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	})
}
