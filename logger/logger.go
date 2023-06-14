package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"k8sManagerApi/config"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

var lg *zap.Logger

// InitLogger 初始化Logger
func InitLogger(cfg *config.LogConfig) (err error) {
	writeSyncer := getLogWriter(cfg.Filename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge, cfg.Compress)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return err
	}
	core := zapcore.NewCore(encoder, writeSyncer, l) // 只是写入文件

	//consoleSyncer := zapcore.AddSync(os.Stdout)                                                  // 终端输出
	//core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writeSyncer, consoleSyncer), l) // 同时写入文件和终端输出

	lg = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(lg) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	return nil
}

// getLogWriter 获取日志写入器
func getLogWriter(fileName string, maxSize, maxBackup, maxAge int, compress bool) zapcore.WriteSyncer {
	// 拼接日志存放路径，存放在项目的logs目录下，文件名以配置的为准
	logFilePath := fmt.Sprintf("logs/%s", fileName)

	// 使用lumberjack实现日志文件的切割和清理
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logFilePath, // 日志文件路径和名称
		MaxSize:    maxSize,     // 日志文件的最大大小（以MB为单位）
		MaxAge:     maxAge,      // 保留日志文件的最大天数
		MaxBackups: maxBackup,   // 最多保留的旧日志文件的数量
		Compress:   compress,    // 是否开启日志压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

// getEncoder 获取日志编码器
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05") // 时间格式
	encoderConfig.TimeKey = "time"                                                // 时间字段的键名
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder                       // 将日志级别编码为大写形式
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder                 // 将时间间隔编码为秒
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder                       // 编码调用者信息（文件名和行号）为短格式
	return zapcore.NewJSONEncoder(encoderConfig)
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		ctx.Next()

		cost := time.Since(start)
		lg.Info(path,
			zap.Int("status", ctx.Writer.Status()),                                 // HTTP状态码
			zap.String("method", ctx.Request.Method),                               // 请求方法
			zap.String("path", path),                                               // 请求路径
			zap.String("query", query),                                             // 查询参数
			zap.String("ip", ctx.ClientIP()),                                       // 客户端IP地址
			zap.String("user-agent", ctx.Request.UserAgent()),                      // 客户端User-Agent信息
			zap.String("errors", ctx.Errors.ByType(gin.ErrorTypePrivate).String()), // 错误信息
			zap.Duration("cost", cost),                                             // 请求耗时
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(ctx.Request, false)
				if brokenPipe {
					lg.Error(ctx.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)

					ctx.Error(err.(error))
					ctx.Abort()
					return
				}

				if stack {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					lg.Error("[Recovery from panic]",
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