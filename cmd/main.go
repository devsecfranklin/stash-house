package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type AppConfig struct {
	App struct {
		Name string `mapstructure:"name"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"app"`
	NS    string `mapstructure:"namespace"`
	Owner string `mapstructure:"owner"`
}

func main() {
	// Set up viper to read the config.yaml file
	viper.SetConfigName("config") // Config file name without extension
	viper.SetConfigType("yaml")   // Config file type
	viper.AddConfigPath(".")      // Look for the config file in the current directory

	/*
	   AutomaticEnv will check for an environment variable any time a viper.Get request is made.
	   It will apply the following rules.
	       It will check for an environment variable with a name matching the key uppercased and prefixed with the EnvPrefix if set.
	*/
	viper.AutomaticEnv()
	viper.SetEnvPrefix("env")                              // will be uppercased automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // this is useful e.g. want to use . in Get() calls, but environmental variables to use _ delimiters (e.g. app.port -> APP_PORT)

	// Read the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Set up environment variable mappings if necessary
	/*
	   BindEnv takes one or more parameters. The first parameter is the key name, the rest are the name of the environment variables to bind to this key.
	   If more than one are provided, they will take precedence in the specified order. The name of the environment variable is case sensitive.
	   If the ENV variable name is not provided, then Viper will automatically assume that the ENV variable matches the following format: prefix + "_" + the key name in ALL CAPS.
	   When you explicitly provide the ENV variable name (the second parameter), it does not automatically add the prefix.
	       For example if the second parameter is "id", Viper will look for the ENV variable "ID".
	*/
	viper.BindEnv("app.name", "APP_NAME") // Bind the app.name key to the APP_NAME environment variable

	// Get the values, using env variables if present
	appName := viper.GetString("app.name")
	namespace := viper.GetString("namespace") // AutomaticEnv will look for an environment variable called `ENV_NAMESPACE` ( prefix + "_" + key in ALL CAPS)
	appPort := viper.GetInt("app.port")       // AutomaticEnv will look for an environment variable called `ENV_APP_PORT` ( prefix + "_" + key in ALL CAPS with _ delimiters)

	// Output the configuration values
	fmt.Printf("App Name: %s\n", appName)
	fmt.Printf("Namespace: %s\n", namespace)
	fmt.Printf("App Port: %d\n", appPort)

	// Create an instance of AppConfig
	var config AppConfig
	// Unmarshal the config file into the AppConfig struct
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	// Output the configuration values
	fmt.Printf("Config: %v\n", config)
}
