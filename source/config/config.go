package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	// "path/filepath"
	// "runtime"

	errors "github.com/Lucasvmarangoni/logella/err"
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

type conf struct {
	database
	rabbitMQ
	// grpc
	storage
	jwt
	authz
	password
	concurrency
}

func GetEnvString(strucT, field string) string {
	key := fmt.Sprintf("%s_%s", strucT, field)
	errors.PanicBool(viper.IsSet(key), fmt.Sprintf("The key %s is not defined", key))
	return viper.GetString(key)
}

func GetEnvInt(strucT, field string) int {
	key := fmt.Sprintf("%s_%s", strucT, field)
	errors.PanicBool(viper.IsSet(key), fmt.Sprintf("The key %s is not defined", key))
	return viper.GetInt(key)
}

func GetEnvBool(strucT, field string) bool {
	key := fmt.Sprintf("%s_%s", strucT, field)
	errors.PanicBool(viper.IsSet(key), fmt.Sprintf("The key %s is not defined", key))
	return viper.GetBool(key)
}

func GetTokenAuth() *jwtauth.JWTAuth {
	tokenAuth := jwtauth.New("HS256", []byte(ReadSecretString(GetEnvString("jwt", "secret"))), nil)
	return tokenAuth
}

func init() {
	viper.AutomaticEnv()
	cfg := &conf{}
	err := viper.Unmarshal(&cfg)
	errors.PanicErr(err, "viper.Unmarshal")
}

func ReadSecretString(secretPath string) (string) {
	content, err := os.ReadFile(secretPath)
	if err != nil {
		errors.PanicErr(err, "os.ReadFile")
	}
	secret := strings.TrimSpace(string(content))
	return secret
}

func ReadSecretInt(secretPath string) (int) {
	content, err := os.ReadFile(secretPath)
	if err != nil {
		errors.PanicErr(err, "os.ReadFile")
	}
	secretStr := strings.TrimSpace(string(content))
	secretInt, err := strconv.Atoi(secretStr)
	if err != nil {
		errors.PanicErr(err, "strconv.Atoi")
	}
	return secretInt
}