package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/abhishekbotx/golang-restapi/internal/config"
	"github.com/abhishekbotx/golang-restapi/internal/types"
	_ "github.com/mattn/go-sqlite3" // sqlite driver
	// "golang.org/x/tools/go/analysis/passes/defers"
)

// Sqlite wrapper
type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	// Create table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students( 
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		age INTEGER
	)`)
	if err != nil {
		return nil, err
	}
	slog.Info("table created or already exists")

	// // Insert test data (ignore if duplicate by checking later)
	// _, err = db.Exec(`INSERT INTO students(name, email, age) VALUES (?, ?, ?)`,
	// 	"Test User", "test@example.com", 20)
	// if err != nil {
	// 	slog.Error("failed to insert", "err", err)
	// }

	// ✅ Correct way to fetch data
	rows, err := db.Query(`SELECT id, name, email, age FROM students`)
	if err != nil {
		slog.Error("failed to fetch students", "err", err)
		return nil, err
	}
	defer rows.Close()

	slog.Info("till here")

	for rows.Next() {
		var id int
		var name, email string
		var age int
		if err := rows.Scan(&id, &name, &email, &age); err != nil {
			slog.Error("failed to scan row", "err", err)
			continue
		}
		slog.Info("student record",
			"id", id,
			"name", name,
			"email", email,
			"age", age,
		)
	}

	return &Sqlite{
		Db: db,
	}, nil
}
//remember:Any struct with a method matching this signature implements the interface.
//why pointer? coz even in the struct we are receiving it
func (s *Sqlite) CreateStudent (name string, email string, age int) (int64,error){ //now this method is attached to sqlite struct ln:12
	//Now here the sqlite struct is implementing the 
	stmt,err:=s.Db.Prepare("INSERT INTO students (name,email,age) VALUES(?,?,?)")
	if err !=nil{
		return 0, err
	}
	defer stmt.Close()//function close ke baad automatically staatmnet close ho jayega

	result, err:=stmt.Exec(name,email,age) //res2 
	if err !=nil{
		return 0,err
	}

	lastId, err:=result.LastInsertId()
	if err !=nil{
		return 0,err
	}

	return lastId,nil

	
	// return 0,nil //for integer and error
}

/*
	Storage interface defines the contract.
	Sqlite struct implements the contract.
	Handler depends on the interface, not the concrete type.
*/

func (s *Sqlite) GetStudentById(id int64) (types.Students, error) {
	stmt,err:=s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")
	if err !=nil{
		return  types.Students{}, err //returning empty struct
	}

	defer stmt.Close()

	var student types.Students
	//passing reference of student to scan method so that it can fill the values in it
	err=stmt.QueryRow(id).Scan(&student.Id,&student.Name,&student.Email,&student.Age) //id is coming from parameter

	if err !=nil{
		if err== sql.ErrNoRows{
			return types.Students{}, fmt.Errorf("student with id %d not found", id)
		}
		return types.Students{}, fmt.Errorf("query error: %w", err) //wrapping error with fmt.Errorf
	}

	return student, nil

}

func(s *Sqlite) GetStudents()([]types.Students,error){
	stmt,err:=s.Db.Prepare("SELECT id,name,email,age FROM students")
	if err!=nil{
		return nil,err
	}

	defer stmt.Close()

	rows, err := stmt.Query() //exec query

	if err!=nil{
		return nil, err
	}

	defer rows.Close()

	var students []types.Students

	for rows.Next(){//next returns true if there is a next row
		var student types.Students
		err:=rows.Scan(&student.Id,&student.Name,&student.Email,&student.Age)
		if err!=nil{
			return nil,err			
		}
		students=append(students,student)
	}

	slog.Info("fetched all students",slog.Int("count",len(students)))

	return students,nil
}