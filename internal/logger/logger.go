package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/labinc-nico/labinc-backuper/internal/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var once sync.Once
var logger *zap.Logger

// Encodeur personnalisé pour masquer les informations sensibles

func getLogName() string {
	return filepath.Join(viper.GetString("logger.path"), "labinc-backuper.log")
}

func initLogger() *zap.Logger {

	err := config.GetLoggerConfig()
	if err != nil {
		panic(err)
	}
	stdout := zapcore.AddSync(os.Stdout)
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   getLogName(),
		MaxSize:    viper.GetInt("logger.max_size"), // megabytes
		MaxBackups: viper.GetInt("logger.max_backups"),
		MaxAge:     viper.GetInt("logger.max_age"), // days
	})
	logLevel, err := zapcore.ParseLevel(viper.GetString("logger.level"))
	if err != nil {
		fmt.Printf("Niveau de log invalide: %s\n", err)
	}
	level := zap.NewAtomicLevelAt(logLevel)
	productionCfg := zap.NewProductionEncoderConfig()

	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
	return zap.New(core).With(zap.String("uuid", uuid.New().String()))
}

func GetLogger() *zap.SugaredLogger {

	once.Do(func() {
		logger = initLogger()

	})
	return logger.Sugar()
}
