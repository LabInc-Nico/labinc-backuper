package config

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var appConfig AppConfig
var loggerConfig LoggerConfig
var once sync.Once

type AppConfig struct {
	Backup struct {
		Directory         string `mapstructure:"directory" validate:"required,dirpath"`
		MaxFilesToKeep    int    `mapstructure:"max_files_to_keep" validate:"required,min=1"`
		DumpFolderName    string `mapstructure:"dump_folder_name" validate:"required"`
		ArchiveFolderName string `mapstructure:"archive_folder_name" validate:"required"`
	} `mapstructure:"backup" validate:"required"`
	Database struct {
		User     string `mapstructure:"user" validate:"required"`
		Password string `mapstructure:"password" validate:"required"`
	} `mapstructure:"database" validate:"required"`
}
type LoggerConfig struct {
	Logger struct {
		Level      string  `mapstructure:"level" default:"info" validate:"oneof=info debug warn error fatal"`
		Path       string  `mapstructure:"path" validate:"required" default:"logs"`
		MaxSize    float64 `mapstructure:"max_size" validate:"required" default:"10"`
		MaxAge     float64 `mapstructure:"max_age" validate:"required" default:"7"`
		MaxBackups float64 `mapstructure:"max_backups" validate:"required" default:"3"`
	}
}

func initConfig() error {
	once.Do(func() {
		viper.SetConfigName(".labinc-backuper") // name of config file (without extension)
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME/.config/labinc-backuper")
		viper.SetEnvPrefix("LABINC")
		viper.AutomaticEnv()

	})
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %s", err)
	}
	return nil
}

func GetAppConfig() error {
	err := initConfig()
	if err != nil {
		return err
	}
	if err := viper.Unmarshal(&appConfig); err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}
	validate := validator.New()
	if err := validate.Struct(appConfig); err != nil {
		fmt.Println("Validation errors:", err)
		// Print the entire configuration struct for debugging:
		fmt.Printf("debug : %+v\n", appConfig)
		return fmt.Errorf("configuration errors: %s", err)
	}

	return nil
}
func GetLoggerConfig() error {
	err := initConfig()
	if err != nil {
		return err
	}
	// get default values by reflect config
	// ex: Path       string  `mapstructure:"path" validate:"required" default:"logs"` => get "logs"
	fields := reflect.VisibleFields(reflect.TypeOf(LoggerConfig{}.Logger))
	for _, field := range fields {
		viper.SetDefault(fmt.Sprintf("logger.%s", field.Name), field.Tag.Get("default"))
	}

	if err := viper.Unmarshal(&loggerConfig); err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}
	validate := validator.New()
	if err := validate.Struct(loggerConfig); err != nil {
		fmt.Println("Validation errors:", err)
		// Print the entire configuration struct for debugging:
		fmt.Printf("debug : %+v\n", loggerConfig)
		return fmt.Errorf("configuration errors: %s", err)
	}

	return nil
}
