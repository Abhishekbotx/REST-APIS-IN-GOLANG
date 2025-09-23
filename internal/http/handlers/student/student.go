package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	// "github.com/abhishekbotx/golang-restapi/internal/http/handlers/student"
	"github.com/abhishekbotx/golang-restapi/internal/storage"
	"github.com/abhishekbotx/golang-restapi/internal/types"
	"github.com/abhishekbotx/golang-restapi/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc { //this means you can pass any struct which implements the storage interface
	slog.Info("Creating a student")
	return func(res http.ResponseWriter, req *http.Request) { //req is also of same interface as io.reader
		var student types.Students
		//.decode converts Json into govalse
		err := json.NewDecoder(req.Body).Decode(&student) //newdecode is of io.reader type// giving address of var student

		if errors.Is(err, io.EOF) { //When body is empty i.e-->EOF
			// response.WriteJson(res, http.StatusBadRequest, response.GeneralError(err))//either send this or custome as belows
			response.WriteJson(res, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty Body")))
			return
		}

		if err != nil {
			response.WriteJson(res, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//request validation

		if err := validator.New().Struct(student); err != nil { //line 19 var

			validateErrs := err.(validator.ValidationErrors) //typecasting to the type validation error requires
			response.WriteJson(res, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("user created successfully", slog.String("userId", fmt.Sprintf("%d", lastId)))

		if err != nil {
			response.WriteJson(res, http.StatusInternalServerError, err)
			return
		}

		// res.Write([]byte("welcome to students api"))
		response.WriteJson(res, http.StatusCreated, map[string]int64{"id": lastId})
	}
}


func GetById(storage storage.Storage) http.HandlerFunc { //this means you can pass any struct which implements the storage interface
	
	return func(res http.ResponseWriter, req *http.Request) { //req is also of same interface as io.reader
		id:=req.PathValue("id")
		slog.Info("getting student by id",slog.String("id:",id))
		intId, err:=strconv.ParseInt(id,10,64) //parseint returns int64 and error
		if err !=nil{
			response.WriteJson(res, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id format")))
			return
		}
		student,err:=storage.GetStudentById(intId) //converting string to int coz id in db is of int type

		if err !=nil{
			slog.Error("error getting user",slog.String("id",id))
			response.WriteJson(res, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(res, http.StatusOK, student)
	}

}

func GetStudents(storage storage.Storage) http.HandlerFunc { //this means you can pass any struct which implements the storage interface
	
	return func(res http.ResponseWriter, req *http.Request) { //req is also of same interface as io.reader
		slog.Info("getting all students")
		students,err:=storage.GetStudents() //converting string to int coz id in db is of int type
		// slog.Info("fetched all students here",students)
		if err !=nil{
			slog.Error("error getting all users")
			response.WriteJson(res, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		slog.Info("students to be sent in response", slog.Any("students", students))
		response.WriteJson(res, http.StatusOK, students)
	}
	

}