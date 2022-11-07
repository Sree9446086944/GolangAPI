package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shopspring/decimal"
)

// create database go_rest_api_crud;
// go get github.com/gorilla/mux  -  creating routes and http handlers

// go get github.com/jinzhu/gorm - Go ORM for mysql

// go get github.com/go-sql-driver/mysql - mysql driver

// github.com/jinzhu/gorm - version1 and gorm.io/gorm - version2

var db *gorm.DB
var err error

// Product is a representation of a product
type Product struct {
	ID    int             `form:"id" json:"id"`
	Code  string          `form:"code" json:"code"`
	Name  string          `form:"name" json:"name"`
	Price decimal.Decimal `form:"price" json:"price" sql:"type:decimal(16,2);"`
}

// Result is an arrayb of products
type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {

	db, err = gorm.Open("mysql", "root:<password>@/go_rest_api_crud?charset=utf8&parseTime=True")

	if err != nil {
		log.Println("Connection failed", err)
	} else {
		log.Println("Connection established")
	}

	db.AutoMigrate(&Product{}) // migrate/create Product table to db automatically
	//creates table named products in db

	handleRequests()
}

func handleRequests() {
	log.Println("Start the development server at http://127.0.0.1:9999")

	myRouter := mux.NewRouter().StrictSlash(true) //created route, StrictSlash(true) will add / at end of route

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/api/products", createProduct).Methods("POST")
	myRouter.HandleFunc("/api/products", getProducts).Methods("GET")
	myRouter.HandleFunc("/api/products/{id}", getProduct).Methods("GET")
	myRouter.HandleFunc("/api/products/{id}", updateProduct).Methods("PUT")
	myRouter.HandleFunc("/api/products/{id}", deleteProduct).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9999", myRouter))

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome") // formats and write to w
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Create product")
	//get req body from request r
	payLoad, _ := ioutil.ReadAll(r.Body)

	var product Product
	//json to object
	json.Unmarshal(payLoad, &product)

	//create record in table
	db.Create(&product)

	//creating response body
	res := Result{Code: 200, Data: product, Message: "Success create product"}

	//convert response struct to json
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//setting response with header,and result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	products := []Product{}

	db.Find(&products)

	res := Result{Code: 200, Data: products, Message: "Success get products"}

	//obj to json
	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//setting response with header,and result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["id"]

	var product Product
	db.First(&product, productId)

	res := Result{Code: 200, Data: product, Message: "Success get products"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["id"]

	payLoad, _ := ioutil.ReadAll(r.Body)

	var productUpdates Product
	//json to object
	json.Unmarshal(payLoad, &productUpdates)

	var product Product
	db.First(&product, productId)
	db.Model(&product).Updates(productUpdates)

	//creating response body
	res := Result{Code: 200, Data: product, Message: "Success create product"}

	//convert response struct to json
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//setting response with header,and result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["id"]

	var product Product

	db.First(&product, productId)
	db.Delete(&product)

	res := Result{Code: 200, Message: "Success delete product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
