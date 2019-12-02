package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// V0ScaleZookeeper will scale the zookeeper service
// example:
// curl -X GET 127.0.0.1:10000/v0/zookeeper/scale/{count of instances} -d 'JSON'
func V0ScaleZookeeper(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars == nil {
		return
	}

	d := []byte("{}")

	logrus.Debug("HTTP GET V0ScaleZookeeper: ", string(d))

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	w.Write(d)
}
