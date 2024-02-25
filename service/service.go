/**
 * @Author: Keven5
 * @Description:
 * @File:  service
 * @Version: 1.0.0
 * @Date: 2024/2/24 11:21
 */

package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"tiny_distributed_system/registry"
)

// Start the service
func Start(ctx context.Context, host, port string, reg registry.Registration, registerHandlerFunc func()) (context.Context, error) {
	registerHandlerFunc()
	ctx = startService(ctx, reg.ServiceName, host, port)
	err := registry.RegisterService(reg)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

// startService start the service of serviceName
func startService(ctx context.Context, serviceName registry.ServiceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	srv.Addr = host + ":" + port

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		fmt.Printf("%v started. Press any key to stop. \n", serviceName)
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		// remove the service registration
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		cancel()
	}()

	return ctx
}
