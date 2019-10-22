package mesos

import (
	"encoding/json"
	"strconv"

	"../proto"

	cfg "../types"
	"github.com/sirupsen/logrus"
)

func startZookeeper(id, max int) {
	var cmd cfg.Command

	cmd.ContainerType = "DOCKER"
	cmd.ContainerImage = "zookeeper"
	cmd.NetworkMode = "bridge"
	cmd.Shell = false
	cmd.TaskName = "Zookeeper" + strconv.Itoa(id)
	sI := strconv.Itoa(id)
	cmd.Environment.Variables = []*mesosproto.Environment_Variable{
		{
			Name:  func() *string { x := "ZOO_MY_ID"; return &x }(),
			Value: &sI,
		}, {
			Name:  func() *string { x := "ZOO_SERVERS"; return &x }(),
			Value: getZookeeperServerString(id, max),
		},
	}

	cmd.Hostname = "zookeeper" + strconv.Itoa(id) + "." + config.Domain

	d, _ := json.Marshal(&cmd)
	logrus.Debug("Start Container: ", string(d))

	config.CommandChan <- cmd
	logrus.Info("Scheduled Container: ", cmd.Command)
}

func getZookeeperServerString(id int, max int) *string {
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

func createZookeeperServerString(max int) {
	var server string
	for i := 1; i <= max; i++ {
		sI := strconv.Itoa(i)
		server += "zookeeper" + sI + "." + config.Domain + ":2181 "
	}

	config.ZookeeperServers = server
}
