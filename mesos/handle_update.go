package mesos

import (
	"encoding/json"
	"io/ioutil"

	mesosproto "../proto"
	"github.com/sirupsen/logrus"
)

// HandleUpdate will handle the offers event of mesos
func HandleUpdate(event *mesosproto.Event) error {
	update := event.Update

	msg := &mesosproto.Call{
		Type: mesosproto.Call_ACKNOWLEDGE.Enum(),
		Acknowledge: &mesosproto.Call_Acknowledge{
			AgentId: update.Status.AgentId,
			TaskId:  update.Status.TaskId,
			Uuid:    update.Status.Uuid,
		},
	}

	// Save state of the task
	taskID := *update.Status.GetTaskId().Value
	tmp := config.State[taskID]
	tmp.Status = update.Status

	logrus.Debug("HandleUpdate: ", update.Status)

	switch *update.Status.State {
	case mesosproto.TaskState_TASK_FAILED:
		deleteOldTask(tmp.Status.TaskId)
	case mesosproto.TaskState_TASK_KILLED:
		deleteOldTask(tmp.Status.TaskId)
	case mesosproto.TaskState_TASK_LOST:
		deleteOldTask(tmp.Status.TaskId)
	}

	// Update Framework State File
	config.State[taskID] = tmp
	persConf, _ := json.Marshal(&config)
	ioutil.WriteFile(config.FrameworkInfoFile, persConf, 0644)

	return Call(msg)
}
