package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"strconv"
)

var logErrors bool

// Config all vars
type Config struct {
	Base struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Dbname   string `json:"dbname"`
		LogSQL   bool   `json:"logsql"`
		LogErr   bool   `json:"logerr"`
	} `json:"base"`
	Web struct {
		Host string `json:"host"`
		Port string `json:"port"`
		Log  bool   `json:"log"`
	} `json:"web"`
}

func getConfig() (Config, error) {
	var c Config
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		errmsg("getConfig ReadFile", err)
		return c, err
	}
	if err = json.Unmarshal(file, &c); err != nil {
		errmsg("getConfig Unmarshal", err)
		return c, err
	}
	logErrors = c.Base.LogErr
	if c.Base.Dbname == "" {
		err = errors.New("Error: empty database name in config")
		errmsg("getConfig", err)
	}
	return c, err
}

func toInt(num string) int64 {
	id, err := strconv.ParseInt(num, 10, 64)
	errchkmsg("toInt", err)
	return id
}

func round(v float64, decimals int) float64 {
	var pow float64 = 1
	for i := 0; i < decimals; i++ {
		pow *= 10
	}
	return float64(int((v*pow)+0.5)) / pow
}

func errmsg(str string, err error) {
	if logErrors {
		log.Println("Error in", str, err)
	}
}

func errchkmsg(str string, err error) {
	if logErrors && err != nil {
		log.Println("Error in", str, err)
	}
}
