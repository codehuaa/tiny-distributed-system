/**
 * @Author: Keven5
 * @Description:
 * @File:  service
 * @Version: 1.0.0
 * @Date: 2024/2/24 11:06
 */

package log

import (
	stlog "log"
	"os"
)

var log *stlog.Logger

type fileLog string

func (fl fileLog) Write(data []byte) (int, error) {
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(data)
}

func write(msg string) {
	log.Printf("%v\n", msg)
}
