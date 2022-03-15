package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/google/goblet"
)

var (
	port      = flag.Int("port", 8080, "port to listen to")
	cacheRoot = flag.String("cache_root", "", "Root directory of cached repositories")
)

func main() {
	flag.Parse()

	var urlCanonicalizer func(u *url.URL) (*url.URL, error) = func(u *url.URL) (*url.URL, error) {
		ret := url.URL{}
		ret.Scheme = "https"
		ret.Host = u.Host
		ret.Path = u.Path

		if strings.HasSuffix(ret.Path, "/info/refs") {
			ret.Path = strings.TrimSuffix(ret.Path, "/info/refs")
		} else if strings.HasSuffix(ret.Path, "/git-upload-pack") {
			ret.Path = strings.TrimSuffix(ret.Path, "/git-upload-pack")
		} else if strings.HasSuffix(ret.Path, "/git-receive-pack") {
			ret.Path = strings.TrimSuffix(ret.Path, "/git-receive-pack")
		}
		ret.Path = strings.TrimSuffix(ret.Path, ".git")

		return &ret, nil
	}

	var requestAuthorizer func(*http.Request) error = func(r *http.Request) error {
		return nil
	}

	// Using modified version of `managed_repository` which no longer use `TokenSource`.
	// var tokenSource = oauth2.StaticTokenSource(&oauth2.Token{})

	var errorReporter func(*http.Request, error)

	var requestLogger func(r *http.Request, status int, requestSize, responseSize int64, latency time.Duration) = func(r *http.Request, status int, requestSize, responseSize int64, latency time.Duration) {
		dump, err := httputil.DumpRequest(r, false)
		if err != nil {
			return
		}
		log.Printf("%q %d reqsize: %d, respsize %d, latency: %v", dump, status, requestSize, responseSize, latency)
	}

	config := &goblet.ServerConfig{
		LocalDiskCacheRoot: *cacheRoot,
		URLCanonializer:    urlCanonicalizer,
		RequestAuthorizer:  requestAuthorizer,
		// TokenSource:        tokenSource,
		ErrorReporter: errorReporter,
		RequestLogger: requestLogger,
	}
	http.Handle("/", goblet.HTTPHandler(config))
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
