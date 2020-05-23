package monitor

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type MonitorServer struct {
}

func (m *MonitorServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(strings.Replace(r.URL.Path, "/", "", -1))
	//path = strings.Replace(r.URL.Path, ".", "", -1)
	if r.URL.Path == "/" {
		path = "index.html"
	}

	if strings.HasSuffix(r.URL.Path, ".css") {
		w.Header().Set("Content-Type", "text/css")
	} else if strings.HasSuffix(r.URL.Path, ".js") {
		w.Header().Set("Content-Type", "text/javascript")
	} else if strings.HasSuffix(r.URL.Path, ".png") {
		w.Header().Set("Content-Type", "image/png")
	} else {
		w.Header().Set("Content-Type", "text/html")
	}
	if strings.HasPrefix(path, "data-dir") {
		return
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	w.Write(bytes)
}
