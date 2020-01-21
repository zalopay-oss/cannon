# ZaloPay Cannon

![version](https://img.shields.io/badge/version-0.1.1-red) [![issues](https://img.shields.io/badge/open%20issues-0-orange)]() [![contributors](https://img.shields.io/badge/contributors-3-blue)]()

<img src="images/cannon-logo.jpg" alt="logo" width="600"/>

## Introduction

ZaloPay Cannon is a benchmark system for ZaloPay's internal services. The aim is to build a multi-tennant system which provides intuitive UI/UX for users to submit tasks and perform benchmark.  

## Architecture

![architecture](images/architecture.png)  

## Features

- Benchmark gRPC service with given proto.  
- Distributed testing: run tests on multiple slaves.  
- Automatically generate input data.  
- Visualize data: metrics data is stored and visualized using InfluxDB.  
- Only support Unary RPCs.  

## Requirements

- Golang 1.13.1
- Locust
- Influxdb 2.0
- Python 3.7.3  

## Usage

- Make sure Locust, InfluxDB and gRPC server are running.

```bash
Usage:
  cannon run [flags]

Flags:
  -c, --config string        path of config file
  -r, --hatchRate int        config hatch rate (users spawned/second) (default 10)
  -h, --help                 help for run
  -H, --host string          target gRPC host (default "localhost")
      --locust-host string   host of locust master (default "localhost")
      --locust-port int      port of locust master (default 5557)
      --locust-web string    locust web target (default "http://localhost:8089")
  -m, --method string        method name
  -w, --no-workers int       number of workers to simulate (default 10)
  -P, --port int             target gRPC port (default 8000)
  -p, --proto string         path of proto file

```

- Cannon example config file

```yaml
# Locust Config
LocustWebPort: "http://0.0.0.0:7000/"
LocustHost: "127.0.0.1"
LocustPort: 5557
NoWorkers: 80 # Number of connections
HatchRate: 10 # Hatch rate

# InfluxDB Config
IsPersistent: "true" # turn on if you want save metrics in influxDB
Bucket: "benchmark-results" # InfluxDB bucket's name
Origin: "zlp-osss"
DatabaseAddr: "http://0.0.0.0:9999"
Token: "egc6_K6V0pCmEwIahIzmnoneommTcsa7TS5XtmcSBnR9VeX31dMsRJ_STN-bUqOwWW77vPiU0aM9RGMQFwxT-A=="

# gRPC benchmark target
GRPCHost: "localhost"
GRPCPort: 4770
Proto: "./proto-name.proto"
Method: "serviceName.methodName"
```

- Monitor at dashboard locust

<img src="images/dashboard-example.png" alt="logo" width="800"/>

## Example

Read the [example](example/README.md)

## Roadmap

Read the [roadmap](docs/ROADMAP.md)

## Acknowledgements

- Thanks to @anhld2, which served as an inspiration and guide in building this project.
- Special thanks to @thinhda, @tranndc and @quyenpt for their work on making component-based theming a reality.
