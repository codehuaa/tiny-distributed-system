/**
 * @Author: Keven5
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2024/2/24 17:47
 */

package main

import (
	"context"
	"fmt"
	stlog "log"
	"tiny_distributed_system/grades"
	"tiny_distributed_system/log"
	"tiny_distributed_system/registry"
	"tiny_distributed_system/service"
)

func main() {
	host, port := "localhost", "6000"
	serviceAddr := fmt.Sprintf("http://%s:%s", host, port)

	r := registry.Registration{
		ServiceName:       registry.GradingService,
		ServiceUrl:        serviceAddr,
		RequiredServices:  []registry.ServiceName{registry.LogService},
		ServiceUpdatedURL: serviceAddr + "/services",
		HeartbeatURL:      serviceAddr + "/heartbeat",
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		grades.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatal(err)
	}

	if logProvider, err := registry.GetProviders(registry.LogService); err == nil {
		fmt.Printf("Logging service found at: %s\n", logProvider)
		log.SetClientLogger(logProvider, r.ServiceName)
	} else {
		fmt.Printf("Logging service not found. ")
	}

	<-ctx.Done()
	fmt.Println("Shutting down the grading service.")
}
