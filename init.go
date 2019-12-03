package main

import (
	"os"
	"strconv"

	cfg "./types"
)

var config cfg.Config

func init() {
	config.ZookeeperMax = 0
	config.KafkaMax = 0
	config.ZookeeperCount = 1
	config.KafkaCount = 1

	config.FrameworkUser = os.Getenv("FRAMEWORK_USER")
	config.FrameworkName = os.Getenv("FRAMEWORK_NAME")
	config.FrameworkPort = os.Getenv("FRAMEWORK_PORT")
	config.Principal = os.Getenv("MESOS_PRINCIPAL")
	config.Username = os.Getenv("MESOS_USERNAME")
	config.Password = os.Getenv("MESOS_PASSWORD")
	config.MesosMasterServer = os.Getenv("MESOS_MASTER")
	config.LogLevel = os.Getenv("LOGLEVEL")
	config.Domain = os.Getenv("DOMAIN")
	config.ZookeeperMax, _ = strconv.Atoi(os.Getenv("ZOOKEEPER_COUNT"))
	config.KafkaMax, _ = strconv.Atoi(os.Getenv("KAFKA_COUNT"))
	config.ResCPU, _ = strconv.ParseFloat(os.Getenv("RES_CPU"), 64)
	config.ResMEM, _ = strconv.ParseFloat(os.Getenv("RES_MEM"), 64)
	config.AppName = "Mesos Kafka Framework"
}
