package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	mesos "go-mesos-framework-kafka/mesos"
	mesosproto "go-mesos-framework-kafka/proto"
)

// V0ScaleZookeeper will scale the zookeeper service
// example:
// curl -X GET 127.0.0.1:10000/v0/zookeeper/scale/{count of instances} -d 'JSON'
func V0ScaleZookeeper(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := CheckAuth(r, w)

	if vars == nil || !auth {
		return
	}

	d := []byte("nok")

	if vars["count"] != "" {
		newCount, _ := strconv.Atoi(vars["count"])
		oldCount := config.ZookeeperMax
		logrus.Debug("V0ScaleZookeeper: oldCount: ", oldCount)
		config.ZookeeperMax = newCount
		i := (newCount - oldCount)
		// change the number to be positiv
		if i < 0 {
			i = i * -1
		}

		// Scale Up
		if newCount > oldCount {
			logrus.Info("Zookeeper Scale Up ", i)
			revive := &mesosproto.Call{
				Type: mesosproto.Call_REVIVE.Enum(),
			}
			mesos.Call(revive)
		}

		// Scale Down
		if newCount < oldCount {
			logrus.Info("Zookeeper Scale Down ", i)

			for x := newCount; x < oldCount; x++ {
				task := mesos.StatusZookeeper(x)
				if task.Status.TaskId != nil {
					id := *task.Status.TaskId.Value
					ret := mesos.Kill(id)

					logrus.Info("V0TaskKill: ", ret)
					config.ZookeeperCount--
				} else {
					logrus.Debug("V0ScaleZookeeper: Missing TaskID")
				}
			}
		}

		d = []byte(strconv.Itoa(newCount - oldCount))
	}

	logrus.Debug("HTTP GET V0ScaleZookeeper: ", string(d))

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	w.Write(d)
}
