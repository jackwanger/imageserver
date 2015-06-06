package main

import (
	"crypto/sha256"
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/pierrre/githubhook"
	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	imageserver_cache_memory "github.com/pierrre/imageserver/cache/memory"
	imageserver_cache_redis "github.com/pierrre/imageserver/cache/redis"
	imageserver_http "github.com/pierrre/imageserver/http"
	imageserver_http_nfntresize "github.com/pierrre/imageserver/http/nfntresize"
	imageserver_native "github.com/pierrre/imageserver/native"
	_ "github.com/pierrre/imageserver/native/bmp"
	_ "github.com/pierrre/imageserver/native/gif"
	_ "github.com/pierrre/imageserver/native/jpeg"
	imageserver_native_nfntresize "github.com/pierrre/imageserver/native/nfntresize"
	_ "github.com/pierrre/imageserver/native/png"
	_ "github.com/pierrre/imageserver/native/tiff"
	imageserver_testdata "github.com/pierrre/imageserver/testdata"
)

func main() {
	var httpAddr string
	flag.StringVar(&httpAddr, "http", ":8080", "Http")
	var gitHubWebhookSecret string
	flag.StringVar(&gitHubWebhookSecret, "github-webhook-secret", "", "GitHub webhook secret")
	flag.Parse()

	log.Println("Start")
	log.Printf("Go version: %s", runtime.Version())
	log.Printf("Go max procs: %d", runtime.GOMAXPROCS(0))

	startHTTPServerAddr(httpAddr, gitHubWebhookSecret)
}

func startHTTPServerAddr(addr string, gitHubWebhookSecret string) {
	http.Handle("/", newImageHTTPHandler())
	if gitHubWebhookSecret != "" {
		http.Handle("/github_webhook", newGitHubWebhookHTTPHandler(gitHubWebhookSecret))
	}
	log.Printf("Start HTTP server on %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Panic(err)
	}
}

func newGitHubWebhookHTTPHandler(secret string) http.Handler {
	return &githubhook.Handler{
		Secret: secret,
		Delivery: func(event string, deliveryID string, payload interface{}) {
			log.Printf("Received GitHub webhook: %s", event)
			if event == "push" {
				delay := time.Duration(5 * time.Second)
				log.Printf("Killing process in %s", delay)
				time.AfterFunc(delay, func() {
					log.Println("Killing process now")
					os.Exit(0)
				})
			}
		},
		Error: func(err error, req *http.Request) {
			log.Printf("GitHub webhook error: %s", err)
		},
	}
}

func newImageHTTPHandler() http.Handler {
	handler := http.Handler(&imageserver_http.Handler{
		Parser:   newParser(),
		Server:   newServer(),
		ETagFunc: imageserver_http.NewParamsHashETagFunc(sha256.New),
		ErrorFunc: func(err error, request *http.Request) {
			log.Println("Error:", err)
		},
	})
	handler = &imageserver_http.ExpiresHandler{
		Handler: handler,
		Expires: time.Duration(7 * 24 * time.Hour),
	}
	return handler
}

func newParser() imageserver_http.Parser {
	return &imageserver_http.ListParser{
		&imageserver_http.SourceParser{},
		&imageserver_http_nfntresize.Parser{},
		&imageserver_http.FormatParser{},
		&imageserver_http.QualityParser{},
	}
}

func newServer() imageserver.Server {
	server := newServerTestData()
	server = newServerNative(server)
	server = newServerLimit(server)
	server = newServerCache(server)
	return server
}

func newServerTestData() imageserver.Server {
	return imageserver_testdata.Server
}

func newServerNative(server imageserver.Server) imageserver.Server {
	return &imageserver_native.Server{
		Server:    server,
		Processor: &imageserver_native_nfntresize.Processor{},
	}
}

func newServerLimit(server imageserver.Server) imageserver.Server {
	return imageserver.NewLimitServer(server, runtime.GOMAXPROCS(0)*2)
}

func newServerCache(server imageserver.Server) imageserver.Server {
	keyGenerator := imageserver_cache.NewParamsHashKeyGenerator(sha256.New)
	server = &imageserver_cache.Server{
		Server:       server,
		Cache:        newCacheRedis(),
		KeyGenerator: keyGenerator,
	}
	server = &imageserver_cache.Server{
		Server:       server,
		Cache:        newCacheMemory(),
		KeyGenerator: keyGenerator,
	}
	return server
}

func newCacheRedis() imageserver_cache.Cache {
	cache := imageserver_cache.Cache(&imageserver_cache_redis.Cache{
		Pool: &redigo.Pool{
			Dial: func() (redigo.Conn, error) {
				return redigo.Dial("tcp", "localhost:6379")
			},
			MaxIdle: 50,
		},
		Expire: time.Duration(7 * 24 * time.Hour),
	})
	cache = &imageserver_cache.IgnoreError{
		Cache: cache,
	}
	cache = &imageserver_cache.Async{
		Cache: cache,
	}
	return cache
}

func newCacheMemory() imageserver_cache.Cache {
	return imageserver_cache_memory.New(100 * 1024 * 1024)
}
