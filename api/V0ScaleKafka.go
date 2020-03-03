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
		logrus.Debug("V0ScaleKafka: oldCount: ", oldCount)
		config.KafkaMax = newCount
		i := (newCount - oldCount)
		// change the number to be positiv
		if i < 0 {
			i = i * -1
		}

		// Scale Up
		if newCount > oldCount {
			logrus.Info("Kafka Scale Up ", i)
		}

		// Scale Down
		if newCount < oldCount {
			logrus.Info("Kafka Scale Down ", i)

			for x := newCount; x < oldCount; x++ {
				task := mesos.StatusKafka(x)
				if task.Status.TaskId != nil {
					id := *task.Status.TaskId.Value
					ret := mesos.Kill(id)

					logrus.Info("V0TaskKill: ", ret)
					config.KafkaCount--
				} else {
					logrus.Debug("V0ScaleKafka: Missing TaskID")
				}
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
