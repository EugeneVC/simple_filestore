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

	v := url.Values{}
	v.Add("body",str)

	resp, _ := http.Post("http://" + config.BindUrl + "/put","application/x-www-form-urlencoded", strings.NewReader(v.Encode()))

	bodyText, err := ioutil.ReadAll(resp.Body)
	fmt.Print(string(bodyText))
}
