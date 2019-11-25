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
	idName := "Kafka" + strconv.Itoa(id)
	for _, element := range config.State {
		if element.Command.TaskName == idName {
			return &element
		}
	}
	return nil
}

func startKafka() {
	var cmd cfg.Command

	// Get ZookeeperStatus
	zookeeperState := statusZookeeper(1)
	if zookeeperState == nil {
		return
	}

	if config.KafkaCount <= config.KafkaMax && zookeeperState.Status.GetState() == 1 {
		// If there was no kafka started before, then do not check the state
		if config.KafkaCount > 1 {
			state := statusKafka(config.KafkaCount - 1)
			if state == nil {
				return
			}
			if state.Status.GetState() != 1 {
				logrus.Debug("TaskStatus ", state.Status.GetState())
				return
			}
		}

		cmd.ContainerType = "DOCKER"
		cmd.ContainerImage = "wurstmeister/kafka"
		cmd.NetworkMode = "bridge"
		cmd.Shell = false
		cmd.TaskName = "Kafka" + strconv.Itoa(config.KafkaCount)
		cmd.Hostname = "kafka" + strconv.Itoa(config.KafkaCount) + "." + config.Domain
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

		config.KafkaCount++
	}
}
