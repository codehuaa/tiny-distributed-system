/**
 * @Author: Keven5
 * @Description:
 * @File:  grades
 * @Version: 1.0.0
 * @Date: 2024/2/24 16:48
 */

package grades

type GradeType string

const (
	GradeQuiz = GradeType("Quiz")
	GradeTest = GradeType("Test")
	GradeExam = GradeType("Exam")
)

type Grade struct {
	Title string
	Type  GradeType
	Score float32
}
