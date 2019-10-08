package api

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	cfg "../types"
)

// V0StartContainer will start a mesos container
// example:
// curl -X POST 127.0.0.1:10000/v0/sandbox/start -d 'JSON'
/*
{
  "Command": "./download",
  "uris": [
    {
      "value": "https://<BINFILE>",
      "extract": false,
      "executable": true,
      "cache": false,
    }
  ]
}
*/
func V0StartContainer(w http.ResponseWriter, r *http.Request) {
	var cmd cfg.Command
	err := json.NewDecoder(r.Body).Decode(&cmd)
	cmd.ContainerType = "DOCKER"

	if err != nil {
		logrus.Error("Start Container: ", err)
	}

	d, _ := json.Marshal(&cmd)
	logrus.Debug("Start Container: ", string(d))

	config.CommandChan <- cmd
	w.WriteHeader(http.StatusAccepted)
	logrus.Info("Scheduled Container: ", cmd.Command)
}
