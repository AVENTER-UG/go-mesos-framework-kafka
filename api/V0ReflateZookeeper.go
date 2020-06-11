package api

import (
	"net/http"

	mesos "../mesos"
	mesosproto "../proto"
	"github.com/sirupsen/logrus"
)

// V0ReflateMissingZookeeper will scale the zookeeper service
// example:
// curl -X GET 127.0.0.1:10000/v0/zookeeper/reflate -d 'JSON'
func V0ReflateZookeeper(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("HTTP GET V0RestartMissingZookeeper")
	auth := CheckAuth(r, w)

	if !auth {
		return
	}

	revive := &mesosproto.Call{
		Type: mesosproto.Call_REVIVE.Enum(),
	}
	mesos.Call(revive)

	mesos.SearchMissingZookeeper()

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	w.Write([]byte("ok"))
}
