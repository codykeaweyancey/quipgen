package main

import (
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

const (
	assetFilePath = "frontend"
	assetURIPath = "asset/"
	apiURIPath = "api/"
	indexHTML = "index.html"
	notfoundHTML = "notfound.html"
)

var contentTypes = map[string]string{
	"html": "text/html",
	"css": "text/css",
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	workDir, err := os.Getwd()
	panicOnErr(err)
	handler := &HttpHandler{AssetsDir: filepath.Join(workDir, assetFilePath)}
	handler.LoadAssets()
	go func () {
		for {
			<- time.After(1 * time.Second)
			handler.LoadAssets()
		}
	}()
	println("Starting web server...")
	http.ListenAndServe(":8080", handler)
}

type HttpHandler struct {
	AssetsDir string
	Assets atomic.Value
}

func (w *HttpHandler) LoadAssets() error {
	assets := map[string][]byte{}
	err := filepath.Walk(w.AssetsDir, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			buf, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			assets[strings.TrimPrefix(strings.TrimPrefix(path, w.AssetsDir), string(filepath.Separator))] = buf
		}
		return nil
	})
	if err != nil {
		return err
	}
	w.Assets.Store(assets)
	return nil
}

func (w *HttpHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	assets := w.Assets.Load().(map[string][]byte)

	normalizedPath := strings.TrimPrefix(req.URL.Path, "/")
	println(normalizedPath)
	if strings.HasPrefix(normalizedPath, apiURIPath) {

	} else if strings.HasPrefix(normalizedPath, assetURIPath) {
		assetName := strings.TrimPrefix(normalizedPath, assetURIPath)
		asset, hasAsset := assets[assetName]
		if hasAsset {
			res.Header().Set("Content-Type", contentTypes[filepath.Ext(assetName)])
			res.Write(asset)
		} else {
			res.WriteHeader(404)
			res.Write(assets[notfoundHTML])
		}
	} else {
		println("Index Served")
		res.Write(assets[indexHTML])
	}
}