package mesos

import (
	"encoding/json"
	"strconv"

	mesosproto "../proto"

	cfg "../types"
	"github.com/sirupsen/logrus"
)

// Get out Status of the given zookeeper ID
func statusZookeeper(id int) *cfg.State {
	idName := "Zookeeper" + strconv.Itoa(id)
	for _, element := range config.State {
		if element.Command.TaskName == idName {
			return &element
		}
	}
	return nil
}

// Start a zookeeper container, but only if the foreunner  zookeeper is in the running state
func startZookeeper() {
	var cmd cfg.Command

	if config.ZookeeperCount <= config.ZookeeperMax {
		// If there was no zookeper started before, then do not check the state
		if config.ZookeeperCount > 1 {
			state := statusZookeeper(config.ZookeeperCount - 1)
			if state == nil {
				return
			}
			if state.Status.GetState() != 1 {
				logrus.Debug("TaskStatus ", state.Status.GetState())
				return
			}
		}
		cmd.ContainerType = "DOCKER"
		cmd.ContainerImage = "zookeeper"
		cmd.NetworkMode = "bridge"
		cmd.Shell = false
		cmd.TaskName = "Zookeeper" + strconv.Itoa(config.ZookeeperCount)
		sI := strconv.Itoa(config.ZookeeperCount)
		cmd.Environment.Variables = []*mesosproto.Environment_Variable{
			{
				Name:  func() *string { x := "ZOO_MY_ID"; return &x }(),
				Value: &sI,
			}, {
				Name:  func() *string { x := "ZOO_SERVERS"; return &x }(),
				Value: getZookeeperServerString(config.ZookeeperCount),
			},
		}

		cmd.Hostname = "zookeeper" + strconv.Itoa(config.ZookeeperCount) + "." + config.Domain

		d, _ := json.Marshal(&cmd)
		logrus.Debug("Start Container: ", string(d))

		config.CommandChan <- cmd
		logrus.Info("Scheduled Container")

		config.ZookeeperCount++
	}
}

func getZookeeperServerString(id int) *string {
	max := config.ZookeeperMax
	var server string
	for i := 1; i <= max; i++ {
		sI := strconv.Itoa(i)
		if i == id {
			server += "server." + sI + "=0.0.0.0:2888:3888;2181 "
		} else {
			server += "server." + sI + "=zookeeper" + sI + "." + config.Domain + ":2888:3888;2181 "
		}
	}

	return &server
}

func createZookeeperServerString() {
	max := config.ZookeeperMax
	var server string
	for i := 1; i <= max; i++ {
		sI := strconv.Itoa(i)
		server += "zookeeper" + sI + "." + config.Domain + ":2181 "
	}

	config.ZookeeperServers = server
}
