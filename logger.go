package bulrush

import (
		  "github.com/2637309949/bulrush/utils"
		  "github.com/gin-gonic/gin"
	gPath "path"
		  "time"
		  "strings"
		  "path/filepath"
		  "fmt"
		  "os"
		  "io"
)

// LOGLEVEL -
type LOGLEVEL int
const (
	// SYSLEVEL  -
	SYSLEVEL LOGLEVEL = 0 + iota
	// USERLEVEL -
	USERLEVEL
)

const (
	// SYSSTROBE -
	SYSSTROBE  = 1024 * 1024 * 5
	// USERSTROBE -
	USERSTROBE = 1024 * 1024 * 5
)

// ensureDir -
func ensureDir(path string) {
	if err := os.MkdirAll(path, os.ModeDir); err != nil {
		fmt.Println("cannot create hidden directory.")
	} 
}

// levelStr -
func levelStr(level LOGLEVEL) string {
	var levelStr string
	switch level {
		case SYSLEVEL:
			levelStr = "0"
		case USERLEVEL:
			levelStr = "1"
	}
	return levelStr
}

// getLogFile -
func getLogFile(level LOGLEVEL, path string) string {
	var filePath string
	levelStr := levelStr(level)
    filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if filePath != "" {
			return nil
		}
		fileName := info.Name()
		fileSize := info.Size()
		levelMatch := strings.HasPrefix(fileName, levelStr)
		sizeMatch  := false
		if level == SYSLEVEL {
			sizeMatch = fileSize < SYSSTROBE
		} else if level == USERLEVEL {
			sizeMatch = fileSize < USERSTROBE
		}
		
		if levelMatch && sizeMatch {
			filePath = path
		}
        return nil
	})
	if filePath != "" {
		return filePath
	}
	// create level log file
	fileName := time.Now().Format("2006.01.02 15.04.05")
	fileName  = fmt.Sprintf("%s_" + fileName + ".log", levelStr)
	filePath  = gPath.Join(path, fileName)
	return filePath
}

// createLog -
func createLog(path string) io.Writer {
	f, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0600)
	writer := io.MultiWriter(f, os.Stdout)
	return writer
}

// LoggerWithWriter -
// fileName start with "0"
// Gin Middles
// System level
func LoggerWithWriter(bulrush *Bulrush) gin.HandlerFunc {
	return func(c *gin.Context) {
		logsDir  := utils.Some(utils.LeftV(bulrush.config.String("logs")), 	"logs").(string)
		logsDir   = gPath.Join(".", logsDir)
		logPath  := getLogFile(SYSLEVEL, logsDir)
		writer   := createLog(logPath)
		out 	 := writer
		start 	 := time.Now()
		path 	 := c.Request.URL.Path
		raw 	 := c.Request.URL.RawQuery
		c.Next()
		end 	 := time.Now()
		latency  := float64(end.Sub(start) / time.Millisecond)
		clientIP := c.ClientIP()
		method   := c.Request.Method
		if raw != "" {
			path = path + "?" + raw
		}
		fmt.Fprintf(out, "[%v]<-E-> %.2fms %s %6s %s\n", end.Format("2006/01/02 15:04:05"), latency, clientIP, method, path)
	}
}

// LoggerWrap -
// fileName start with "1"
// User level
func LoggerWrap(wc *WellCfg) func(string){
	return func(info string) {
		logsDir  := utils.Some(utils.LeftV(wc.String("logs")), 	"logs").(string)
		logsDir   = gPath.Join(".", logsDir)
		logPath  := getLogFile(USERLEVEL, logsDir)
		writer   := createLog(logPath)
		out 	 := writer
		start 	 := time.Now()
		fmt.Fprintf(out, "[%v]<-S-> %s\n", start.Format("2006/01/02 15:04:05"), info)
	}		
}