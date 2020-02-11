package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	mesos "../mesos"
)

// V0ScaleKafka will scale the kafka service
// example:
// curl -X GET 127.0.0.1:10000/v0/kafka/scale/{count of instances} -d 'JSON'
func V0ScaleKafka(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	auth := CheckAuth(r, w)

	if vars == nil || !auth {
		return
	}

	d := []byte("nok")

	if vars["count"] != "" {
		newCount, _ := strconv.Atoi(vars["count"])
		oldCount := config.KafkaMax
		config.KafkaMax = newCount
		i := (newCount - oldCount) * -1

		// Scale Up
		if newCount > oldCount {
			logrus.Info("Kafka Scale Up ", i)
			for x := oldCount; x < newCount; x++ {
				mesos.StartKafka(x)
			}
		}

		// Scale Down
		if newCount < oldCount {
			logrus.Info("Kafka Scale Down ", i)

			for x := newCount; x < oldCount; x++ {
				task := mesos.StatusKafka(x)
				id := *task.Status.TaskId.Value
				ret := mesos.Kill(id)

				logrus.Info("V0TaskKill: ", ret)
				config.KafkaCount--
			}
		}

		d = []byte(strconv.Itoa(newCount - oldCount))
	}

	logrus.Debug("HTTP GET V0ScaleKafka: ", string(d))

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	w.Write(d)
}
