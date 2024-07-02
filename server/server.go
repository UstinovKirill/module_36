package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	api "module_31/pkg/api"
	"module_31/pkg/rss"
	storage "module_31/pkg/storage"
	db "module_31/pkg/storage/db"
	"net/http"
	"os"
	"time"
)

type server struct {
	db  storage.Interface
	api *api.API
}

type config struct {
	Period  int      `json:"request_period"`
	LinkArr []string `json:"rss"`
}

func main() {

	var srv server

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	db, err := db.New(ctx, "postgres://postgres:rootroot@localhost:5432/aggregator")
	if err != nil {
		log.Fatal(err)
	}

	srv.db = db

	srv.api = api.New(srv.db)

	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config config
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err)
	}

	chanPosts := make(chan []storage.Post)
	chanErrs := make(chan error)

	myLinks := getRss("config.json", chanErrs)
	for i := range myLinks.LinkArr {
		go parseNews(myLinks.LinkArr[i], chanErrs, chanPosts, config.Period)
	}

	go func() {
		for posts := range chanPosts {
			for i := range posts {
				db.AddPost(posts[i])
			}
		}
	}()

	go func() {
		for err := range chanErrs {
			log.Println("ошибка:", err)
		}
	}()

	err = http.ListenAndServe(":80", srv.api.Router())
	if err != nil {
		log.Fatal(err)
	}

}

func getRss(fileName string, errors chan<- error) config {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		errors <- err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var links config

	json.Unmarshal(byteValue, &links)

	return links
}

func parseNews(link string, errs chan<- error, posts chan<- []storage.Post, period int) {
	for {
		newsPosts, err := rss.RssToStruct(link)
		if err != nil {
			errs <- err
			continue
		}
		posts <- newsPosts
		time.Sleep(time.Minute * time.Duration(period))
	}
}
