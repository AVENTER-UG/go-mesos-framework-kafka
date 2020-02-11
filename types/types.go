package types

import mesosproto "../proto"

// Config is a struct of the framework configuration
type Config struct {
	FrameworkPort         string
	FrameworkBind         string
	FrameworkUser         string
	FrameworkName         string
	FrameworkInfo         mesosproto.FrameworkInfo
	FrameworkInfoFile     string
	FrameworkInfoFilePath string
	Principal             string
	Username              string
	Password              string
	MesosMasterServer     string
	MesosSSL              bool
	MesosStreamID         string
	TaskID                uint64
	SSL                   bool
	LogLevel              string
	MinVersion            string
	AppName               string
	EnableSyslog          bool
	Hostname              string
	Listen                string
	CommandChan           chan Command `json:"-"`
	State                 map[string]State
	Domain                string
	ZookeeperServers      string
	ZookeeperCount        int
	ZookeeperMax          int
	KafkaCount            int
	KafkaMax              int
	ResCPU                float64
	ResMEM                float64
	Credentials           UserCredentials
}

// Command is a chan which include all the Information about the started tasks
type Command struct {
	ContainerImage string                        `json:"container_image,omitempty"`
	ContainerType  string                        `json:"container_type,omitempty"`
	TaskName       string                        `json:"task_name,omitempty"`
	Command        string                        `json:"command,omitempty"`
	Hostname       string                        `json:"hostname,omitempty"`
	Privileged     bool                          `json:"privileged,omitempty"`
	NetworkMode    string                        `json:"network_mode,omitempty"`
	Volumes        []*mesosproto.Volume          `protobuf:"bytes,2,rep,name=volumes" json:"volumes,omitempty"`
	Shell          bool                          `protobuf:"varint,6,opt,name=shell,def=1" json:"shell,omitempty"`
	Uris           []*mesosproto.CommandInfo_URI `protobuf:"bytes,1,rep,name=uris" json:"uris,omitempty"`
	Environment    mesosproto.Environment        `protobuf:"bytes,2,opt,name=environment" json:"environment,omitempty"`
	Arguments      []string                      `protobuf:"bytes,7,rep,name=arguments" json:"arguments,omitempty"`
	Executor       mesosproto.ExecutorInfo
	InternalID     int
	IsZookeeper    bool
	IsKafka        bool
}

// State will have the state of all tasks stated by this framework
type State struct {
	Command Command                `json:"command"`
	Status  *mesosproto.TaskStatus `json:"status"`
}

// UserCredentials - The Username and Password to authenticate against this framework
type UserCredentials struct {
	Username string
	Password string
}
