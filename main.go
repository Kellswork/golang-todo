package main

// step 1 import packages
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// step 2 create variables
var rnd *renderer.Render
var client *mongo.Client
var db *mongo.Database

const (
	dbName         string = "todo-example"
	collectionName string = "todo"
)

// step 3 create a todo model struct type fior the mongo db database and a todo struct type for the frontend
type (
	TodoModel struct {
		ID        primitive.ObjectID `bson:"id,omitempty"`
		Title     string             `bson:"title"`
		Completed bool               `bson:"completed"`
		CreatedAt time.Time          `bson:"created_at"`
	}

	Todo struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"created_at"`
	}

	GetTodoResponse struct {
		Message string `json:"message"`
		Data    []Todo `json:"data"`
	}
)

func init() {
	// step 4 create init function and connect to database
	fmt.Println("init function running")
	rnd = renderer.New()
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	// log error for database connect failure
	checkError(err)
	err = client.Ping(ctx, readpref.Primary())
	checkError(err)
	db = client.Database(dbName)

}

func homeHandler(rw http.ResponseWriter, r *http.Request) {
	filePath := "./README.md"
	err := rnd.FileView(rw, http.StatusOK, filePath, "readme.md")
	checkError(err)
}

func getTodos(rw http.ResponseWriter, r *http.Request) {
	// initialise a variable and assign empty array with type todomodel
	var todoListFromDB = []TodoModel{}
	// fetch all the todos stored in the databse collection
	filter := bson.D{}
	cursor, err := db.Collection(collectionName).Find(context.Background(), filter)
	if err != nil {
		log.Printf("failed to fetch todo records from the db: %v\n", err.Error())
		// render a json error message and the error
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "could not fetch the todo collection",
			"error":   err.Error(),
		})
		return
	}
	// create a todo list variable and assign it an empty array with type
	todoList := []Todo{}
	if err = cursor.All(context.Background(), &todoListFromDB); err != nil {
		checkError(err)
	}

	// loop through the database array, convert to json using the todomodel and append to the todolist array.
	for _, td := range todoListFromDB {
		todoList = append(todoList, Todo{
			ID:        td.ID.Hex(),
			Title:     td.Title,
			Completed: td.Completed,
			CreatedAt: td.CreatedAt,
		})
	}

	// render a JSON response for successfully fetching the data with a custom getTodoResponse type
	rnd.JSON(rw, http.StatusOK, GetTodoResponse{
		Message: "All todos retrieved",
		Data:    todoList,
	})
}

func createTodo(rw http.ResponseWriter, r *http.Request) {
	// create variable todo for storing user input
	var todo Todo
	// process the client input, if it returns an error , send a http response of bad request
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		log.Printf("failed to decode json data: %v\n", err.Error())
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "could not decode data",
		})
		return
	}
	// check if the title is empty, return a response error message of title required
	if todo.Title == "" {
		log.Println("no title added to response body")
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "please add a title",
		})
		return
	}
	// create a todomodel for adding a todo to the database
	todoModel := TodoModel{
		ID:        primitive.NewObjectID(),
		Title:     todo.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	// add the todo to the database
	data, err := db.Collection(collectionName).InsertOne(r.Context(), todoModel)
	// return http status response if todo failed to save to the database
	if err != nil {
		log.Printf("failed to insert data into the database: %v\n", err.Error())
		rnd.JSON(rw, http.StatusInternalServerError, renderer.M{
			"message": "Failed to insert data into the database",
			"error":   err.Error(),
		})
		return
	}
	// if successfull, return a httpo status response success with the id.
	rnd.JSON(rw, http.StatusCreated, renderer.M{
		"message": "Todo created successfully",
		"ID":      data.InsertedID,
	})
}

func updateTodo(rw http.ResponseWriter, r *http.Request) {
	// get the id from the url params
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	// check if the id is a hex value because we stored it as a hex value, if error return a message with id invalid
	res, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("the id param is not a valid hex value: %v\n", err.Error())
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "The id is invalid",
			"error":   err.Error(),
		})
		return
	}

	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		log.Printf("failed to decode the json response body data: %v\n", err.Error())
		rnd.JSON(rw, http.StatusInternalServerError, err.Error())
	}
	if todo.Title == "" {
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "Title cannot be empty",
		})
		return
	}
	// update the todo in the database
	filter := bson.M{"id": res}
	update := bson.M{"$set": bson.M{"title": todo.Title, "completed": todo.Completed}}
	data, err := db.Collection(collectionName).UpdateOne(r.Context(), filter, update)

	if err != nil {
		log.Printf("failed to update db collection: %v\n", err.Error())
		rnd.JSON(rw, http.StatusInternalServerError, renderer.M{
			"message": "Failed to update data in the database",
			"error":   err.Error(),
		})
		return
	}
	rnd.JSON(rw, http.StatusOK, renderer.M{
		"message": "Todo updated successfully",
		"data":    data.ModifiedCount,
	})
}

func deleteTodo(rw http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	res, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("invalid id: %v\n", err.Error())
		rnd.JSON(rw, http.StatusBadRequest, err.Error())
		return
	}

	filter := bson.M{"id": res}
	// delete that todo entry in the database
	if data, err := db.Collection(collectionName).DeleteOne(r.Context(), filter); err != nil {
		log.Printf("could not delete item from database: %v\n", err.Error())
		rnd.JSON(rw, http.StatusInternalServerError, renderer.M{
			"message": "an error eccoured while deleting todo item",
			"error":   err.Error(),
		})
	} else {
		rnd.JSON(rw, http.StatusOK, renderer.M{
			"message": "item deleted successfully",
			"data":    data,
		})
	}

}

func main() {
	// step 8 create router and route handlers for home
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", homeHandler)          // for homepage and todo routes
	router.Mount("/todo", todoHandlers()) // Mount attaches another http.Handler along ./pattern/*

	// step 6 connect to a server
	server := &http.Server{
		Addr:         ":9000",
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	// step 7: stop the server
	// create channel to reccieve signal
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	// start the server in a seperate go routine. why? when done, try this without the go func , also ask vic if this is neccessary, that's the go func
	go func() {
		fmt.Println("Server started on port", 9000)
		if err := server.ListenAndServe(); err != nil {
			log.Printf("listen:%s\n", err)
		}
	}()

	// wait for a signal to shutdown the server
	sig := <-stopChan
	log.Printf("shutting down server: %v\n", sig)
	if err := client.Disconnect(context.Background()); err != nil {
		panic(err)
	}

	// create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v\n", err)
	}
	log.Println("Server shutdown gracefully")

}

// step 9: create a group route for todo routers
func todoHandlers() http.Handler {
	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		r.Get("/", getTodos)
		r.Post("/", createTodo)
		r.Put("/{id}", updateTodo)
		r.Delete("/{id}", deleteTodo)
	})
	return router
}

// step 5 define checkErroro function
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
