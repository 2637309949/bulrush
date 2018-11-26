package middles

import (
	"os"
	"path"
	"strings"
	"net/http"
	"github.com/gin-gonic/gin"
)

// index html
const index = "index.html"

// ServeFileSystem -
type ServeFileSystem interface {
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
func (l *LocalFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		name := path.Join(l.root, p)
		stats, err := os.Stat(name)
		if err != nil {
			return false
		}
		if stats.IsDir() {
			if !l.indexes {
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

// Serve -
type Serve struct {
	URLPrefix string
	Fs ServeFileSystem
}

// Inject for gin
func (serve *Serve) Inject(injects map[string]interface{}) {
	engine, _ := injects["Engine"].(*gin.Engine);
	fileserver := http.FileServer(serve.Fs)
	if serve.URLPrefix != "" {
		fileserver = http.StripPrefix(serve.URLPrefix, fileserver)
	}
	engine.GET(serve.URLPrefix + "/*any", func(c *gin.Context) {
		if serve.Fs.Exists(serve.URLPrefix, c.Request.URL.Path) {
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	})
}
