package scraper

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestTransport(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadFile("_examples/challenge.html")
		if err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Server", "cloudflare-nginx")
		w.WriteHeader(503)
		w.Write(b)
	}))
	defer ts.Close()

	scraper, err := NewTransport(http.DefaultTransport)
	if err != nil {
		t.Fatal(err)
	}

	c := http.Client{
		Transport: scraper,
	}

	res, err := c.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestAll(t *testing.T) {

	c, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	res, err := c.Get("https://hidemyna.me/en/proxy-list")
	if err != nil {
		t.Fatal(err)
	}

	var body bytes.Buffer
	if _, err := io.Copy(&body, res.Body); err != nil {
		t.Fatal(err)
	}

	r, _ := regexp.Compile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)
	fmt.Println(body.String())
	if len(r.FindAllString(body.String(), -1)) != 65{
		t.Fatal("should be 65 ips")
	}
}