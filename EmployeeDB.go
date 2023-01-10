package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"strconv"

	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

type Employee struct {
	EID        string `json:"eid"`
	Name       string `json:"name"`
	Salary     int    `json:"salary"`
	ContactNum int    `json:"contactnum"`
}
type Employ struct {
	EID        string `field:"EID"`
	Name       string `field:"NAME"`
	Salary     int    `field:"SALARY"`
	ContactNum int    `field:"CONTACTNUM"`
}

var cfg = mysql.Config{
	User:   "root",
	Passwd: "Shettyrahul3@1",
	Net:    "tcp",
	Addr:   "127.0.0.1:3306",
	DBName: "MySQLTestWithGO",
}

func GetEmpById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	eid, er := params["eid"]

	if !er {
		return
	}

	var E Employ

	if err := db.QueryRow("SELECT * from EmployeeData where EID = "+eid).Scan(&E.EID, &E.Name, &E.Salary, &E.ContactNum); err != nil {
		//json.NewEncoder(w).Encode(E2)
		fmt.Printf("The given id %v doesn't exist\n "+err.Error()+"\n", eid)
		return
	}

	E2 := Employee(E)

	json.NewEncoder(w).Encode(E2)
}

func DelEmp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	eid, er := params["eid"]

	if !er {
		return
	}

	if _, err := db.Exec("DELETE from EmployeeData where eid = " + eid); err != nil {
		fmt.Printf("The given id: %v doesn't exist\n", eid)
		return
	}
}

func GetEmpSal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	name := r.URL.Query().Get("name")

	var E string

	if err := db.QueryRow("SELECT salary from EmployeeData where name = " + name).Scan(&E); err != nil {
		fmt.Printf("The given name: %v doesn't exist\n", name)
		return
	}

	json.NewEncoder(w).Encode(E)
}

func NewEmp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	Nam := r.URL.Query().Get("name")
	Sal := r.URL.Query().Get("salary")
	id := r.URL.Query().Get("id")
	CNo := r.URL.Query().Get("CNo")

	Slr, _ := strconv.Atoi(Sal)
	E2 := Employee{EID: id, Name: Nam, Salary: Slr, ContactNum: 345676}

	json.NewEncoder(w).Encode(E2)

	_, er := db.Exec("INSERT INTO EmployeeData VALUES(\"" + id + "\"," + Nam + "," + Sal + "," + "\"" + CNo + "\"" + ");")

	if er != nil {
		fmt.Println("Error Adding Value to DataBase" + er.Error())
	}

}

func main() {

	r := mux.NewRouter()

	d, err := sql.Open("mysql", cfg.FormatDSN())
	db = d
	if err != nil {
		log.Fatal(err)
		return
	}

	r.HandleFunc("/Employee/{eid}", GetEmpById).Methods("GET")
	r.HandleFunc("/Employee/", GetEmpSal).Methods("GET")
	r.HandleFunc("/Employee/", NewEmp).Methods("POST")
	r.HandleFunc("/Employee/{eid}", DelEmp).Methods("DELETE")

	fmt.Println("Starting Server at Port 8000")

	//fmt.Println(len(empl))
	log.Fatal(http.ListenAndServe(":8000", r))

	defer (*db).Close()

}
