package storage

import "github.com/abhishekbotx/golang-restapi/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Students, error)
	GetStudents() ([]types.Students, error)
}