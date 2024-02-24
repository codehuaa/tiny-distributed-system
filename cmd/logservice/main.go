/**
 * @Author: Keven5
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2024/2/24 11:21
 */

package main

import (
	"context"
	"fmt"
	stlog "log"
	"tiny_distributed_system/log"
	"tiny_distributed_system/registry"
	"tiny_distributed_system/service"
)

func main() {
	log.Run("./distributed.log")
	host, port := "localhost", "4000"
	serviceAddr := fmt.Sprintf("http://%s:%s", host, port)
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		registry.Registration{
			ServiceName:       registry.LogService,
			ServiceUrl:        serviceAddr,
			RequiredServices:  make([]registry.ServiceName, 0),
			ServiceUpdatedURL: serviceAddr + "/services",
		},
		log.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatalln(err)
	}

	<-ctx.Done()

	fmt.Println("Shutting down the log service.")
}
