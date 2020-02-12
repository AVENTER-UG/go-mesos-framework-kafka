# go-mesos-framework-kafka

<a href="https://liberapay.com/mesos" target="_new"><img src="https://img.shields.io/liberapay/receives/AVENTER.svg?logo=liberapay"></a>


Dies ist ein Kafka Framework für Apache Mesos

## Vorraussetzung

Dieses Framework ist aktuell so erstellt, dass es MESOS mit SSL Verschlüsselung und Authentication benötigt.

## Framework starten

```Bash

export FRAMEWORK_USER="root"
export FRAMEWORK_NAME="kafkaframework"
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

go run init.go app.go
```

Dies startet das Framework. Es wird sich an den Mesos Master anmelden. Nach wenigen Sekunden kann man "kafkaframework" als Eintrag in der Mesos UI sehen. Sobald das Framework erfolgreich gestartet wurde, startet es Zookeeper. Ist Zookeeper erfolgreich gestartet, werden die Kafka Container gestartet.

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
