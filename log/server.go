/**
 * @Author: Keven5
 * @Description:
 * @File:  service
 * @Version: 1.0.0
 * @Date: 2024/2/24 17:05
 */

package log

import (
	"io"
	stlog "log"
	"net/http"
)

func Run(dest string) {
	log = stlog.New(fileLog(dest), "[go] ", stlog.LstdFlags)
}

func RegisterHandlers() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			msg, err := io.ReadAll(r.Body)
			if err != nil || len(msg) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			write(string(msg))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}
