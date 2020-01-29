package api

import (
	"net/http"

	cfg "../types"
	"github.com/sirupsen/logrus"
)

// V0StartCommand will start a simple command
// example:
// curl -X POST 127.0.0.1:10000/v0/cmd/start\?cmd\=python%20-m%20SimpleHTTPServer%209033\&shell=true/false
func V0StartCommand(w http.ResponseWriter, r *http.Request) {
	auth := CheckAuth(r, w)

	if !auth {
		return
	}

	r.ParseForm()
	var cmd cfg.Command
	cmd.Command = r.Form["cmd"][0]
	cmd.TaskName = r.Form["name"][0]
	cmd.ContainerType = "NONE"

	if r.Form["shell"][0] == "true" {
		cmd.Shell = true
	} else {
		cmd.Shell = false
	}

	//write command to channel
	config.CommandChan <- cmd
	w.WriteHeader(http.StatusAccepted)
	logrus.Info("Scheduled Command: ", cmd.Command)
}
