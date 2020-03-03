package mesos

import (
	"encoding/json"
	"strconv"

	mesosproto "../proto"
	cfg "../types"
	"github.com/sirupsen/logrus"
)

// SearchMissingKafka Check if all kafkas are running. If one is missing, restart it.
func SearchMissingKafka() {
	if config.State != nil {
		for i := 0; i < config.KafkaMax; i++ {
			state := *StatusKafka(i).Status.State
			if state != mesosproto.TaskState_TASK_RUNNING {
				logrus.Debug("Missing Kafka: ", i)
				CreateZookeeperServerString()
				StartKafka(i)
			}
		}
	}
}

// Get out Status of the given kafka ID
func StatusKafka(id int) *cfg.State {
	if config.State != nil {
		for _, element := range config.State {
			if element.Status != nil {
				if element.Command.InternalID == id {
					return &element
				}
			}
		}
	}
	return nil
}

// start kafka with the given id
func StartKafka(id int) {
	var cmd cfg.Command

	// be sure, that there is no kafka with this id already running
	status := StatusKafka(id)
	if status != nil {
		if status.Status.State == mesosproto.TaskState_TASK_STAGING.Enum() {
			logrus.Info("startKafka: kafka is staging ", id)
			return
		}
		if status.Status.State == mesosproto.TaskState_TASK_STARTING.Enum() {
			logrus.Info("startKafka: kafka is starting ", id)
			return
		}
		if status.Status.State == mesosproto.TaskState_TASK_RUNNING.Enum() {
			logrus.Info("startKafka: kafka already running ", id)
			return
		}
	}

	cmd.ContainerType = "DOCKER"
	cmd.ContainerImage = "wurstmeister/kafka"
	cmd.NetworkMode = "bridge"
	cmd.Shell = false
	cmd.InternalID = id
	cmd.IsKafka = true
	cmd.TaskName = "kafka" + strconv.Itoa(id)
	cmd.Hostname = "kafka" + strconv.Itoa(id) + config.KafkaCustomString + "." + config.Domain
	cmd.Environment.Variables = []*mesosproto.Environment_Variable{
		{
			Name:  func() *string { x := "SERVICE_NAME"; return &x }(),
			Value: &cmd.TaskName,
		},
		{
			Name:  func() *string { x := "KAFKA_ADVERTISED_HOST_NAME"; return &x }(),
			Value: &cmd.Hostname,
		},
		{
			Name:  func() *string { x := "KAFKA_ZOOKEEPER_CONNECT"; return &x }(),
			Value: &config.ZookeeperServers,
		},
		{
			Name:  func() *string { x := "KAFKA_LISTENER_SECURITY_PROTOCOL_MAP"; return &x }(),
			Value: func() *string { x := "INSIDE:PLAINTEXT,OUTSIDE:PLAINTEXT"; return &x }(),
		},
		{
			Name:  func() *string { x := "KAFKA_ADVERTISED_LISTENERS"; return &x }(),
			Value: func() *string { x := "INSIDE://:9092,OUTSIDE://_" + cmd.Hostname + ":9094"; return &x }(),
		},
		{
			Name:  func() *string { x := "KAFKA_LISTENERS"; return &x }(),
			Value: func() *string { x := "INSIDE://:9092,OUTSIDE://:9094"; return &x }(),
		},
		{
			Name:  func() *string { x := "KAFKA_INTER_BROKER_LISTENER_NAME"; return &x }(),
			Value: func() *string { x := "INSIDE"; return &x }(),
		},
	}

	d, _ := json.Marshal(&cmd)
	logrus.Debug("Scheduled Kafka: ", string(d))

	config.CommandChan <- cmd
	logrus.Info("Scheduled Kafka")

}

// the first run should be in ta strict order.
func initStartKafka() {
	// Start kafka only if the zookeeper is running
	zookeeperState := StatusZookeeper(config.ZookeeperMax - 1)
	if zookeeperState == nil {
		return
	}

	if config.KafkaCount <= (config.KafkaMax-1) && zookeeperState.Status.GetState() == 1 {
		StartKafka(config.KafkaCount)
		config.KafkaCount++
	}
}
