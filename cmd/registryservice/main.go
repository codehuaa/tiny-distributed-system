/**
 * @Author: Keven5
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2024/2/24 14:57
 */

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"tiny_distributed_system/registry"
)

func main() {
	http.Handle("/services", &registry.RegistryService{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var srv http.Server
	srv.Addr = registry.ServerPort

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		fmt.Printf("Registry service started. Press any key to stop. \n")
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		cancel()
	}()

	<-ctx.Done()
	fmt.Println("Shutting down the registry service")
}
