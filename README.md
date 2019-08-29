# Service-Prober

Sync kubernetes service to cloudprober service( run inside the same container )

## Usage

```
# Service
apiVersion: v1
kind: Service
metadata:
  name: fs
  namespace: yunwei
  annotations:
    prober.haodai.net/enable: "true"
```

> remove annotation will remove the probe entry.

## Configure path

service annotations:

| annotation | example value |
|---|--|
|prober.haodai.net/enable: | "true" |
|prober.haodai.net/path: | "/healthz" |


## Cloudprober status page

![cloudprober status](doc/status.png)

## Demo grafana

![grafana](doc/grafana.png)

[grafana dashboard](grafana-dashboard.json)

> we use prometheus auto service discovery to let prometheus pull metrics

example config ( depends on your prometheus config )

```
apiVersion: v1
kind: Service
metadata:
  name: service-prober
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/port: "9313"
```