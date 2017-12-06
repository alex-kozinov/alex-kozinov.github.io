package main

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

var urls map[string]string

func makeShorter() string {
	a := len(urls) + 1
	s := ""
	for a > 0 {
		s += string('a' + (a % 26))
		a /= 26
	}
	return s
}

func shorter(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var body []byte
		defer r.Body.Close()
		body, _ = ioutil.ReadAll(r.Body)
		var req map[string]string
		err := json.Unmarshal(body, &req)
		if err == nil {
			url, ok := req["url"]
			if ok {
				key := makeShorter()
				urls[key] = url
				respMap := make(map[string]string)
				respMap["key"] = key
				resp, _ := json.Marshal(respMap)
				w.Write(resp)
				return
			}
		}
		http.Error(w, "", 400);
	} else {
		url, ok := urls[r.RequestURI[1:]]
		if !ok {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, url, 301)
	}

}

func main() {
	urls = make(map[string]string)
	http.HandleFunc("/", shorter)
	http.ListenAndServe(":8082", nil)
}
