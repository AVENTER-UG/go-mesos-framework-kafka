# go-mesos-framework-kafka

[![Donate](https://img.shields.io/liberapay/receives/AVENTER.svg?logo=liberapay)](https://liberapay.com/mesos)
[![Support Chat](https://img.shields.io/static/v1?label=Chat&message=Support&color=brightgreen)](https://riot.im/app/#/room/#support:matrix.aventer.biz)

Dies ist ein Kafka Framework für Apache Mesos

## Voraussetzung

- Apache Mesos ab 1.6.0
- Mesos mit SSL und Authentication ist Optional
- Persistent Storage

## Framework starten

Mit den folgenden Umgebungsvariablen kann das Framework konfiguriert werden. Nach dem Starten wird es sich an den Mesos Master anmelden und erscheint als "kafkaframework" in der Mesos UI. Sobald sich das Framework erfolgreich initialisiert hat, startet es Zookeeper und anschließend Kafka.



```Bash
export FRAMEWORK_USER="root"
export FRAMEWORK_NAME="kafkaframework"
export FRAMEWORK_PORT="10000"
export FRAMEWORK_ROLE="kafka"
export FRAMEWORK_STATEFILE_PATH="/tmp"
export MESOS_PRINCIPAL="<mesos_principal>"
export MESOS_USERNAME="<mesos_user>"
export MESOS_PASSWORD="<mesos_password>"
export MESOS_MASTER="<mesos_master_server>:5050"
export LOGLEVEL="DEBUG"
export DOMAIN="weave.local"
export ZOOKEEPER_COUNT=1
export KAFKA_COUNT=3
export RES_CPU=0.1
export RES_MEM=3200
export AUTH_PASSWORD="password"
export AUTH_USERNAME="user"
export MESOS_SSL="true"
export ZOOKEEPER_CUSTOM_DOMAIN=""
export KAFKA_CUSTOM_DOMAIN=""
export IMAGE_KAFKA="confluentinc/cp-kafka:5.4.1"
export IMAGE_ZOOKEEPER="zookeeper:3.5.7"
export VOLUME_DRIVER="local"
export VOLUME_KAFKA="/tmp/kafka1,/tmp/kafka2,/tmp/kafka3"
export VOLUME_ZOOKEEPER="/tmp/zookeeper1,/tmp/zookeeper2,/tmp/zookeeper3"

go run init.go app.go
```

![Kafka Framework in Mesos](kafka_mesos.gif)

Faellt ein Container aus, wird dieser neugestartet.

![Kafka Framework in Mesos](kafka_mesos1.gif)

## Task Status Abfragen

Um den Status eines Tasks über das Framework abzufragen, folgendes Kommando verwenden:

```Bash
curl -X GET 127.0.0.1:10000/v0/container/<taskId> -d 'JSON'  | jq
```

## Fehlende Kafka oder Zookeeper Starten

Sollte aus bestimmten Gründen der Healthcheck Status im Framework nicht mit der Realität übereinstimmen, kann über den nachfolgenden Aufruf erzwungen werden die fehlenden Container zu starten.

```Bash
curl -X GET 127.0.0.1:10000/v0/<kafka|zookeeper>/reflate -d 'JSON'
```

## Kafka oder Zookeeper skalieren

Um Kafka oder Zookeeper im Betrieb zu skalieren, wird dem Framework die zu laufende Anzahl an Containern angegeben. Soll der Zookeeper also drei mal laufen, muss als <count> eine "3" angegeben werden. Beim Scaledown wird der zuletzt hinzugefügte Container entfernt.

```Bash
curl -X GET 127.0.0.1:10000/v0/<kafka|zookeeper>/scale/<count> -d 'JSON'
```

## Task killen

Sollte es notwendig sein einen Task zu benden, erfolgt dies mit dem nachfolgenden aufruf. Der beendete Container wird nicht automatisch neu gestartet.

```Bash
curl -X GET 127.0.0.1:10000/v0/task/kill/<taskId> -d 'JSON'
```
