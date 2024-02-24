/**
 * @Author: Keven5
 * @Description:
 * @File:  provider
 * @Version: 1.0.0
 * @Date: 2024/2/24 19:31
 */

package registry

import (
	"fmt"
	"sync"
)

type providers struct {
	services map[ServiceName][]string
	mutex    *sync.RWMutex
}

var prov = providers{
	services: make(map[ServiceName][]string),
	mutex:    new(sync.RWMutex),
}

func (p *providers) Update(pat patch) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	// add the services
	for _, patchEntry := range pat.Added {
		// if the service doesn't exist in prov
		if _, ok := p.services[patchEntry.Name]; !ok {
			p.services[patchEntry.Name] = make([]string, 0)
		}
		p.services[patchEntry.Name] = append(p.services[patchEntry.Name], patchEntry.URL)
	}

	// remove the services
	for _, patchEntry := range pat.Removed {
		if providerURLs, ok := p.services[patchEntry.Name]; ok {
			for i := range providerURLs {
				if providerURLs[i] == patchEntry.URL {
					p.services[patchEntry.Name] = append(p.services[patchEntry.Name][:i], p.services[patchEntry.Name][i+1:]...)
				}
			}
		}
	}
}

// get the service URLs from the name
func (p *providers) get(name ServiceName) ([]string, error) {
	providers, ok := p.services[name]
	if !ok {
		return nil, fmt.Errorf("no providers available for service %v", name)
	}
	return providers, nil
}

func GetProviders(name ServiceName) ([]string, error) {
	return prov.get(name)
}
