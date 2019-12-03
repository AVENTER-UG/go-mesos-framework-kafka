package api

import (
	"net/http"

	mesos "../mesos"
	"github.com/sirupsen/logrus"
)

// V0ReflateKafka will restart all missing kafka containers
// example:
// curl -X GET 127.0.0.1:10000/v0/kafka/reflate -d 'JSON'
func V0ReflateKafka(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("HTTP GET V0ReflateKafka")

	mesos.SearchMissingKafka()

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	w.Write([]byte("ok"))
}
