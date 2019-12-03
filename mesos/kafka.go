package mesos

import (
	"encoding/json"
	"strconv"

	mesosproto "../proto"
	cfg "../types"
	"github.com/sirupsen/logrus"
)

// Get out Status of the given kafka ID
func statusKafka(id int) *cfg.State {
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

func startKafka(id int) {
	var cmd cfg.Command

	status := statusKafka(id)
	if status != nil {
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
	cmd.TaskName = "Kafka" + strconv.Itoa(id)
	cmd.Hostname = "kafka" + strconv.Itoa(id) + "." + config.Domain
	cmd.Environment.Variables = []*mesosproto.Environment_Variable{
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
	logrus.Debug("Start Container: ", string(d))

	config.CommandChan <- cmd
	logrus.Info("Scheduled Container: ", cmd.Command)

}

func initStartKafka() {
	// Get ZookeeperStatus
	zookeeperState := statusZookeeper(config.ZookeeperMax)
	if zookeeperState == nil {
		return
	}

	if config.KafkaCount <= config.KafkaMax && zookeeperState.Status.GetState() == 1 {
		startKafka(config.KafkaCount)

		config.KafkaCount++
	}
}
