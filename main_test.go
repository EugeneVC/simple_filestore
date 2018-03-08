package main

import (
	. "testing"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"strings"
	"net/url"
	_ "fmt"
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
	strMD5 := MD5(str)

	//POST FILE TEST
	v := url.Values{}
	v.Add("body",str)
	res, err := http.Post("http://" + config.BindUrl + "/put", "application/x-www-form-urlencoded",strings.NewReader(v.Encode()))
	if err != nil {
		t.Error(err.Error())
		return
	}

	if status := res.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
		return
	}

	bodyBytes,err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if string(bodyBytes) != strMD5 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			string(bodyBytes), strMD5)
		return
	}

	////GET FILE TEST
	res, err = http.Get("http://" + config.BindUrl + "/get/"+strMD5)
	if err != nil {
		t.Error(err.Error())
		return
	}


	if status := res.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
		return
	}

	bodyBytes,err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if string(bodyBytes) != str {
		t.Errorf("handler returned unexpected body: got %v want %v",
			string(bodyBytes), str)
		return
	}

	//DELETE FILE
	client := &http.Client{}

	req,err := http.NewRequest("DELETE", "http://" + config.BindUrl + "/delete/"+strMD5, nil)
	if err != nil {
		t.Error(err.Error())
		return
	}

	// Fetch Request
	res, err = client.Do(req)
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer res.Body.Close()

	if status := res.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
		return
	}
}
