package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/alexgtn/go-middleware-metrics/github.com/middleware"
	"github.com/gorilla/mux"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Employee request model
type Employee struct {
	// Id of the employee
	ID string `json:"id"`
	//Isbn of the employee
	Isbn string `json:"isbn"`
	// First Name of the employee
	Firstname string `json:"fname"`
	// Last Name of the employee
	Lastname string `json:"lname"`
}

var employees []Employee

// Employee response payload
// swagger:response employeeRes
type swaggEmployeeRes struct {
	// in:body
	Body Employee
}

// Success response
// swagger:response okResp
type swaggRespOk struct {
	// in:body
	Body struct {
		// HTTP status code 200 - OK
		Code int `json:"code"`
	}
}

// Error Bad Request
// swagger:response badReq
type swaggReqBadRequest struct {
	// in:body
	Body struct {
		// HTTP status code 400 -  Bad Request
		Code int `json:"code"`
	}
}

// Error Not Found
// swagger:response notFoundReq
type swaggReqNotFound struct {
	// in:body
	Body struct {
		// HTTP status code 404 -  Not Found
		Code int `json:"code"`
	}
}

func getSwag(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello")
	w.Header().Set("Content-Type", "application/json")
	file, err := os.Open("swaggerui/swagger.json") // For read access.
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(w, file)
}

func getEmployees(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	json.NewEncoder(w).Encode(employees)
}

func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)

	for index, item := range employees {

		if item.ID == params["id"] {
			employees = append(employees[:index], employees[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(employees)
}

func getEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	for _, item := range employees {
		if item.ID == params["id"] {

			json.NewEncoder(w).Encode(item)
			return

		}
	}
}

func createEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var employee Employee
	_ = json.NewDecoder(r.Body).Decode(&employee)
	employee.ID = strconv.Itoa(rand.Intn(100000000))
	employees = append(employees, employee)
	json.NewEncoder(w).Encode(employee)
}

func updateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var employee Employee
	_ = json.NewDecoder(r.Body).Decode(&employee)

	for index, item := range employees {
		if item.ID == employee.ID {
			employees[index] = employee
			json.NewEncoder(w).Encode(employee)
			return
		}
	}
}
func main() {
	employees = append(employees, Employee{ID: "1", Isbn: "12345", Firstname: "Anand", Lastname: "Pandey"})
	employees = append(employees, Employee{ID: "2", Isbn: "13245", Firstname: "Siddharth", Lastname: "Soni"})

	r := mux.NewRouter()
	metricsMiddleware := middleware.NewMetricsMiddleware()

	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/lemon", lemonHandler).Methods(http.MethodGet)
	r.HandleFunc("/potato", potatoHandler).Methods(http.MethodPost)
	r.HandleFunc("/employees", getEmployees).Methods("GET")
	r.HandleFunc("/employees/{id}", getEmployee).Methods("GET")
	r.HandleFunc("/employees", createEmployee).Methods("POST")
	r.HandleFunc("/employees", updateEmployee).Methods("PUT")
	r.HandleFunc("/employees/{id}", deleteEmployee).Methods("DELETE")
	r.HandleFunc("/v2/api-docs", getSwag).Methods("GET")

	fs := http.FileServer(http.Dir("./swaggerui"))
	r.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", fs))

	r.Use(metricsMiddleware.Metrics)
	fmt.Println("Server has started on 8080: ")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func lemonHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Lemon"))
}

func potatoHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Potato"))
}
