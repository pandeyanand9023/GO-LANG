// Package classification Employee API.
//
// These are the services in the form of RESTful APIs in golang.
// Here Swagger is used  to provide a detailed overview of the language specs
//
// Terms Of Service:
//
//     Schemes: http, https
//     Host: localhost:8080
//     Version: 1.0.0
//     Contact: Anand Pandey<pandey_anand@optum.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - api_key:
//
//     SecurityDefinitions:
//     api_key:
//          type: apiKey
//          name: KEY
//          in: header
//
// swagger:met
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

	// swagger:operation GET /employees employees getEmployees
	// ---
	// summary: Returns list of Employees
	// description: Return the list of Employees.
	// parameters:
	// responses:
	//   "200":
	//     "$ref": "#/responses/employeeRes"
	//   "400":
	//     "$ref": "#/responses/badReq"
	//   "404":
	//     "$ref": "#/responses/notFoundReq"
	r.HandleFunc("/employees", getEmployees).Methods("GET")

	// swagger:operation GET /employees/{id} employees getEmployee
	// ---
	// summary: Return an Employee provided by the id.
	// description: If the employee is found, employee will be returned else Error Not Found (404) will be returned.
	// parameters:
	// - name: id
	//   in: path
	//   description: id of the employee
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/employeeRes"
	//   "400":
	//     "$ref": "#/responses/badReq"
	//   "404":
	//     "$ref": "#/responses/notFoundReq"
	r.HandleFunc("/employees/{id}", getEmployee).Methods("GET")

	// swagger:operation POST /employees employees createEmployee
	// ---
	// summary: Creates a new employee.
	// description: If employee creation is success, employee will be  returned with Created (201).
	// parameters:
	// - name: employee
	//   description: employee to add to the list of accounts
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/Employee"
	// responses:
	//   "200":
	//     "$ref": "#/responses/okResp"
	//   "400":
	//     "$ref": "#/responses/badReq"
	r.HandleFunc("/employees", createEmployee).Methods("POST")

	// swagger:operation PUT /employees/{id} employees updateEmployee
	// ---
	// summary: Updates a current employee.
	// description: If employee updation is success, employee will be returned with  (201).
	// parameters:
	// - name: id
	//   in: path
	//   description: employee id
	//   type: string
	//   required: true
	// - name: employee
	//   description: Information of employee to add to the list of accounts
	//   in: body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/Employee"
	// responses:
	//   "200":
	//     "$ref": "#/responses/okResp"
	//   "400":
	//     "$ref": "#/responses/badReq"
	//   "404":
	//     "$ref": "#/responses/notFoundReq"
	r.HandleFunc("/employees", updateEmployee).Methods("PUT")

	// swagger:operation DELETE /employees/{id} employees deleteEmployee
	// ---
	// summary: Deletes requested employee by employee id.
	// description: Depending on the employee id, HTTP Status Not Found (404) or HTTP Status OK (200) may be returned.
	// parameters:
	// - name: id
	//   in: path
	//   description: employee id
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/okResp"
	//   "400":
	//     "$ref": "#/responses/badReq"
	//   "404":
	//     "$ref": "#/responses/notFoundReq"
	r.HandleFunc("/employees/{id}", deleteEmployee).Methods("DELETE")

	fs := http.FileServer(http.Dir("./swaggerui"))
	r.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", fs))

	fmt.Println("Server has started on 8080: ")
	log.Fatal(http.ListenAndServe(":8080", r))
}
