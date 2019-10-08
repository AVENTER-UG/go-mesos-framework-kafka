package mesos

import (
	"github.com/sirupsen/logrus"

	"../proto"
)

// HandleUpdate will handle the offers event of mesos
func HandleUpdate(event *mesosproto.Event) error {

	taskStatus := event.GetUpdate().GetStatus()

	taskID := *taskStatus.TaskId.Value

	state := config.State[taskID]
	state.Status = taskStatus.GetState().String()
	config.State[taskID] = state

	logrus.Debug("HandleUpate cmd: ", state)

	return nil

}
