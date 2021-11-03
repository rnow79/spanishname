package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

//configuration struct
type Config struct {
	Port     string
	Database struct {
		Type     string
		Host     string
		Db       string
		Username string
		Password string
	}
	Tls         bool
	Certificate struct {
		Crt string
		Key string
	}
	LogRequests bool
}

var config Config

//fomat name function
func formatName(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.RequestURI())
	if err != nil {
		log.Fatal(err)
	}
	oname, err := url.QueryUnescape(u.RequestURI())
	if err != nil {
		log.Fatal(err)
	}
	oname = strings.Replace(oname, "/", "", -1)
	oname = strings.Replace(oname, "favicon.ico", "", -1)
	name := strings.ToLower(oname)
	name = strings.Trim(name, " ")
	if len(name) == 0 {
		return
	}
	var nname string
	words := strings.Fields(name)
	for _, word := range words {
		subwords := strings.Split(word, "-")
		sep := " "
		for i, subword := range subwords {
			nword, err := ioutil.ReadFile("words/" + subword)
			if i > 0 {
				sep = "-"
			}
			if err != nil {
				nname = strings.Trim(nname+sep+subword, " ")
				continue
			}
			nsword := strings.ReplaceAll(string(nword), "\r", "")
			nsword = strings.ReplaceAll(nsword, "\n", "")
			nname = strings.Trim(nname+sep+nsword, " ")
		}
	}
	nname = strings.Title(nname)
	nname = strings.ReplaceAll(nname, " De ", " de ")
	nname = strings.ReplaceAll(nname, " Del ", " del ")
	nname = strings.ReplaceAll(nname, " La ", " la ")
	nname = strings.ReplaceAll(nname, " Las ", " las ")
	nname = strings.ReplaceAll(nname, " Los ", " los ")
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "%s", nname)
	if config.LogRequests {
		log.Printf("[%s] -> [%s]\n", oname, nname)
	}
}

func main() {
	//load configuration json file
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("config file: ", err)
	}
	defer file.Close()
	//decode json into Config struct
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	//starting application
	log.Println("starting")
	log.Println("configuration loaded")
	http.HandleFunc("/", formatName)

	if config.Tls {
		//https server, ensure to put both crt and key files in root directory
		err = http.ListenAndServeTLS(":"+config.Port, config.Certificate.Crt, config.Certificate.Key, nil)
	} else {
		//http server
		err = http.ListenAndServe(":"+config.Port, nil)
	}
	if err != nil {
		log.Print("listen: ", err)
	}
	log.Println("exiting")
}
