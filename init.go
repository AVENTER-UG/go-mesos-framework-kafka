package main

import (
	"os"

	cfg "./types"
)

var config cfg.Config

func init() {
	config.FrameworkUser = os.Getenv("FRAMEWORK_USER")
	config.FrameworkName = os.Getenv("FRAMEWORK_NAME")
	config.Principal = os.Getenv("MESOS_PRINCIPAL")
	config.Username = os.Getenv("MESOS_USERNAME")
	config.Password = os.Getenv("MESOS_PASSWORD")
	config.MesosMasterServer = os.Getenv("MESOS_MASTER")
	config.LogLevel = os.Getenv("LOGLEVEL")
	config.Domain = os.Getenv("DOMAIN")

	config.AppName = "Mesos TestFramework"
}
