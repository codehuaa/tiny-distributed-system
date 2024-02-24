/**
 * @Author: Keven5
 * @Description:
 * @File:  client
 * @Version: 1.0.0
 * @Date: 2024/2/24 19:50
 */

package log

import (
	"bytes"
	"fmt"
	stlog "log"
	"net/http"
	"tiny_distributed_system/registry"
)

type clientLogger struct {
	urls []string
}

func (c *clientLogger) Write(data []byte) (int, error) {
	b := bytes.NewBuffer(data)
	res, err := http.Post(fmt.Sprintf("%s/log", c.urls), "text/plain", b)
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to send log message. Service respond")
	}
	return len(data), nil
}

func SetClientLogger(serviceURLs []string, clientService registry.ServiceName) {
	stlog.SetPrefix(fmt.Sprintf("[%v] - ", clientService))
	stlog.SetFlags(0)
	stlog.SetOutput(&clientLogger{serviceURLs})
}
