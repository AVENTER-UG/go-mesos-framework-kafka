package mesos

import (
	mesosproto "../proto"
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
	taskId := *update.Status.GetTaskId().Value
	tmp := config.State[taskId]
	tmp.Status = update.Status
	config.State[taskId] = tmp

	return Call(msg)
}
