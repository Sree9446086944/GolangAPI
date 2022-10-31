package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

//create restapi - gin pkg
//setup dependency tracking in app - go mod init
//go get github.com/gin-gonic/gin

// define our todo structure using struct
type todo struct {
	ID        string `json:"id"` // like a template
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

// array of todos , this is different datastructre not json
var todos = []todo{
	{ID: "1", Item: "Clean Room", Completed: false},
	{ID: "2", Item: "Read Book", Completed: false},
	{ID: "3", Item: "Record Video", Completed: false},
}

// client and server interact using json
//convert our datastructure to json

// main run at start by default
func main() {
	//create server
	router := gin.Default() // this router is our server

	//create endpoint
	router.GET("/todos", getTodos) // GET(path to append to url, fn returning json data)

	//add todo
	router.POST("/todos", addTodo)

	//get todo by id
	router.GET("/todos/:id", getTodo)

	//PATCH since updating already existing todo parameter
	router.PATCH("/todos/:id", toggleTodoStatus)

	//to run the server
	router.Run("localhost:9090") //default 8080, if not given
}

func getTodos(context *gin.Context) { // this context have info about the incoming http req like req body, params, header etc

	//converting []todos datastructure to json
	context.IndentedJSON(http.StatusOK, todos) // (status code, obj to convert to json)
}

func addTodo(context *gin.Context) {
	var newTodo todo
	//take json from req body and bind to newTodo using BindJSON(), &newTodo since there is change in newTodo
	if err := context.BindJSON(&newTodo); err != nil { //throw error if error
		// if error no need to continueś
		return
	}
	//if no error
	todos = append(todos, newTodo)
	//return the newTodo in json format
	context.IndentedJSON(http.StatusCreated, newTodo)
}

func getTodo(context *gin.Context) {
	//all req info in context, get id from path param
	id := context.Param("id")
	todo, err := getTodoById(id)
	//if err , then give custom message and return, dont continue
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{ //gin.H{} - for custom message as map
			"message": "Todo not found",
		})
		return
	}
	//if err is nil
	context.IndentedJSON(http.StatusOK, todo)
}

func getTodoById(id string) (*todo, error) { //if todo type return then error is nil, if error todo is nil, here *todo is pointer type
	for i, t := range todos {
		if t.ID == id {
			return &todos[i], nil
		}
	}
	return nil, errors.New("todo not found") //error if not found
}

func toggleTodoStatus(context *gin.Context) {
	id := context.Param("id")
	todo, err := getTodoById(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{
			"message": "Todo not found",
		})
		return
	}
	//if not error, toggle/update the completed status
	todo.Completed = !todo.Completed
	//return the updated todo
	context.IndentedJSON(http.StatusOK, todo)
}

// http://localhost:9090/todos/1 - PATCH - toggles completed status each time hit req

//http://localhost:9090/todos/1 - GET - get by id

//http://localhost:9090/todos - GET - get all todos

//http://localhost:9090/todos - POST
/*{
	"id": "1",
	"item": "Make Bed",
	"completed": false
}*/
