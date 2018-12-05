package bulrush

import (
	"os"
	"fmt"
	"strings"
	"github.com/gin-gonic/gin"
)

// IsDebugging -
func IsDebugging() bool {
	return gin.Mode() == gin.DebugMode
}

// debugPrint -
func debugPrint(format string, values ...interface{}) {
	if IsDebugging() {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(os.Stderr, "[bh-debug] "+format, values...)
	}
}