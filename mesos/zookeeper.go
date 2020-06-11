package mesos

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync/atomic"

	mesosproto "../proto"

	cfg "../types"
	"github.com/sirupsen/logrus"
)

// SearchMissingZookeeper Check if all zookeepers are running. If one is missing, restart it.
func SearchMissingZookeeper() {
	if config.State != nil {
		for i := 0; i < config.ZookeeperMax; i++ {
			state := StatusZookeeper(i)
			if state != nil {
				if *state.Status.State != mesosproto.TaskState_TASK_RUNNING {
					logrus.Debug("Missing Zookeeper: ", i)
					GetZookeeperServerString(i)
					StartZookeeper(i)
				}
			}
		}
	}
}

// StatusZookeeper Get out Status of the given zookeeper ID
func StatusZookeeper(id int) *cfg.State {
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

// StartZookeeper is starting a zookeeper container with the given IDs
func StartZookeeper(id int) {
	newTaskID := atomic.AddUint64(&config.TaskID, 1)

	var cmd cfg.Command

	// before we will start a new zookeeper, we should be sure its not already running
	status := StatusZookeeper(id)
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

	networkIsolator := "weave"
	var hostport, containerport uint32
	hostport = 31210 + uint32(newTaskID)
	containerport = 2181
	protocol := "tcp"

	cmd.TaskID = newTaskID

	cmd.ContainerType = "DOCKER"
	cmd.ContainerImage = config.ImageZookeeper
	cmd.NetworkMode = "bridge"
	cmd.NetworkInfo = []*mesosproto.NetworkInfo{{
		Name: &networkIsolator,
	}}
	cmd.DockerPortMappings = []*mesosproto.ContainerInfo_DockerInfo_PortMapping{{
		HostPort:      &hostport,
		ContainerPort: &containerport,
		Protocol:      &protocol,
	}}

	cmd.Shell = false
	cmd.TaskName = "av_zookeeper" + strconv.Itoa(id)
	cmd.Hostname = "av_zookeeper" + strconv.Itoa(id) + config.ZookeeperCustomString + "." + config.Domain
	cmd.InternalID = id
	cmd.IsZookeeper = true
	cmd.Volumes = []*mesosproto.Volume{
		{
			HostPath:      func() *string { x := config.VolumeZookeeper + "/" + strconv.Itoa(id); return &x }(),
			ContainerPath: func() *string { x := "/datalog"; return &x }(),
			Mode:          mesosproto.Volume_RW.Enum(),
		},
	}
	sI := strconv.Itoa(id)
	cmd.Environment.Variables = []*mesosproto.Environment_Variable{
		{
			Name:  func() *string { x := "SERVICE_NAME"; return &x }(),
			Value: &cmd.TaskName,
		},

		{
			Name:  func() *string { x := "ZOO_MY_ID"; return &x }(),
			Value: &sI,
		},
		{
			Name:  func() *string { x := "ZOO_SERVERS"; return &x }(),
			Value: GetZookeeperServerString(id),
		},
	}

	d, _ := json.Marshal(&cmd)
	logrus.Debug("Scheduled Zookeeper: ", string(d))

	config.CommandChan <- cmd
	logrus.Info("Scheduled Zookeeper")
}

// The first run have to be in a right sequence
func initStartZookeeper() {
	if config.ZookeeperCount <= (config.ZookeeperMax - 1) {
		StartZookeeper(config.ZookeeperCount)
		config.ZookeeperCount++
	}
}

// GetZookeeperServerString create the zookeeper connection string for every zookeeper container
func GetZookeeperServerString(id int) *string {
	max := config.ZookeeperMax
	var server string
	for i := 0; i < max; i++ {
		sI := strconv.Itoa(i)
		if i == id {
			server += "server." + sI + "=0.0.0.0:2888:3888;2181 "
		} else {
			server += "server." + sI + "=av_zookeeper" + sI + config.ZookeeperCustomString + "." + config.Domain + ":2888:3888;2181 "
		}
	}

	return &server
}

// CreateZookeeperServerString create the zookeeper connection string for every kafka container
func CreateZookeeperServerString() {
	max := config.ZookeeperMax
	var server string
	for i := 0; i < max; i++ {
		sI := strconv.Itoa(i)
		server += "av_zookeeper" + sI + config.ZookeeperCustomString + "." + config.Domain + ":2181,"
	}

	// remote last char
	server = strings.TrimSuffix(server, ",")

	config.ZookeeperServers = server
}
