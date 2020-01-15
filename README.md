# Benchmark System

## Requirement

- Support load test as code
- Support simulate multiple protocol (gRPC, HTTP,...)
- Flexible configuration
- UI to client can: config task, submit task, view report
- Persistent data
  
## Specification

### Flowchart

```plantuml
@startuml
participant Client
participant LocustController
participant Locust
participant Worker
participant Database

== ExecuteTest ==

Client -> LocustController: upload config
LocustController -> LocustController: parsing config
LocustController -> Database: save config and task
LocustController -> Locust: execute benchmark task
Locust -> Worker: run task
Worker --> Locust: recorded metrics

== QueryMetrics ==

LocustController -> Locust: get recorded metrics
Locust --> LocustController: metrics
LocustController -> Database: save metrics

@enduml
```

## Planning

### Sprint 1 (v0.1.1)

Overview:

- Basic benchmark tool.
- CLI tool features:
  - Connect Locust and work as slave.
  - Auto generate data and request server by specific proto file.
  - Test configuration: hatch rate, concurrent user, target service.
  - Render report to LocustUI.

Deliveriable:

- Source code
- Documentations: sequence diagram , usecase, internal service, technical stack.

[Detail](src/README.md)

### Sprint 2

`TBD`