package api

import (
	"net/http"

	mesos "../mesos"
	"github.com/sirupsen/logrus"
)

// V0ReflateMissingZookeeper will scale the zookeeper service
// example:
// curl -X GET 127.0.0.1:10000/v0/zookeeper/reflate -d 'JSON'
func V0ReflateZookeeper(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("HTTP GET V0RestartMissingZookeeper")

	mesos.SearchMissingZookeeper()

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	w.Write([]byte("ok"))
}