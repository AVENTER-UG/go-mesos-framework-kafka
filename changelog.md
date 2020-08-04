# Changelog

## dev

- Add persistent storage Support
- Change to go mod

## v0.0.2

- Add start mesos container
- Add start command
- Add persist framework info therefore the framework know what to do after a crash
- Add save task state
- Add start kafka and zookeeper
- Add persist task state
- Add multinode support
- Add container monitor to restart if a container failed
- Add RestCall to reflate missing zookeeper and kafka processes, for the case the monitoring can not find the problem
- Add RestCall scale up and down for zookeeper
- Add RestCall to kill a task
- Add Authentication
- Add Support to configure (non)SSL Support to Mesos
- Add Custom Domain to Zookeeper and Kafka to match Consul DNS
- Add Service Name ENV Variable to Zookeeper and Kafka to match Consul DNS
- Add Call_Suppress to tell mesos it does not send us offers until we ask
- Add default values for some env
- Add custom image name via env
