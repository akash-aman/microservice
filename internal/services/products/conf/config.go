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
	"pkg/grpc"
	"pkg/helper"
	http "pkg/http/server"
	"pkg/logger"
	"pkg/otel/conf"
	"pkg/websocket/gobwas"
	"products/app/core/models"
	"runtime"
	"strings"

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
	Service          *models.Service         `mapstructure:"service" validate:"required"`
	Echo             *http.EchoConfig        `mapstructure:"echo" validate:"required"`
	Logger           *logger.LoggerConfig    `mapstructure:"logger" validate:"required"`
	Sql              *db.SQLConfig           `mapstructure:"sql" validate:"required"`
	GraphQL          *gql.GraphQLConfig      `mapstructure:"graphql" validate:"required"`
	Otel             *conf.OtelConfig        `mapstructure:"telemetry" validate:"required"`
	WSConfig         *gobwas.WebSocketConfig `mapstructure:"websocket" validate:"required"`
	GrpcConfig       *grpc.GrpcConfig        `mapstructure:"grpc_server" validate:"required"`
	GrpcClientConfig *grpc.GrpcClientConfig  `mapstructure:"grpc_client" validate:"required"`
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

func InitConfig() (*Config, *http.EchoConfig, *logger.LoggerConfig, *db.SQLConfig, *gql.GraphQLConfig, *conf.OtelConfig, *gobwas.WebSocketConfig, *grpc.GrpcConfig, error) {
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
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			configPath = d
		}
	}

	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.AddConfigPath(configPath)
	viper.SetConfigType("json")

	// Enable environment variable substitution
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Custom function to replace ${VAR} with environment variables
	viper.SetTypeByDefaultValue(true)

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Error reading config file:", err)
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	// Process environment variable substitutions
	configMap := viper.AllSettings()
	processEnvVars(configMap)

	// Convert back to viper
	err := viper.MergeConfigMap(configMap)
	if err != nil {
		log.Println("Error merging processed config:", err)
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	if err := viper.Unmarshal(cnf); err != nil {
		log.Println("Error unmarshalling config file:", err)
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	log.Println("Config loaded successfully from:", configPath)

	return cnf, cnf.Echo, cnf.Logger, cnf.Sql, cnf.GraphQL, cnf.Otel, cnf.WSConfig, cnf.GrpcConfig, nil
}

/**
 * processEnvValue handles a single
 * environment variable substitution
 */
func processEnvValue(value string) string {
	if !strings.HasPrefix(value, "${") || !strings.HasSuffix(value, "}") {
		return value
	}

	envVar := strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")

	if envVar == "HOSTNAME" {
		return helper.GetHostname()
	}

	if envValue := os.Getenv(envVar); envValue != "" {
		return envValue
	}

	return value
}

/**
 * processEnvVars recursively processes
 * environment variables in configuration
 */
func processEnvVars(configMap map[string]interface{}) {
	for key, value := range configMap {
		switch v := value.(type) {
		case string:
			configMap[key] = processEnvValue(v)

		case map[string]interface{}:
			processEnvVars(v)

		case []interface{}:
			for i, item := range v {
				switch typedItem := item.(type) {
				case string:
					v[i] = processEnvValue(typedItem)
				case map[string]interface{}:
					processEnvVars(typedItem)
				}
			}
		}
	}
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
