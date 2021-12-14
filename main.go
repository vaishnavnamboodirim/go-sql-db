package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Student struct {
	Id    string  `json:"id"`
	Name  string  `json:"student1"`
	Class string  `json:"class"`
	Marks float64 `json:"marks"`
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "root:Vaishnav@1234@tcp(127.0.0.1:3306)/<sql-go-db>")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/students", getStudents).Methods("GET")
	r.HandleFunc("/students", createStudent).Methods("POST")
	r.HandleFunc("/students/{id}", getStudent).Methods("GET")
	r.HandleFunc("/students/{id}", deleteStudent).Methods("DELETE")

	http.ListenAndServe(":3306", r)

}

func getStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var students []Student

	result, err := db.Query("SELECT id, name, class, marks from posts")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {
		var student Student
		err := result.Scan(&student.Id, &student.Name, &student.Class, &student.Marks)
		if err != nil {
			panic(err.Error())
		}
		students = append(students, student)
	}

	json.NewEncoder(w).Encode(students)
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stmt, err := db.Prepare("INSERT INTO students(Id, name, class, marks) VALUES(?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	var student Student

	json.Unmarshal(body, &student)

	student2 := Student{"104", "Student2", "5A", 450}
	newId := student2.Id
	newName := student2.Name
	newClass := student2.Class
	newMarks := student2.Marks

	_, err = stmt.Exec(newId, newName, newClass, newMarks)
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "New student was created")

}

func getStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := db.Query("SELECT id,name,class,marks FROM students WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var student Student

	for result.Next() {
		err := result.Scan(&student.Id, &student.Name, &student.Class, &student.Marks)
		if err != nil {
			panic(err.Error())
		}
	}

	json.NewEncoder(w).Encode(student)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	stmt, err := db.Prepare("DELETE FROM posts WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "post with ID = %s was deleted", params["id"])

}
