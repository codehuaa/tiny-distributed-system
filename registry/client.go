/**
 * @Author: Keven5
 * @Description:
 * @File:  client
 * @Version: 1.0.0
 * @Date: 2024/2/24 15:22
 */

package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func RegisterService(r Registration) error {
	// heartbeat handler binding
	heartbeatURL, err := url.Parse(r.HeartbeatURL)
	if err != nil {
		return err
	}
	http.HandleFunc(heartbeatURL.Path, func(w http.ResponseWriter, r *http.Request) {
		// we could do some other monitor actions here, such as Cpu, Memory etc.
		w.WriteHeader(http.StatusOK)
	})

	// update handler binding
	serviceUpdateURL, err := url.Parse(r.ServiceUpdatedURL)
	if err != nil {
		return err
	}
	http.Handle(serviceUpdateURL.Path, &serviceUpdateHandler{})

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(r)
	if err != nil {
		return err
	}
	resp, err := http.Post(ServicesUrl, "application/json", buf)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to register service. Registy service responded with code %v", resp.StatusCode)
	}

	return nil
}

func ShutdownService(url string) error {
	req, err := http.NewRequest(http.MethodDelete, ServicesUrl, bytes.NewBuffer([]byte(url)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/palin")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to deregister service. Registy service responded with code %v", resp.StatusCode)
	}
	return nil
}
