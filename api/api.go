package api

import (

	//"encoding/json"

	"github.com/gorilla/mux"
	//"io/ioutil"
	"net/http"

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
	rtr.HandleFunc("/v0/container/{taskID}", V0StatusContainer).Methods("GET")
	rtr.HandleFunc("/v0/zookeeper/scale/{count}", V0ScaleZookeeper).Methods("GET")
	rtr.HandleFunc("/v0/zookeeper/reflate", V0ReflateZookeeper).Methods("GET")
	rtr.HandleFunc("/v0/kafka/scale/{count}", V0ScaleKafka).Methods("GET")
	rtr.HandleFunc("/v0/kafka/reflate", V0ReflateKafka).Methods("GET")
	rtr.HandleFunc("/v0/task/kill/{id}", V0KillTask).Methods("GET")

	return rtr
}

// CheckAuth will check if the token is valid
func CheckAuth(r *http.Request, w http.ResponseWriter) bool {
	// if no credentials are configured, then we dont have to check
	if config.Credentials.Username == "" || config.Credentials.Password == "" {
		return true
	}

	username, password, ok := r.BasicAuth()

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	if username == config.Credentials.Username && password == config.Credentials.Password {
		w.WriteHeader(http.StatusOK)
		return true
	}

	w.WriteHeader(http.StatusUnauthorized)
	return false
}
