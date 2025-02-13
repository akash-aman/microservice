package conf

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"pkg/db"
	"pkg/gql"
	http "pkg/http/server"
	"pkg/logger"
	"products/app/core/models"
	"runtime"

	"github.com/spf13/viper"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "path to config file")
}

/**
 * Config - Centralized configuration for all the service present in application
 */
type Config struct {
	Service *models.Service      `mapstructure:"service" validate:"required"`
	Echo    *http.EchoConfig     `mapstructure:"echo" validate:"required"`
	Logger  *logger.LoggerConfig `mapstructure:"logger" validate:"required"`
	Sql     *db.SQLConfig        `mapstructure:"sql" validate:"required"`
	GraphQL *gql.GraphQLConfig   `mapstructure:"graphql" validate:"required"`
}

/**
 * InitConfig - Initializes and loads the configuration for the application.
 *
 * Details:
 * - Determines the environment from the "APP_ENV" environment variable, defaulting to "development" if not set.
 * - The configuration file path can be supplied through a flag or the "CONFIG_PATH" environment variable.
 * - If neither is provided, it uses the directory from where this function is called.
 * - The configuration file is expected to be in JSON format and named "config.<env>.json".
 * - Uses the Viper library to read and unmarshal the configuration into a Config struct.
 * - Returns the loaded Config struct, EchoConfig struct, and an error if any occurs during the process.
 */

func InitConfig() (*Config, *http.EchoConfig, *logger.LoggerConfig, *db.SQLConfig, *gql.GraphQLConfig, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	cnf := &Config{}

	/**
	 * configPath : Supplied through flag.
	 * configPathFromEnv : Supplied through env variable.
	 * dirname() : directory from where this function is called.
	 */
	if configPath == "" {
		configPathFromEnv := os.Getenv("CONFIG_PATH")
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			d, err := CallerDirPath()
			if err != nil {
				log.Println("Error getting current directory:", err)
				return nil, nil, nil, nil, nil, err
			}
			configPath = d
		}
	}

	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.AddConfigPath(configPath)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Error reading config file:", err)
		return nil, nil, nil, nil, nil, err
	}

	if err := viper.Unmarshal(cnf); err != nil {
		log.Println("Error unmarshalling config file:", err)
		return nil, nil, nil, nil, nil, err
	}

	log.Println("Config loaded successfully from:", configPath)

	return cnf, cnf.Echo, cnf.Logger, cnf.Sql, cnf.GraphQL, nil
}

/**
 * CallerFilename returns the filename of the caller.
 * It uses the runtime.Caller function to get the current filename.
 * If it fails to retrieve the filename, it returns an error.
 */
func CallerFilename() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}
	return filename, nil
}

/**
 * CallerDirPath returns the directory path of the caller.
 * It calls CallerFilename to get the current filename and then
 * uses filepath.Dir to get the directory path.
 * If it fails to retrieve the filename, it returns an error.
 */
func CallerDirPath() (string, error) {
	filename, err := CallerFilename()
	if err != nil {
		return "", err
	}
	return filepath.Dir(filename), nil
}
