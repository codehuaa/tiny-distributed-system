/**
 * @Author: Keven5
 * @Description:
 * @File:  server
 * @Version: 1.0.0
 * @Date: 2024/2/24 11:45
 */

package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

const ServerPort = ":3000"
const ServicesUrl = "http://localhost" + ServerPort + "/services"

type registry struct {
	mutex         *sync.RWMutex
	registrations []Registration
}

// reg is the instance of registry.
// registrations is the container of service registration, including serviceName and serviceUrl
// mutex allows the concurrent security
var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.RWMutex),
}

// add a service to the registry service
func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	err := r.sendRequiredServices(reg)
	return err
}

func (r *registry) sendRequiredServices(reg Registration) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var p patch
	for _, serviceReg := range r.registrations {
		for _, reqService := range reg.RequiredServices {
			if serviceReg.ServiceName == reqService {
				p.Added = append(p.Added, patchEntry{
					Name: serviceReg.ServiceName,
					URL:  serviceReg.ServiceUrl,
				})
			}
		}
	}
	err := r.sendPatch(p, reg.ServiceUpdatedURL)
	if err != nil {
		return err
	}
	return nil
}

func (r *registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return err
	}
	return nil
}

// remove could remove the service from the registry service
func (r *registry) remove(url string) error {
	for i := range reg.registrations {
		if reg.registrations[i].ServiceUrl == url {
			r.mutex.Lock()
			reg.registrations = append(reg.registrations[:i], reg.registrations[i+1:]...)
			r.mutex.Unlock()
			return nil
		}
	}
	return fmt.Errorf("Service at URL %s is not found.", url)
}

// RegistryService implements the function ServeHTTP of http.Handler
type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("Request received")
	switch req.Method {
	// Post will add the registration
	case http.MethodPost:
		// dec :=
		var r Registration
		if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding service: %v with URL: %s\n", r.ServiceName, r.ServiceUrl)
		if err := reg.add(r); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	// Delete will remove the registration
	case http.MethodDelete:
		payload, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		url := string(payload)
		log.Printf("Removing service at URL %s", url)
		err = reg.remove(url)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	// Registry Service will not respond to other requests except Post & Delete
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
