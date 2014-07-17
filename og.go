package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"log"
	"os"
)

var (
	config Config
	Info *log.Logger
	Debug *log.Logger
	Error   *log.Logger
)

// ----------------------------- Structs and interfaces -----------------------------
type Config struct {
	Url            string          `json:"url"`
	Port           int32           `json:"port"`
	Authentication *Authentication `json:"authentication"`
}

type Authentication struct {
	Credentials []Credential `json:"credentials"`
}

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (this Config) Validate(username, password string) bool {
	for _, credential := range this.Authentication.Credentials {
		if credential.Username == username && credential.Password == password {
			return true
		}

	}
	return false
}

// ----------------------------- Structs and interfaces -----------------------------

func main() {

	Info = log.New(os.Stdout, "INFO: ",log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(os.Stdout, "DEBUG: ",log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR: ",log.Ldate|log.Ltime|log.Lshortfile)
	
	Info.Println("Starging og reverse proxy")
	
	location := flag.String("config", "", "og config file")
	
	flag.Parse()

	if len(*location) == 0 {
		
		Error.Println("You need to specify a config file in order to start the proxy")
		return
	}
	err := loadConfig(*location)
	
	Info.Printf("Listening on port %v \n", fmt.Sprintf("%v", config.Port))
	
	if err != nil {
		Error.Println(err)
		return
	}
	
	local, err := url.Parse(config.Url)
	
	if err != nil {
		Error.Panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(local)
	http.Handle("/", auth(proxy))
	err = http.ListenAndServe(":"+fmt.Sprintf("%v", config.Port), nil)
	if err != nil {
		Error.Panic(err)
	}
}

func loadConfig(location string) error {
	contents, err := ioutil.ReadFile(location)
	if err != nil {
		return fmt.Errorf("Could not open file at location %s", location)
	}
	err = json.Unmarshal(contents, &config)
	if err != nil {
		return err
	}
	return nil
}

func auth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		if len(req.Header.Get("Authorization")) == 0 {
			http.Error(w, "Authentication required", http.StatusForbidden)
			return
		}
		auth := strings.SplitN(req.Header["Authorization"][0], " ", 2)
		if auth == nil || len(auth) != 2 || auth[0] != "Basic" {
			http.Error(w, "bad syntax", http.StatusBadRequest)
			return
		}
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !config.Validate(pair[0], pair[1]) {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}
		ip, _, _ := net.SplitHostPort(req.RemoteAddr)
		req.Header.Add("X-Forwarded-For", ip)
		Info.Printf("Forwarding [%v] : [%v]",req.Method, req.RequestURI)
		handler.ServeHTTP(w, req)
	})

}

func Validate(username, password string) bool {
	if username == "vinicius" && password == "carvalho" {
		return true
	}
	return false
}
