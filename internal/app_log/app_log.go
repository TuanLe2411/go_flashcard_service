package app_log

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LOG_KQI_TYPE string

const (
	FUNCTION LOG_KQI_TYPE = "FUNCTION"
	API      LOG_KQI_TYPE = "API"
	DATABASE LOG_KQI_TYPE = "DATABASE"
)

type LOG_KQI_HTTP_METHOD string

type KQI struct {
	TrackingId   string       `json:"trackingId"`
	LogType      LOG_KQI_TYPE `json:"logType"`
	HttpMethod   string       `json:"httpMethod"`
	HttpPath     string       `json:"httpPath"`
	FunctionName string       `json:"functionName"`
	OriginalName string       `json:"originalName"`
	ServiceName  string       `json:"serviceName"`
	Description  string       `json:"description"`
	StartTme     time.Time    `json:"startTime"`
	EndTime      time.Time    `json:"endTime"`
	DurationMs   int64        `json:"durationMs"`
	IsError      bool         `json:"isError"`
	ErrorMessage string       `json:"errorMessage"`
	ResponseCode int          `json:"responseCode"`
}

func InitLogger() {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./log/app.log",
		MaxSize:    50,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}

	writer := zerolog.ConsoleWriter{
		Out:        lumberjackLogger,
		TimeFormat: time.RFC3339,
		NoColor:    true,
		FormatLevel: func(i any) string {
			return fmt.Sprintf("[%s]", i)
		},
		FormatCaller: func(i any) string {
			if caller, ok := i.(string); ok {
				parts := strings.Split(caller, ":")
				if len(parts) >= 2 {
					file := parts[0]
					line := parts[1]
					fileName := filepath.Base(file)
					return fmt.Sprintf("[%s:%s]", fileName, line)
				}
				return fmt.Sprintf("[%s]", caller)
			}
			return fmt.Sprintf("[%v]", i)
		},
		FormatFieldName: func(i any) string {
			return fmt.Sprintf("%s=", i)
		},
	}

	logLevel := os.Getenv("LOG_LEVEL")

	switch logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Logger = zerolog.New(writer).With().Caller().Timestamp().Logger()
}

func LogKQI(l KQI) {
	detail, _ := json.Marshal(l)
	log.Info().RawJSON("KQI", detail).Msg("")
}
