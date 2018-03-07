package main

import(
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
)

type Config struct{
	BindUrl string
	FilePathWithData string
}

func acessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		//fmt.Println("acessLogMiddleware",r.URL.Path)
		start := time.Now()
		next.ServeHTTP(w,r)
		fmt.Printf("[%s] %s, %s %s\n",r.Method,r.RemoteAddr, r.URL.Path,time.Since(start))
	})
}

func MD5(text string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func makeFilePath(rootPath string, filename string) string {
	return path.Join(rootPath, filename[0:1],filename[1:2],filename[2:3],filename[3:4],filename[4:5])
}


func main() {

	jsonBytes,err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Error read config file: ", err)
	}

	config := Config{}
	err = json.Unmarshal(jsonBytes,&config)
	if err != nil {
		log.Fatal("Error parse config file: ", err)
	}

	log.Print("URL: ",config.BindUrl," Path: ", config.FilePathWithData)
	log.Print("Listening....")

	// site multiplexer
	siteMux := http.NewServeMux()

	//main page
	siteMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w,"Hello. I`m simple filestore")
	})

	//save file in storage
	siteMux.HandleFunc("/put", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintln(w,"Only methon POST allowed")
			return
		}

		//parse agruments
		r.ParseForm()  // parse arguments, you have to call this by yourself
		fmt.Println(r.Form)  // print form information in server side

		value,exists := r.Form["body"]
		if !exists {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w,"Body param empty")
			return
		}

		//calc MD5 sum
		body := strings.Join(value,"")
		bodyMd5 := MD5(body)
		fmt.Println("Body",body,bodyMd5)

		//save file on disk - make path 5 level
		pathFile := makeFilePath(config.FilePathWithData,bodyMd5)
		err = os.MkdirAll(pathFile,0777)
		if err != nil{
			w.WriteHeader(http.StatusExpectationFailed)
			fmt.Fprintln(w,"Can`t create file path")
			return
		}
		pathFile = path.Join(pathFile,bodyMd5)

		fl,err = os.Create(pathFile)
		if err!=nil{
			w.WriteHeader(http.StatusExpectationFailed)
			fmt.Fprintln(w,"Can`t write file")
			return
		}

		fmt.Println(pathFile)

		fmt.Fprintln(w,bodyMd5)

	})

	//get file from storage
	siteMux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintln(w,"Only methon GET allowed")
			return
		}
		fmt.Fprintln(w,"get file")

	})

	//get file from storage
	siteMux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintln(w,"Only methon DELETE allowed")
			return
		}
		fmt.Fprintln(w,"delete file")

	})

	siteHandler := acessLogMiddleware(siteMux)

	err = http.ListenAndServe(config.BindUrl, siteHandler) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
