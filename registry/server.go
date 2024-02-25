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
	"time"
)

const ServerPort = ":3000"
const ServicesUrl = "http://localhost" + ServerPort + "/services"

// registrations is the container of service registration, including serviceName and serviceUrl
// mutex allows the concurrent security
type registry struct {
	mutex         *sync.RWMutex
	registrations []Registration
}

// reg is the instance of registry.
var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.RWMutex),
}

// add a service to the registry service
func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	// reg get the required services if they were registered
	err := r.sendRequiredServices(reg)
	if err != nil {
		return err
	}
	// reg notify the services which required it
	r.notify(patch{
		Added: []patchEntry{
			{
				Name: reg.ServiceName,
				URL:  reg.ServiceUrl,
			},
		},
	})
	return nil
}

// sendRequiredServices will determine that whether there exists a service that reg required.
// Then it will send post request to the 'service updated url' to update the prov, which provide the providers server
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

// notify will notify the services which required the services in patch.Added
func (r *registry) notify(fullPatch patch) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, reg := range reg.registrations {
		// improve the efficiency of traversal by goroutine
		go func(reg Registration) {
			for _, requiredService := range reg.RequiredServices {
				p := patch{Added: make([]patchEntry, 0), Removed: make([]patchEntry, 0)}
				sendUpdateFlag := false
				for _, added := range fullPatch.Added {
					if added.Name == requiredService {
						p.Added = append(p.Added, added)
						sendUpdateFlag = true
					}
				}
				for _, removed := range fullPatch.Removed {
					if removed.Name == requiredService {
						p.Removed = append(p.Removed, removed)
						sendUpdateFlag = true
					}
				}
				if sendUpdateFlag {
					err := r.sendPatch(p, reg.ServiceUpdatedURL)
					if err != nil {
						log.Println(err)
						return
					}
				}
			}
		}(reg)
	}
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
			r.notify(patch{
				Removed: []patchEntry{
					{
						Name: reg.registrations[i].ServiceName,
						URL:  url,
					},
				},
			})
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

var once sync.Once

func SetupRegistryService() {
	once.Do(func() {
		go reg.heartbeat(3 * time.Second)
	})
}

// heartbeat sends requests to all registrations periodically to make sure that they are all connected and no problems
func (r *registry) heartbeat(freq time.Duration) {
	for {
		var wg sync.WaitGroup
		for _, reg := range r.registrations {
			wg.Add(1)
			go func(reg Registration) {
				defer wg.Done()
				success := true
				// retry for 3 times with 1 second interval
				for attempts := 0; attempts < 3; attempts++ {
					res, err := http.Get(reg.HeartbeatURL)
					if err != nil {
						log.Println(err)
					} else if res.StatusCode == http.StatusOK {
						log.Printf("Heartbeat check passed for %v", reg.ServiceName)
						if !success {
							err := r.add(reg)
							if err != nil {
								return
							}
						}
						break
					}
					log.Printf("Heartbeat failed for %v", reg.ServiceName)
					if success {
						success = false
						err := r.remove(reg.ServiceUrl)
						if err != nil {
							return
						}
					}
					time.Sleep(1 * time.Second)
				}
			}(reg)
		}
		wg.Wait()
		time.Sleep(freq)
	}
}
