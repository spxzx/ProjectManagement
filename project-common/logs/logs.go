package logs

import (
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

var Log *zap.Logger

type LogConfig struct {
	DebugFileName string `json:"debugFileName"`
	InfoFileName  string `json:"infoFileName"`
	WarnFileName  string `json:"warnFileName"`
	MaxAge        int    `json:"maxAge"`
	MaxSize       int    `json:"maxSize"`
	MaxBackups    int    `json:"maxBackups"`
}

// getLogWriter 调用日志分割存储日志到文件中
func getLogWriter(filename string, maxAge, maxSize, maxBackup int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxAge:     maxAge,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// InitLogger 初始化Logger
func InitLogger(cfg *LogConfig) (err error) {
	writeSyncerDebug := getLogWriter(cfg.DebugFileName, cfg.MaxAge, cfg.MaxSize, cfg.MaxBackups)
	writeSyncerInfo := getLogWriter(cfg.InfoFileName, cfg.MaxAge, cfg.MaxSize, cfg.MaxBackups)
	writeSyncerWarn := getLogWriter(cfg.WarnFileName, cfg.MaxAge, cfg.MaxSize, cfg.MaxBackups)

	encoder := getEncoder()
	// 文件输出
	debugCore := zapcore.NewCore(encoder, writeSyncerDebug, zapcore.DebugLevel)
	infoCore := zapcore.NewCore(encoder, writeSyncerInfo, zapcore.InfoLevel)
	warnCore := zapcore.NewCore(encoder, writeSyncerWarn, zapcore.WarnLevel)

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	// 标准输出
	std := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)
	core := zapcore.NewTee(debugCore, infoCore, warnCore, std)
	Log = zap.New(core, zap.AddCaller())
	// 替换zap包中全局logger实例，后续再其他包中只需要使用zap.L()调用即可
	zap.ReplaceGlobals(Log)
	return
}

func GinLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		ctx.Next()
		cost := time.Since(start)
		Log.Info(path,
			zap.Int("status", ctx.Writer.Status()),
			zap.Duration("cost", cost),
			zap.String("method", ctx.Request.Method),
			zap.String("path", path),
			zap.String("ip", ctx.ClientIP()),
			zap.String("query", query),
			zap.String("project-grpc-agent", ctx.Request.UserAgent()),
			zap.String("errors", ctx.Errors.ByType(gin.ErrorTypePrivate).String()),
		)
	}
}

func GinRecovery(stack bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 检查连接是否断开，因为这并不是真正需要紧急堆栈跟踪的条件。
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(ctx.Request, false)
				if brokenPipe {
					Log.Error(ctx.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					ctx.Error(err.(error))
					ctx.Abort()
					return
				}
				if stack {
					Log.Error("[recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					Log.Error("[recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		ctx.Next()
	}
}
