package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Scheme == "" {
		uStr := fmt.Sprintf("http://%s%s", r.Host, r.RequestURI)
		u, _ := url.Parse(uStr)
		r.URL = u
	}

	//http: Request.RequestURI can't be set in client requests.
	r.RequestURI = ""
	r.Header.Del("Proxy-Connection")
	r.Header.Del("Connection")

	cookie, err := cookiejar.New(nil)
	cookie.SetCookies(r.URL, r.Cookies())
	client := &http.Client{Jar: cookie}
	resp, err := client.Do(r)
	if err != nil {
		log.Println("client Do err :", err.Error())
		return
	}
	defer resp.Body.Close()

	 for k, v := range resp.Header {
	 	for _, vv := range v {
	 		w.Header().Add(k, vv)
	 	}
	 }

	// for _, c := range resp.Cookies() {
	// 	w.Header().Set(c.Name, c.Value)
	// }
	w.WriteHeader(resp.StatusCode)

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		log.Println("ioutil ReadAll err :", err)
		panic(err)
	}
	w.Write(result)
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Start serving on port 7788")
	http.ListenAndServe(":7788", nil)
	os.Exit(0)
}
