package web

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed index.html assets/*
var EmbeddedFiles embed.FS

var IndexHTML []byte

func AssetFS() fs.FS {
	assets, err := fs.Sub(EmbeddedFiles, "assets")
	if err != nil {
		log.Panicf("load embedded assets: %v", err)
	}

	return assets
}

func init() {
	html, err := EmbeddedFiles.ReadFile("index.html")
	if err != nil {
		log.Panicf("load embedded index: %v", err)
	}

	IndexHTML = html
}
