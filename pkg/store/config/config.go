package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"time"
)

type Config struct {
	JWT     JWT
	Worker  Worker
	Sqlite3 string
}

type Worker struct {
	Concurrency uint
	Redis       Redis
}

type Redis struct {
	MaxActive int
	MaxIdle   int
	Wait      bool
	Address   string
}

type JWT struct {
	PrivateKey *rsa.PrivateKey
	ExpireIn   time.Duration
	Issuer     string
}

func New() *Config {
	return &Config{
		JWT: JWT{
			PrivateKey: readJWTPrivateKey(),
			ExpireIn:   time.Second * 30,
			Issuer:     "boiler",
		},
		Worker: Worker{
			Concurrency: 10,
			Redis: Redis{
				MaxActive: 5,
				MaxIdle:   5,
				Wait:      true,
				Address:   ":6379",
			},
		},
		Sqlite3: "./db.sqlite3",
	}
}

func readJWTPrivateKey() *rsa.PrivateKey {
	var privateKey *rsa.PrivateKey
	if rawPrivateKey, err := ioutil.ReadFile("jwt.pem"); err != nil {
		log.Panic(err)
	} else {
		block, _ := pem.Decode(rawPrivateKey)
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			log.Fatal(err)
		}
	}

	return privateKey
}
