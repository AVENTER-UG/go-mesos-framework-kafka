package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	"./api"
	"./mesos"
	mesosproto "./proto"
	cfg "./types"

	util "git.aventer.biz/AVENTER/util"
	"github.com/Showmax/go-fqdn"
	"github.com/sirupsen/logrus"
)

func main() {
	util.SetLogging(config.LogLevel, config.EnableSyslog, config.AppName)
	logrus.Println(config.AppName + " build" + config.MinVersion)

	hostname := fqdn.Get()
	listen := ":10000"

	logrus.Info(hostname)

	failoverTimeout := 5000.0
	checkpoint := true
	webuiurl := fmt.Sprintf("http://%s%s", hostname, listen)

	config.FrameworkInfoFile = fmt.Sprintf("%s/%s", os.TempDir(), "framework.json")
	config.CommandChan = make(chan cfg.Command, 100)
	config.Hostname = hostname
	config.Listen = listen

	util.SetLogging(config.LogLevel, config.EnableSyslog, config.AppName)

	config.State = map[string]cfg.State{}

	config.FrameworkInfo.User = &config.FrameworkUser
	config.FrameworkInfo.Name = &config.FrameworkName
	config.FrameworkInfo.Hostname = &hostname
	config.FrameworkInfo.WebuiUrl = &webuiurl
	config.FrameworkInfo.FailoverTimeout = &failoverTimeout
	config.FrameworkInfo.Checkpoint = &checkpoint
	config.FrameworkInfo.Principal = &config.Principal
	config.FrameworkInfo.Capabilities = []*mesosproto.FrameworkInfo_Capability{
		{Type: mesosproto.FrameworkInfo_Capability_RESERVATION_REFINEMENT.Enum()},
	}

	mesos.SetConfig(&config)
	api.SetConfig(&config)

	http.Handle("/", api.Commands())

	go func() {
		http.ListenAndServe(listen, nil)
	}()
	logrus.Fatal(mesos.Subscribe())
}
