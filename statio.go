package statio

import (
	"github.com/insionng/macross"
	"net/http"
	"os"
	"path"
	"strings"
)

type (
	onlyfilesFS struct {
		fs http.FileSystem
	}
	neuteredReaddirFile struct {
		http.File
	}
)

// Dir returns a http.Filesystem that can be used by http.FileServer(). It is used interally
// in router.Static().
// if listDirectory == true, then it works the same as http.Dir() otherwise it returns
// a filesystem that prevents http.FileServer() to list the directory files.
func Dir(root string, listDirectory bool) http.FileSystem {
	fs := http.Dir(root)
	if listDirectory {
		return fs
	}
	return &onlyfilesFS{fs}
}

// Conforms to http.Filesystem
func (fs onlyfilesFS) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}

// Overrides the http.File default implementation
func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	// this disables directory listing
	return nil, nil
}

type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}

type localFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
}

func LocalFile(root string, indexes bool) *localFileSystem {
	return &localFileSystem{
		FileSystem: Dir(root, indexes),
		root:       root,
		indexes:    indexes,
	}
}

func (l *localFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		name := path.Join(l.root, p)
		stats, err := os.Stat(name)
		if err != nil {
			return false
		}
		if !l.indexes && stats.IsDir() {
			return false
		}
		return true
	}
	return false
}

func ServeRoot(urlPrefix, root string) macross.Handler {
	return Serve(urlPrefix, LocalFile(root, false))
}

// Static returns a middleware handler that serves static files in the given directory.
func Serve(urlPrefix string, fs ServeFileSystem) macross.Handler {
	fileserver := http.FileServer(fs)
	if len(urlPrefix) != 0 {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *macross.Context) error {
		rpath := string(c.Request.URI().Path())
		if fs.Exists(urlPrefix, rpath) {
			content, err := fs.Open(rpath)
			if err != nil {
				return err
			}
			fi, _ := content.Stat()
			err = c.ServeContent(content, fi.Name(), fi.ModTime())
			c.Abort()
			return err
		}
		return c.Next()
	}
}
