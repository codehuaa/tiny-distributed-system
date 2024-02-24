/**
 * @Author: Keven5
 * @Description:
 * @File:  mockData
 * @Version: 1.0.0
 * @Date: 2024/2/24 16:58
 */

package grades

func init() {
	students = []Student{
		{
			ID:        1,
			FirstName: "Zhang",
			LastName:  "San",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 85,
				},
				{
					Title: "Quiz 2",
					Type:  GradeQuiz,
					Score: 90,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 90,
				},
			},
		},
		{
			ID:        2,
			FirstName: "Li",
			LastName:  "Si",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 100,
				},
				{
					Title: "Quiz 2",
					Type:  GradeQuiz,
					Score: 100,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 100,
				},
			},
		},
		{
			ID:        3,
			FirstName: "Wang",
			LastName:  "Wu",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 55,
				},
				{
					Title: "Quiz 2",
					Type:  GradeQuiz,
					Score: 48,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 52,
				},
			},
		},
	}
}
