package mesos

import (
	"fmt"
	"sync/atomic"

	"../proto"
	cfg "../types"
	"git.aventer.biz/AVENTER/util"
)

func prepareTaskInfoExecuteCommand(agent *mesosproto.AgentID, cmd cfg.Command) ([]*mesosproto.TaskInfo, error) {

	newTaskID, _ := util.GenUUID()

	return []*mesosproto.TaskInfo{{
		Name: &cmd.TaskName,
		TaskId: &mesosproto.TaskID{
			Value: &newTaskID,
		},
		AgentId:   agent,
		Resources: defaultResources(),
		Command: &mesosproto.CommandInfo{
			Shell:       &cmd.Shell,
			Value:       &cmd.Command,
			Uris:        cmd.Uris,
			Environment: &cmd.Environment,
		},
	}}, nil
}

func prepareTaskInfoExecuteContainer(agent *mesosproto.AgentID, cmd cfg.Command) ([]*mesosproto.TaskInfo, error) {
	newTaskID := fmt.Sprint(atomic.AddUint64(&config.TaskID, 1))

	networkIsolator := "weave"

	contype := mesosproto.ContainerInfo_DOCKER.Enum()

	// Set Container Network Mode
	networkMode := mesosproto.ContainerInfo_DockerInfo_BRIDGE.Enum()

	if cmd.NetworkMode == "host" {
		networkMode = mesosproto.ContainerInfo_DockerInfo_HOST.Enum()
	}
	if cmd.NetworkMode == "none" {
		networkMode = mesosproto.ContainerInfo_DockerInfo_NONE.Enum()
	}
	if cmd.NetworkMode == "user" {
		networkMode = mesosproto.ContainerInfo_DockerInfo_USER.Enum()
	}
	if cmd.NetworkMode == "bridge" {
		networkMode = mesosproto.ContainerInfo_DockerInfo_USER.Enum()
	}

	// Save state of the task
	state := config.State[newTaskID]
	state.Command = cmd
	config.State[newTaskID] = state

	return []*mesosproto.TaskInfo{{
		Name: &cmd.TaskName,
		TaskId: &mesosproto.TaskID{
			Value: &newTaskID,
		},
		AgentId:   agent,
		Resources: defaultResources(),
		Command: &mesosproto.CommandInfo{
			Shell:       &cmd.Shell,
			Value:       &cmd.Command,
			Uris:        cmd.Uris,
			Environment: &cmd.Environment,
		},
		Container: &mesosproto.ContainerInfo{
			Type:    contype,
			Volumes: cmd.Volumes,
			Docker: &mesosproto.ContainerInfo_DockerInfo{
				Image:      &cmd.ContainerImage,
				Network:    networkMode,
				Privileged: &cmd.Privileged,
			},
			NetworkInfos: []*mesosproto.NetworkInfo{{
				Name: &networkIsolator,
			}},
		},
	}}, nil
}
