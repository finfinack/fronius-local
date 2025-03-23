# fronius-local

Tool to locally pull data from Fronius and send it to an MQTT topic.

## Config

```json
{
  "mqtt": {
    "broker": "tcp://<mqtt-broker>:1883",
    "username": "<username>",
    "password": "<password>",
    "topic": "<topic>"
  },
  "fronius": {
    "backend": "http://<fronius-ip>",
    "username": "customer",
    "password": "<password>"
  },
  "interval": "1h"
}
```

## Run on k8s

- Check Helm chart in `deploy` folder
- Ensure a configmap `fronius-local-config` exists, containing a key `prod` which contains a config as described above.

  ```bash
  kubectl create configmap fronius-local-config --from-file=prod=./config.json
  ```
