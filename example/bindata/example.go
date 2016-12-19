package main

import (
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/insionng/macross"
	"github.com/macross-contrib/statio"
	"net/http"
	"strings"
)

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {

	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func BinaryFileSystem(root string) *binaryFileSystem {
	//fs := &assetfs.AssetFS{Asset, AssetDir, root}
	fs := &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: root}
	return &binaryFileSystem{
		fs,
	}
}

// Usage
// $ go-bindata data/
// $ go build && ./bindata
//
func main() {
	m := macross.New()

	m.Use(statio.Serve("/static", BinaryFileSystem("data")))
	m.Get("/ping", func(self *macross.Context) error {
		return self.String("pong")
	})
	// Listen and Server in 0.0.0.0:8080
	m.Listen(":8080")
}
