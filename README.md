# Tiny Distributed System

![Static Badge](https://img.shields.io/badge/go-b)
![Static Badge](https://img.shields.io/badge/microservice-red)
![Static Badge](https://img.shields.io/badge/distributed-system-blue)

This is a very simple distributed system implemented with Golang. 

In this system, we realize the service startup, the service registry, the service notification, heartbeat and so on. 

> Before you reading the code, make sure that you have learned about the simple using of following packages:
> - content
> - goroutine
> - net/http
> - sync.Mutex
> - channel

## Overview
(TODO)

## Structure
(TODO)

## Run
### 1. start up
Firstly, start the registry service:
```shell
cd .\cmd\registryservice\
go run .\main.go
```

Secondly, open a new terminal and start the logging service, which provide service for grading_service:
```shell
cd .\cmd\logservice\
go run .\main.go
```

Thirdly, open a new terminal, then start the grading service:
```shell
cd .\cmd\gradingservice\
go run .\main.go
```


Then, as the following figure shown, you could see the `registry_service` has found register the `logging_service` and `grading_service`. 
And `grading_service` found the required service `logging_service`. 
Besides, `registry_service` has started the heartbeat to check the status of other services. 

![running.png](pics%2Frunning.png)

### 2. shut down
You could shut down the service by **pressing any key and Enter**. 

If you shut down the other service rather than `registry_service`, you could see `Removing service at URL ...`

![shutdown1.png](pics%2Fshutdown1.png)

And you can simulate the breakdown of service by **enter `Ctrl+c`**. 

Then you will find that `registry_service` heartbeat reports that failure.
`registry_service` will retry to connect the service at backend for 3 times. 
You could see it in `registry/server.go`

