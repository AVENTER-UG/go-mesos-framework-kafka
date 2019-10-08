# go-mesos-framework-basis

Dies ist die Basis für Mesos Frameworks.

## Vorraussetzung

Dieses Basis Framework ist aktuell so erstellt, dass es MESOS mit SSL Verschlüsselung und Authentication benötigt.

## Framework starten

```Bash

export FRAMEWORK_USER="root"
export FRAMEWORK_NAME="test_framework"
export MESOS_PRINCIPAL="<mesos_principal>"
export MESOS_USERNAME="<mesos_user>"
export MESOS_PASSWORD="<mesos_password>"
export MESOS_MASTER="<mesos_master_server>:5050"


go run init.go app.go
```

Dies startet das Framework. Es wird sich an den Mesos Master anmelden. Nach wenigen Sekunden kann man "test_framework" als Eintrag in der Mesos UI sehen. Gleichzeitig öffnet das Framework einen Port auf 10000 auf der Maschine auf dem das Framework gestartet wurde.

## Task Starten

### Command

```Bash
curl -X POST 127.0.0.1:10000/v0/command/start\?cmd\=python%20-m%20SimpleHTTPServer%209033
```

### Mesos Container

Um einen Mesos Container zu starten, muss man der nachfolgenden Aufruf angepasst werden. "Value" bekommt dabei eine URL von dem aus ein Binary heruntergeladen wird. Das Binary wird dann, über "Command" aufgerufen.

```Bash
 curl -X POST 127.0.0.1:10000/v0/container/start -d '{ "command": "./test", "uris": [{ "value": "https://<URL>/test", "extract": false, "executable": true, "cache": false }]}'
```

Auf einem Mesos Agent wird man nun einen entsprechenden Prozess erkennen können.
