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
		logrus.Debug("Subscribe Got RAW: ", data)
		logrus.Info("Subscribe Got: ", event.GetType())

		switch *event.Type {
		case mesosproto.Event_SUBSCRIBED:
			logrus.Info("Subscribed")
			config.FrameworkInfo.Id = event.Subscribed.FrameworkId
			config.MesosStreamID = res.Header.Get("Mesos-Stream-Id")

			// Save framework info
			persConf, _ := json.Marshal(&config)
			ioutil.WriteFile(config.FrameworkInfoFile, persConf, 0644)

			// Start the Zookeeper Container
			startZookeeper(1, 1)
			//startZookeeper(2, 3)
			//startZookeeper(3, 3)
			createZookeeperServerString(1)

			// Start the Kafka Container
			//startKafka(1, 3)
			//startKafka(2, 3)
			//startKafka(3, 3)
		case mesosproto.Event_UPDATE:
			logrus.Info("Update", HandleUpdate(&event))
		case mesosproto.Event_HEARTBEAT:
			logrus.Info("Heartbeat")
		case mesosproto.Event_OFFERS:
			logrus.Info("Offers Returns: ", HandleOffers(event.Offers))
		default:
			logrus.Info("DEFAULT EVENT: ", event.Offers)
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
