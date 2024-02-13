package main

import (
	"os"
	"strconv"
	"strings"

	util "github.com/AVENTER-UG/util/util"

	cfg "go-mesos-framework-kafka/types"
)

var config cfg.Config

func init() {
	config.ZookeeperMax = 0
	config.KafkaMax = 0
	config.ZookeeperCount = 0
	config.KafkaCount = 0

	config.FrameworkUser = util.Getenv("FRAMEWORK_USER", "root")
	config.FrameworkName = util.Getenv("FRAMEWORK_NAME", "kafka")
	config.FrameworkRole = util.Getenv("FRAMEWORK_ROLE", "kafka")
	config.FrameworkPort = util.Getenv("FRAMEWORK_PORT", "10000")
	config.FrameworkInfoFilePath = util.Getenv("FRAMEWORK_STATEFILE_PATH", "/tmp")
	config.Principal = os.Getenv("MESOS_PRINCIPAL")
	config.Username = os.Getenv("MESOS_USERNAME")
	config.Password = os.Getenv("MESOS_PASSWORD")
	config.MesosMasterServer = os.Getenv("MESOS_MASTER")
	config.LogLevel = util.Getenv("LOGLEVEL", "info")
	config.Domain = os.Getenv("DOMAIN")
	config.ZookeeperMax, _ = strconv.Atoi(os.Getenv("ZOOKEEPER_COUNT"))
	config.KafkaMax, _ = strconv.Atoi(os.Getenv("KAFKA_COUNT"))
	config.ResCPU, _ = strconv.ParseFloat(util.Getenv("RES_CPU", "0.1"), 64)
	config.ResMEM, _ = strconv.ParseFloat(util.Getenv("RES_MEM", "1200"), 64)
	config.Credentials.Username = os.Getenv("AUTH_USERNAME")
	config.Credentials.Password = os.Getenv("AUTH_PASSWORD")
	config.AppName = "Mesos Kafka Framework"
	config.ZookeeperCustomString = os.Getenv("ZOOKEEPER_CUSTOM_DOMAIN")
	config.KafkaCustomString = os.Getenv("KAFKA_CUSTOM_DOMAIN")
	config.ImageZookeeper = util.Getenv("IMAGE_ZOOKEEPER", "zookeeper:3.5.7")
	config.ImageKafka = util.Getenv("IMAGE_KAFKA", "confluentinc/cp-kafka:5.4.1")
	config.VolumeDriver = util.Getenv("VOLUME_DRIVER", "local")
	config.VolumeKafka = util.SplitString(util.Getenv("VOLUME_KAFKA", "/data/kafka,"))
	config.VolumeZookeeper = util.SplitString(util.Getenv("VOLUME_ZOOKEEPER", "/data/zookeeper,"))

	// The comunication to the mesos server should be via ssl or not
	if strings.Compare(os.Getenv("MESOS_SSL"), "true") == 0 {
		config.MesosSSL = true
	} else {
		config.MesosSSL = false
	}

}
