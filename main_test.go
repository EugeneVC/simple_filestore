package main

import (
	. "testing"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"strings"
	"net/url"
	"fmt"
)

func TestFileStore(t *T) {

	//read config
	jsonBytes,err := ioutil.ReadFile("config.json")
	if err != nil {
		t.Error("Error read config file: ", err)
	}

	config := Config{}
	err = json.Unmarshal(jsonBytes,&config)
	if err != nil {
		t.Error("Error parse config file: ", err)
	}

	//file body
	str := "This is simple text file body"

	//POST FILE
	v := url.Values{}
	v.Add("body",str)

	resp, err := http.Post("http://" + config.BindUrl + "/put","application/x-www-form-urlencoded", strings.NewReader(v.Encode()))
	if err != nil {
		t.Error("Web server fail",err.Error())
		return
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bodyText))

	//GET FILE
	resp, err = http.Get("http://" + config.BindUrl + "/get")
	if err != nil {
		t.Error("Web server fail",err.Error())
		return
	}

	bodyText, err = ioutil.ReadAll(resp.Body)
	fmt.Println(string(bodyText))
}
