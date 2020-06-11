package mesos

import (
	"strconv"

	mesosproto "../proto"
	cfg "../types"
)

func prepareTaskInfoExecuteContainer(agent *mesosproto.AgentID, cmd cfg.Command) ([]*mesosproto.TaskInfo, error) {
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
		networkMode = mesosproto.ContainerInfo_DockerInfo_BRIDGE.Enum()
	}

	// Save state of the new task
	newTaskID := strconv.Itoa(int(cmd.TaskID))
	tmp := config.State[newTaskID]
	tmp.Command = cmd
	config.State[newTaskID] = tmp

	if cmd.Shell == true {
		return []*mesosproto.TaskInfo{{
			Name: &cmd.TaskName,
			TaskId: &mesosproto.TaskID{
				Value: &newTaskID,
			},
			AgentId:   agent,
			Resources: defaultResources(cmd),
			Command: &mesosproto.CommandInfo{
				Shell:       &cmd.Shell,
				Value:       &cmd.Command,
				Uris:        cmd.Uris,
				Environment: &cmd.Environment,
			},
			Container: &mesosproto.ContainerInfo{
				Type:     contype,
				Volumes:  cmd.Volumes,
				Hostname: &cmd.Hostname,
				Docker: &mesosproto.ContainerInfo_DockerInfo{
					Image:        &cmd.ContainerImage,
					Network:      networkMode,
					PortMappings: cmd.DockerPortMappings,
					Privileged:   &cmd.Privileged,
				},
				NetworkInfos: cmd.NetworkInfo,
			},
		}}, nil
	} else {
		return []*mesosproto.TaskInfo{{
			Name: &cmd.TaskName,
			TaskId: &mesosproto.TaskID{
				Value: &newTaskID,
			},
			AgentId:   agent,
			Resources: defaultResources(cmd),
			Command: &mesosproto.CommandInfo{
				Shell:       &cmd.Shell,
				Uris:        cmd.Uris,
				Environment: &cmd.Environment,
			},
			Container: &mesosproto.ContainerInfo{
				Type:     contype,
				Volumes:  cmd.Volumes,
				Hostname: &cmd.Hostname,
				Docker: &mesosproto.ContainerInfo_DockerInfo{
					Image:        &cmd.ContainerImage,
					Network:      networkMode,
					PortMappings: cmd.DockerPortMappings,
					Privileged:   &cmd.Privileged,
				},
				NetworkInfos: cmd.NetworkInfo,
			},
		}}, nil
	}
}
