package config

import (
	"log"

	"github.com/joho/godotenv"
)

var envs = map[string]string{}
var defaults = map[string]string{
	"PORT": "8080",
}

func LoadEnv() {
	var err error
	envs, err = godotenv.Read()
	if err != nil {
		log.Panic(`Error loading .env file:`, err)
	}
}

func GetEnv(key string) string {
	if val, ok := envs[key]; ok {
		return val
	}
	if val, ok := defaults[key]; ok {
		return val
	}
	return ""
}
