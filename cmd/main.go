package main

import (
	"fmt"
	"os"
	"psj/microservice"
	"psj/microservice/keys"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {

	if err := run(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func run() (err error) {

	if err = godotenv.Load(); err != nil {
		return
	}

	logLevelString := os.Getenv(keys.LogLevel)
	var logLevel logrus.Level
	if logLevel, err = logrus.ParseLevel(logLevelString); err != nil {
		return
	}

	log := logrus.New()
	log.SetLevel(logLevel)

	var dsn string
	if dsn, err = buildDSN(); err != nil {
		return
	}

	var srv *microservice.Server
	if srv, err = microservice.NewServer(dsn, log); err != nil {
		return err
	}

	listenAddress := os.Getenv(keys.ListenAddress)
	if err = srv.Listen(listenAddress); err != nil {
		return
	}

	return
}

func buildDSN() (DSN string, err error) {
	var user, pass, server, port, name string
	if user, err = valFor(keys.DBUser); err != nil {
		return "", err
	}
	if pass, err = valFor(keys.DBPass); err != nil {
		return "", err
	}
	if server, err = valFor(keys.DBServer); err != nil {
		return "", err
	}
	if port, err = valFor(keys.DBPort); err != nil {
		return "", err
	}
	if name, err = valFor(keys.DBName); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", user, pass, server, port, name), nil
}

func valFor(key string) (string, error) {
	var val string
	if val = os.Getenv(key); val == "" {
		return "", fmt.Errorf("key '%s' not found", key)
	}
	return val, nil
}
