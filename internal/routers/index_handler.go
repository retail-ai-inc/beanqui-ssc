package routers

import (
	"net/http"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type Index struct {
}

func (t *Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.RequestURI
	if strings.HasSuffix(url, ".vue") {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	var dir string = "./"
	_, f, _, ok := runtime.Caller(0)
	if ok {
		dir = filepath.Dir(f)
	}

	hdl := http.FileServer(http.Dir(path.Join(dir, "../../ui/")))
	hdl.ServeHTTP(w, r)
}
