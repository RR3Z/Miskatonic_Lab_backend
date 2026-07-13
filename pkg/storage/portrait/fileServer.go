package portrait

import (
	"net/http"
	"path/filepath"
	"strings"
)

const PublicPathPrefix = "/uploads/portraits/"

type FileServer struct {
	directory string
}

func NewFileServer(store *LocalStore) *FileServer {
	return &FileServer{directory: store.directory}
}

func (s *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	name := strings.TrimPrefix(r.URL.Path, PublicPathPrefix)
	if _, ok := managedFileName(StorageKeyPrefix + name); !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	http.ServeFile(w, r, filepath.Join(s.directory, name))
}
