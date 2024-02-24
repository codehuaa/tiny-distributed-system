/**
 * @Author: Keven5
 * @Description:
 * @File:  student
 * @Version: 1.0.0
 * @Date: 2024/2/24 16:57
 */

package grades

import (
	"fmt"
	"sync"
)

type Student struct {
	ID        int
	FirstName string
	LastName  string
	Grades    []Grade
}

type Students []Student

var (
	students      Students
	studentsMutex sync.Mutex
)

func (s *Student) Average() float32 {
	var result float32
	for i := range s.Grades {
		result += s.Grades[i].Score
	}
	return result / float32(len(s.Grades))
}

func (ss Students) GetByID(id int) (*Student, error) {
	for i := range ss {
		if ss[i].ID == id {
			return &ss[i], nil
		}
	}
	return nil, fmt.Errorf("student with Id %d not found", id)
}
