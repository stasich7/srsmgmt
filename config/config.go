package config

import (
	"embed"
	"os"
	"strconv"
	"sync"
)

//go:embed *.tpl
var embedFS embed.FS

type Config struct {
	HTTPAddr    string
	SRSAddr     string
	GRPCAddr    string
	HLSAddr     string
	RTMPAddr    string
	SRTAddr     string
	LiveTSPath  string
	DbURI       string
	ApiKey      string
	CacheTTL    int
	IsPprof     bool
	IsDebug     bool
	SRSConfPath string
	TplStorage  embed.FS
}

var cfg *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		cfg = readConfig()
	})
	return cfg
}

func readConfig() *Config {
	return &Config{
		HTTPAddr:    fromEnv("HTTP_ADDR", "0.0.0.0:8887").(string),
		SRSAddr:     fromEnv("SRS_ADDR", "").(string),
		GRPCAddr:    fromEnv("GRPC_ADDR", "0.0.0.0:9087").(string),
		HLSAddr:     fromEnv("HLS_ADDR", "").(string),
		RTMPAddr:    fromEnv("RTMP_ADDR", "").(string),
		SRTAddr:     fromEnv("SRT_ADDR", "").(string),
		LiveTSPath:  fromEnv("LIVE_TS_PATH", "/tmp").(string),
		DbURI:       fromEnv("DATABASE_URI", "").(string),
		ApiKey:      fromEnv("APIKEY", "").(string),
		CacheTTL:    fromEnv("CACHE_TTL", 3).(int),
		IsPprof:     fromEnv("PPROF_ENABLED", false).(bool),
		IsDebug:     fromEnv("DEBUG", false).(bool),
		SRSConfPath: fromEnv("SRS_CONF_PATH", "/opt/srs/trunk/cfg/hls_transcode.conf").(string),
		TplStorage:  embedFS,
	}
}

func fromEnv(key string, defaultVal interface{}) interface{} {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultVal
	}

	switch defaultVal.(type) {
	case bool:
		return (value == "true")
	case int:
		res, _ := strconv.ParseInt(value, 10, 64)
		return int(res)
	default:
		return value
	}
}
