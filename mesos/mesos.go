package mesos

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	mesosproto "../proto"
	cfg "../types"

	"github.com/golang/protobuf/jsonpb"
	"github.com/sirupsen/logrus"
)

// Service include all the current vars and global config
var config *cfg.Config

// Marshaler to serialize Protobuf Message to JSON
var marshaller = jsonpb.Marshaler{
	EnumsAsInts: false,
	Indent:      " ",
	OrigName:    true,
}

// SetConfig set the global config
func SetConfig(cfg *cfg.Config) {
	config = cfg
}

// Subscribe to the mesos backend
func Subscribe() error {

	subscribeCall := &mesosproto.Call{
		FrameworkId: config.FrameworkInfo.Id,
		Type:        mesosproto.Call_SUBSCRIBE.Enum(),
		Subscribe: &mesosproto.Call_Subscribe{
			FrameworkInfo: &config.FrameworkInfo,
		},
	}
	body, _ := marshaller.MarshalToString(subscribeCall)
	logrus.Debug(body)
	client := &http.Client{}
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	req, _ := http.NewRequest("POST", "https://"+config.MesosMasterServer+"/api/v1/scheduler", bytes.NewBuffer([]byte(body)))
	req.SetBasicAuth(config.Username, config.Password)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		logrus.Fatal(err)
	}

	reader := bufio.NewReader(res.Body)

	line, _ := reader.ReadString('\n')
	bytesCount, _ := strconv.Atoi(strings.Trim(line, "\n"))

	for {
		// Read line from Mesos
		line, _ = reader.ReadString('\n')
		line = strings.Trim(line, "\n")
		// Read important data
		data := line[:bytesCount]
		// Rest data will be bytes of next message
		bytesCount, _ = strconv.Atoi((line[bytesCount:]))

		var event mesosproto.Event
		jsonpb.UnmarshalString(data, &event)
		logrus.Debug("Subscribe Got: ", event.GetType())

		if config.MesosStreamID != "" {
			initStartZookeeper()
			CreateZookeeperServerString()
			initStartKafka()
		}

		if event.Type != nil {
			switch *event.Type {
			case mesosproto.Event_SUBSCRIBED:
				logrus.Info("Subscribed")
				config.FrameworkInfo.Id = event.Subscribed.FrameworkId
				config.MesosStreamID = res.Header.Get("Mesos-Stream-Id")

				// Save framework info
				persConf, _ := json.Marshal(&config)
				ioutil.WriteFile(config.FrameworkInfoFile, persConf, 0644)

			case mesosproto.Event_UPDATE:
				logrus.Debug("Update", HandleUpdate(&event))
			case mesosproto.Event_HEARTBEAT:
			case mesosproto.Event_OFFERS:
				restartFailedContainer()
				logrus.Debug("Offers Returns: ", HandleOffers(event.Offers))
			default:
				logrus.Debug("DEFAULT EVENT: ", event.Offers)
			}
		} else {
			logrus.Error("Framework ID is outdated. Please remove the framework state file.")
		}
	}
}

// Call will send messages to mesos
func Call(message *mesosproto.Call) error {
	message.FrameworkId = config.FrameworkInfo.Id
	body, _ := marshaller.MarshalToString(message)

	client := &http.Client{}
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	req, _ := http.NewRequest("POST", "https://"+config.MesosMasterServer+"/api/v1/scheduler", bytes.NewBuffer([]byte(body)))
	req.SetBasicAuth(config.Username, config.Password)
	req.Header.Set("Mesos-Stream-Id", config.MesosStreamID)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	io.Copy(os.Stderr, res.Body)

	if res.StatusCode != 202 {
		return fmt.Errorf("Error %d", res.StatusCode)
	}

	return fmt.Errorf("Offer Accept %d", res.StatusCode)
}

// If the framework was restartet, we have to reconcile the task states.
func Reconcile() {
	var oldTasks []*mesosproto.Call_Reconcile_Task
	maxID := 0
	for _, t := range config.State {
		if t.Status != nil {
			oldTasks = append(oldTasks, &mesosproto.Call_Reconcile_Task{
				TaskId:  t.Status.TaskId,
				AgentId: t.Status.AgentId,
			})
			numericID, err := strconv.Atoi(t.Status.TaskId.GetValue())
			if err == nil && numericID > maxID {
				maxID = numericID
			}
		}
	}
	atomic.StoreUint64(&config.TaskID, uint64(maxID))
	Call(&mesosproto.Call{
		Type:      mesosproto.Call_RECONCILE.Enum(),
		Reconcile: &mesosproto.Call_Reconcile{Tasks: oldTasks},
	})
}

// Restart failed zookeeper container
func restartFailedContainer() {
	if config.State != nil {
		for _, element := range config.State {
			if element.Status != nil {
				logrus.Debug("restartFailedContainer: ", *element.Status.State, element.Status.TaskId)
				switch *element.Status.State {
				case mesosproto.TaskState_TASK_FAILED:
					if element.Command.IsZookeeper == true {
						logrus.Info("RestartZookeeper: ", element.Status.TaskId)
						StartZookeeper(element.Command.InternalID)
					}
					if element.Command.IsKafka == true {
						logrus.Info("RestartKafka: ", element.Status.TaskId)
						StartKafka(element.Command.InternalID)
					}
					deleteOldTask(element.Status.TaskId)
				case mesosproto.TaskState_TASK_KILLED:
					deleteOldTask(element.Status.TaskId)
				case mesosproto.TaskState_TASK_LOST:
					deleteOldTask(element.Status.TaskId)
				case mesosproto.TaskState_TASK_ERROR:
					deleteOldTask(element.Status.TaskId)
				}
			}
		}
	}
}

// Delete Failed Tasks from the config
func deleteOldTask(taskID *mesosproto.TaskID) {
	copy := make(map[string]cfg.State)

	if config.State != nil {
		for _, element := range config.State {
			if element.Status != nil {
				tmpID := *element.Status.GetTaskId().Value
				if element.Status.TaskId != taskID {
					copy[tmpID] = element
				} else {
					logrus.Debug("Delete Task from config: ", tmpID)
				}
			}
		}

		config.State = copy
	}
}

// Kill a Task with the given taskID
func Kill(taskID string) error {
	task := config.State[taskID]

	logrus.Debug("Kill task %s [%#v]", taskID, task)

	return Call(&mesosproto.Call{
		Type: mesosproto.Call_KILL.Enum(),
		Kill: &mesosproto.Call_Kill{
			TaskId:  task.Status.TaskId,
			AgentId: task.Status.AgentId,
		},
	})
}
