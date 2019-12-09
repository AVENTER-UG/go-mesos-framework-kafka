package api

import (

	//"encoding/json"

	"github.com/gorilla/mux"
	//"io/ioutil"

	cfg "../types"
)

// Service include all the current vars and global config
var config *cfg.Config

// SetConfig set the global config
func SetConfig(cfg *cfg.Config) {
	config = cfg
}

// Commands is the main function of this package
func Commands() *mux.Router {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/v0/cmd/start", V0StartCommand).Methods("POST")
	rtr.HandleFunc("/v0/container/start", V0StartContainer).Methods("POST")
	rtr.HandleFunc("/v0/container/{taskID}", V0StatusContainer).Methods("GET")
	rtr.HandleFunc("/v0/zookeeper/scale/{count}", V0ScaleZookeeper).Methods("GET")
	rtr.HandleFunc("/v0/zookeeper/reflate", V0ReflateZookeeper).Methods("GET")
	rtr.HandleFunc("/v0/kafka/reflate", V0ReflateKafka).Methods("GET")
	rtr.HandleFunc("/v0/task/kill/{id}", V0KillTask).Methods("GET")

	return rtr
}
