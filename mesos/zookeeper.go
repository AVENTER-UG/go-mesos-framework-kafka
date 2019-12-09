package mesos

import (
	"encoding/json"
	"strconv"

	mesosproto "../proto"

	cfg "../types"
	"github.com/sirupsen/logrus"
)

// SearchMissingZookeeper Check if all zookeepers are running. If one is missing, restart it.
func SearchMissingZookeeper() {
	if config.State != nil {
		for i := 1; i <= config.ZookeeperMax; i++ {
			if statusZookeeper(i) == nil {
				logrus.Debug("Missing Zookeeper: ", i)
				statusZookeeper(i)
			}
		}
	}
}

// Get out Status of the given zookeeper ID
func statusZookeeper(id int) *cfg.State {
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

// Start a zookeeper container with the given ID
func startZookeeper(id int) {
	var cmd cfg.Command

	// before we will start a new zookeeper, we should be sure its not already running
	status := statusZookeeper(id)
	if status != nil {
		if status.Status.State == mesosproto.TaskState_TASK_STAGING.Enum() {
			logrus.Info("startZookeeper: zookeeper is staging ", id)
			return
		}
		if status.Status.State == mesosproto.TaskState_TASK_STARTING.Enum() {
			logrus.Info("startZookeeper: zookeeper is starting ", id)
			return
		}
		if status.Status.State == mesosproto.TaskState_TASK_RUNNING.Enum() {
			logrus.Info("startZookeeper: zookeeper already running ", id)
			return
		}
	}

	cmd.ContainerType = "DOCKER"
	cmd.ContainerImage = "zookeeper"
	cmd.NetworkMode = "bridge"
	cmd.Shell = false
	cmd.TaskName = "Zookeeper" + strconv.Itoa(id)
	cmd.InternalID = id
	cmd.IsZookeeper = true
	sI := strconv.Itoa(id)
	cmd.Environment.Variables = []*mesosproto.Environment_Variable{
		{
			Name:  func() *string { x := "ZOO_MY_ID"; return &x }(),
			Value: &sI,
		}, {
			Name:  func() *string { x := "ZOO_SERVERS"; return &x }(),
			Value: getZookeeperServerString(id),
		},
	}

	cmd.Hostname = "zookeeper" + strconv.Itoa(id) + "." + config.Domain

	d, _ := json.Marshal(&cmd)
	logrus.Debug("Scheduled Zookeeper: ", string(d))

	config.CommandChan <- cmd
	logrus.Info("Scheduled Zookeeper")
}

// The first run have to be in a right sequence
func initStartZookeeper() {
	if config.ZookeeperCount <= (config.ZookeeperMax - 1) {
		startZookeeper(config.ZookeeperCount)

		config.ZookeeperCount++
	}
}

// create the zookeeper connection string for every zookeeper container
func getZookeeperServerString(id int) *string {
	max := config.ZookeeperMax
	var server string
	for i := 0; i < max; i++ {
		sI := strconv.Itoa(i)
		if i == id {
			server += "server." + sI + "=0.0.0.0:2888:3888;2181 "
		} else {
			server += "server." + sI + "=zookeeper" + sI + "." + config.Domain + ":2888:3888;2181 "
		}
	}

	return &server
}

// create the zookeeper connection string for every kafka container
func createZookeeperServerString() {
	max := config.ZookeeperMax
	var server string
	for i := 0; i < max; i++ {
		sI := strconv.Itoa(i)
		server += "zookeeper" + sI + "." + config.Domain + ":2181, "
	}

	config.ZookeeperServers = server
}
