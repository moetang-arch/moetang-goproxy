package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/goproxy/goproxy"
	"github.com/goproxy/goproxy/cacher"
)

func main() {

	// setup http proxy for remote connections
	// BEFORE initialize goproxy
	if err := os.Setenv("HTTP_PROXY", "http://127.0.0.1:1080"); err != nil {
		panic(err)
	}
	if err := os.Setenv("HTTPS_PROXY", "http://127.0.0.1:1080"); err != nil {
		panic(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cachePath := fmt.Sprint(wd, string(os.PathSeparator), "data")
	if err := os.MkdirAll(cachePath, os.ModeDir); err != nil {
		panic(err)
	}
	diskCache := &cacher.Disk{
		Root: cachePath,
	}

	log.Println("[DEBUG] disk cache path:", diskCache.Root)

	p := goproxy.New()
	p.Cacher = diskCache
	p.GoBinMaxWorkers = 16

	// setup goproxy configurations
	p.GoBinEnv = append(p.GoBinEnv, "GOPROXY=direct")
	// setup proxy sumdb
	p.ProxiedSUMDBs = append(
		p.ProxiedSUMDBs,
		"sum.golang.org",
		"sum.golang.google.cn",
	)

	if err := http.ListenAndServe(":3008", &LoggingHandler{
		p: p,
	}); err != nil {
		panic(err)
	}
}

type LoggingHandler struct {
	p *goproxy.Goproxy
}

func (this *LoggingHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Println("[DEBUG] request:", req.URL)
	this.p.ServeHTTP(resp, req)
}
