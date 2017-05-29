package main

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tylerb/graceful"

	"io/ioutil"

	"github.com/prosline/jobco/application"
	"github.com/prosline/jobco/libunix"
	"github.com/prosline/jobco/models"
)

func init() {
	gob.Register(&models.UserRow{})
}

func newConfig(secretKey string) (*viper.Viper, error) {
	u, err := libunix.CurrentUser()
	if err != nil {
		return nil, err
	}

	c := viper.New()
	c.SetDefault("dsn", fmt.Sprintf("postgres://%v@localhost:5432/jobco?sslmode=disable", u))
	c.SetDefault("cookie_secret", secretKey)
	c.SetDefault("http_addr", ":8888")
	c.SetDefault("http_cert_file", "")
	c.SetDefault("http_key_file", "")
	c.SetDefault("http_drain_interval", "1s")

	c.AutomaticEnv()

	return c, nil
}
func getKey() (string, error) {
	// create a text file named secret.key and include a secret key to it!
	key, err := ioutil.ReadFile("secret.key")
	if err != nil {
		logrus.Println("Secret key not available....!")
		return string(key), err
	}
	return string(key), nil
}
func main() {
	sk, err := getKey()

	if err != nil {
		logrus.Println("Secret key not assigned...!")
	}

	config, err := newConfig(sk)
	if err != nil {
		logrus.Fatal(err)
	}

	app, err := application.New(config)
	if err != nil {
		logrus.Fatal(err)
	}

	middle, err := app.MiddlewareStruct()
	if err != nil {
		logrus.Fatal(err)
	}

	serverAddress := config.Get("http_addr").(string)

	certFile := config.Get("http_cert_file").(string)
	keyFile := config.Get("http_key_file").(string)
	drainIntervalString := config.Get("http_drain_interval").(string)

	drainInterval, err := time.ParseDuration(drainIntervalString)
	if err != nil {
		logrus.Fatal(err)
	}

	srv := &graceful.Server{
		Timeout: drainInterval,
		Server:  &http.Server{Addr: serverAddress, Handler: middle},
	}

	logrus.Infoln("Running HTTP server on " + serverAddress)

	if certFile != "" && keyFile != "" {
		err = srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		err = srv.ListenAndServe()
	}

	if err != nil {
		logrus.Fatal(err)
	}
}
