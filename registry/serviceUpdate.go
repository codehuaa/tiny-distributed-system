/**
 * @Author: Keven5
 * @Description:
 * @File:  serviceUpdate
 * @Version: 1.0.0
 * @Date: 2024/2/24 19:43
 */

package registry

import (
	"encoding/json"
	"log"
	"net/http"
)

type serviceUpdateHandler struct{}

func (suh *serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var p patch
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	prov.Update(p)
}
