package config

import (
	"flag"
	"log"
	"os"
	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Addr string 
}

// validating env is present via env-required:"true
type Config struct {
	Env string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	//above line mean if coming from yaml the key is env and if coming from env file its ENV
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HttpServer  `yaml:"http_server"`
}

func Mustload() *Config{ //return config struct type
	var configPath string

	configPath = os.Getenv("CONFIG_PATH") // if coming from env variable file , we have getenv method in os

	// First tries to load config file path from environment variable CONFIG_PATH.
	// If not set in env, looks for a command-line flag
	if configPath == "" {
		flags := flag.String("config", "", "path to configuration file")

		flag.Parse()

		configPath = *flags //we are dereferencing here
	//To get the value stored at the address that a pointer points to, you use * (dereferencing operator)
	}

	if configPath == "" {//If still empty → crash.
		log.Fatal("Config path is not set") 
	}

	//🐾 Checks whether file exists, if not → crash.
	if _, err := os.Stat(configPath); os.IsNotExist(err) { //passing err is of IsNotExist
		log.Fatalf("config file doesnt exist: %s", configPath)
	}


	//now if everything works now 
	// Uses cleanenv to read config from the YAML file into struct.
	// Also automatically validates required fields.
	// Also supports env overrides.
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg) //this will return error
	if err != nil {
		log.Fatalf("cannot read config file: %s", err.Error())
	}

// 	What it does : 👆

// 	Opens the YAML (or JSON, TOML) config file at configPath.
// 	1.Reads it.
// 	2.Parses it into the struct you give (&cfg).
// 	3.Fills struct fields using YAML values and/or environment variables.
// 	4.Checks for required fields (env-required:"true") and applies defaults if defined (env-default:"...").
// 	So basically, it maps your config file + env variables → into Go struct.

	return &cfg
}
