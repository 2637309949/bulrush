
/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush LoggerWriter plugin]
 */

 package plugins

 import (
	"io"
	"os"
	"fmt"
	"time"
	"path"
	"strings"
	"path/filepath"
	"github.com/gin-gonic/gin"
	"github.com/2637309949/bulrush"
	"github.com/2637309949/bulrush/utils"
 )

 // LOGLEVEL -
type LOGLEVEL int
const (
	// SYSLEVEL  -
	SYSLEVEL LOGLEVEL = 0 + iota
	// USERLEVEL -
	USERLEVEL
	// SYSSTROBE -
	SYSSTROBE  = 1024 * 1024 * 5
	// USERSTROBE -
	USERSTROBE = 1024 * 1024 * 5
)

type (
	// LoggerWriter plugin
	LoggerWriter struct {
		bulrush.PNBase
		cfg *bulrush.Config
	}
)


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
func getLogFile(level LOGLEVEL, basePath string) string {
	var filePath string
	levelStr := levelStr(level)
    filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
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
	filePath  = path.Join(basePath, fileName)
	return filePath
}

// createLog -
func createLog(path string) io.Writer {
	f, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0600)
	writer := io.MultiWriter(f, os.Stdout)
	return writer
}

// Plugin for Recovery
func(logger *LoggerWriter) Plugin() bulrush.PNRet {
	return func(cfg *bulrush.Config, router *gin.RouterGroup) {
		logger.cfg = cfg
		router.Use(func(c *gin.Context){
			logsDir  := utils.Some(utils.LeftV(cfg.String("logs")), 	"logs").(string)
			logsDir   = path.Join(".", logsDir)
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
		})
	}
}


// Writer -
// fileName start with "1"
// User level
func(logger *LoggerWriter) Writer(info string) {
	logsDir  := utils.Some(utils.LeftV(logger.cfg.String("logs")), 	"logs").(string)
	logsDir   = path.Join(".", logsDir)
	logPath  := getLogFile(USERLEVEL, logsDir)
	writer   := createLog(logPath)
	out 	 := writer
	start 	 := time.Now()
	fmt.Fprintf(out, "[%v]<-S-> %s\n", start.Format("2006/01/02 15:04:05"), info)
}
