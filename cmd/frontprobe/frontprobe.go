// frontprobe.go - quickly test domain fronting availability.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of frontier, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/unkaktus/frontier"
)

func main() {
	log.SetFlags(0)
	var dump = flag.Bool("dump", false, "dump response body")
	var hostname = flag.String("a", "", "TCP address")
	var front = flag.String("n", "", "SNI")
	var hostHeader = flag.String("h", "", "Host header")
	var path = flag.String("p", "/", "path")
	var method = flag.String("m", http.MethodGet, "method")
	flag.Parse()

	fr := frontier.New(http.DefaultTransport, *front, *hostname)
	u := &url.URL{
		Scheme: "https",
		Host:   *hostHeader,
		Path:   *path,
	}
	req, err := http.NewRequest(*method, u.String(), http.NoBody)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := fr.RoundTrip(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Printf("status code: %v", resp.StatusCode)
	location := resp.Header.Get("Location")
	if location != "" {
		log.Printf("location: %s", location)
	}
	if *dump {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("body: %s", body)
	}
}
