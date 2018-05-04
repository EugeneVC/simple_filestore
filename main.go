package main

import (
	_ "fmt"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"fmt"
	_ "time"
	"time"
	"strings"
	"crypto/md5"
	"encoding/hex"
	"path"
	"os"
	"github.com/gorilla/mux"
)

const MAX_UPLOAD_FILE_SIZE = 4*1024*1024*1024 // 4Gb

type Config struct {
	BindUrl          string
	FilePathWithData string
}

var config Config

func acessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("acessLogMiddleware",r.URL.Path)
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("[%s] %s, %s %s\n", r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
	})
}

func MD5(text string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func makeFilePath(rootPath string, filename string) string {
	return path.Join(rootPath, filename[0:1], filename[1:2], filename[2:3], filename[3:4], filename[4:5])
}

func postFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Only methon POST allowed")
		return
	}

	maxlength := int64(MAX_UPLOAD_FILE_SIZE)

	//parse agruments
	//fmt.Println("Content-type",r.Header.Get("Content-type"))
	r.Body = http.MaxBytesReader(w, r.Body, maxlength)
	err:= r.ParseForm()       // parse arguments, you have to call this by yourself
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err.Error())
		return
	}
	//fmt.Println(r.Form) // print form information in server side

	value, exists := r.Form["body"]
	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Body param absent")
		return
	}

	//calc MD5 sum
	body := strings.Join(value, "")
	bodyMd5 := MD5(body)
	//fmt.Println("Body", body, bodyMd5)

	//save file on disk - make path 5 level
	pathFile := makeFilePath(config.FilePathWithData, bodyMd5)
	err = os.MkdirAll(pathFile, 0777)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprintln(w, "Can`t create file path: "+ pathFile)
		return
	}
	pathFile = path.Join(pathFile, bodyMd5)

	fl, err := os.Create(pathFile)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprintln(w, "Can`t create file")
		return
	}
	defer fl.Close()

	_,err = fl.Write([]byte(body))
	if err != nil{
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprintln(w, "Can`t write file")
		return
	}

	//fmt.Println("Content-Length",r.Header.Get("Content-Length"),bodyMd5,pathFile)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(bodyMd5))
}

func getFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Only methon GET allowed")
		return
	}

	vars := mux.Vars(r)
	filename,exist := vars["filename"]
	if !exist {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//read file
	pathFile := makeFilePath(config.FilePathWithData, filename)
	pathFile = path.Join(pathFile, filename)

	body,err := ioutil.ReadFile(pathFile)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

func deleteFile(w http.ResponseWriter, r *http.Request){
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Only methon DELETE allowed")
		return
	}

	vars := mux.Vars(r)
	filename,exist := vars["filename"]
	if !exist {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//make path
	pathFile := makeFilePath(config.FilePathWithData, filename)
	pathFile = path.Join(pathFile, filename)

	//check file exits
	if _, err := os.Stat(pathFile); err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	//delete file
	err := os.Remove(pathFile)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {

	jsonBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Error read config file: ", err)
	}

	err = json.Unmarshal(jsonBytes, &config)
	if err != nil {
		log.Fatal("Error parse config file: ", err)
	}

	log.Print("URL: ", config.BindUrl, " Path: ", config.FilePathWithData)
	log.Print("Listening....")

	// site multiplexer
	siteMux := mux.NewRouter()

	//main page
	siteMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello. I`m simple filestore")
	})

	//save file in storage
	siteMux.HandleFunc("/put", postFile)

	//get file from storage
	siteMux.HandleFunc("/get/{filename}", getFile)

	//get file from storage
	siteMux.HandleFunc("/delete/{filename}", deleteFile)

	siteHandler := acessLogMiddleware(siteMux)

	err = http.ListenAndServe(config.BindUrl, siteHandler) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
