package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	// "path/filepath"
	// "runtime"

	"github.com/Lucasvmarangoni/logella/err"
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
	errs.PanicBool(viper.IsSet(key), fmt.Sprintf("The key %s is not defined", key))
	return viper.GetString(key)
}

func GetEnvInt(strucT, field string) int {
	key := fmt.Sprintf("%s_%s", strucT, field)
	errs.PanicBool(viper.IsSet(key), fmt.Sprintf("The key %s is not defined", key))
	return viper.GetInt(key)
}

func GetEnvBool(strucT, field string) bool {
	key := fmt.Sprintf("%s_%s", strucT, field)
	errs.PanicBool(viper.IsSet(key), fmt.Sprintf("The key %s is not defined", key))
	return viper.GetBool(key)
}
func GetTokenAuth() *jwtauth.JWTAuth {    // ReadSecretString é quando é para secret do docker
	tokenAuth := jwtauth.New("HS256", []byte(ReadSecretString(GetEnvString("jwt", "secret"))), nil)
	return tokenAuth
}

func GetTokenAuthTemp() *jwtauth.JWTAuth {
	tokenAuthTemp := jwtauth.New("HS256", []byte(GetEnvString("jwt", "secret_temp")), nil)
	return tokenAuthTemp
}

func init() {
	viper.AutomaticEnv()
	cfg := &conf{}
	err := viper.Unmarshal(&cfg)
	errs.PanicErr(err, "viper.Unmarshal")
}

func ReadSecretString(secretPath string) string {
	content, err := os.ReadFile(secretPath)
	if err != nil {
		errs.PanicErr(err, "os.ReadFile")
	}
	secret := strings.TrimSpace(string(content))
	return secret
}

func ReadSecretInt(secretPath string) int {
	content, err := os.ReadFile(secretPath)
	if err != nil {
		errs.PanicErr(err, "os.ReadFile")
	}
	secretStr := strings.TrimSpace(string(content))
	secretInt, err := strconv.Atoi(secretStr)
	if err != nil {
		errs.PanicErr(err, "strconv.Atoi")
	}
	return secretInt
}
