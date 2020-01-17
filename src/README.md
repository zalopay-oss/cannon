# Zalopay-Cannon v0.1.1

![version](https://img.shields.io/badge/version-0.1.1-red) [![issues](https://img.shields.io/badge/open%20issues-0-orange)]() [![contributors](https://img.shields.io/badge/contributors-2-blue)]()

## Introduction

ZaloPay Cannon is a benchmark system for ZaloPay's internal service. The aim is to build a multi-tennant system which provides intuitive UI/UX for users to submit tasks and perform benchmark.  

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

## Run

- Edit configuration in file `configs/config.yaml`

```yaml
# Locust config
Locust: "http://0.0.0.0:7000/" # Locust address
NoConns: 1 # Number of connections | number of slaves
NoWorkers: 80 # Number of users
HatchRate: 10 # Hatch rate

# Influx DB config
Bucket: "benchmark-results" # Bucket's name
Origin: "zlp" # Influxdb Origin
DatabaseAddr: "http://0.0.0.0:9999" # Database address
Token: "e_jl06gwSsAmwStymP1hrSp3_-l8s56QFT9jzklJ_B_uTwu6L4h1BtjFRoYk3LgsDGKl562X8msWwbaQN5llQg==" # InfluxDB Token

# GRPC config
GRPCPort: 1234
GRPCHost: "localhost"
Service: "service.KeyValueStoreService.Connect" #  Service's name
Proto: "path/to/name.proto" #link to proto
```

- Make sure Locust, InfluxDB and gRPC server are running.
- Run:

```bash
./run.sh
```
