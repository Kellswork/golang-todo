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

// var client *mongo.Client
// var db *mongo.Database

const (
	dbName         string = "todo-example"
	collectionName string = "todo"
)

// step 3 create struct type
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
	CreateTodo struct {
		Title string `json:"title"`
	}
	UpdateTodo struct {
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}

	// we create a service struct to hold the dependencies required by our handlers in this case only the db connection for now
	// we create a service struct and add the db as as fiekd, this way we can make the db variable available to the handlers. becuase it is no longe a global variable
	Service struct {
		db *mongo.Database
	}
)

func init() {
	// step 4 create init function and connect to database
	fmt.Println("init function running")
	rnd = renderer.New(
		renderer.Options{
			ParseGlobPattern: "html/*.html",
		},
	)

}

// for the function return type, you dont have to add the vairaible names only the types except its a named return type
func initializeDB(uri string) (*mongo.Client, *mongo.Database) {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	// log error for database connect failure
	checkError(err)
	err = client.Ping(ctx, readpref.Primary())
	checkError(err)
	db := client.Database(dbName)
	return client, db
}

func homeHandler(rw http.ResponseWriter, r *http.Request) {
	// filePath := "./README.md"
	// err := rnd.FileView(rw, http.StatusOK, filePath, "readme.md")
	err := rnd.HTML(rw, http.StatusOK, "indexPage", nil)
	checkError(err)
}

func (s Service) getTodos(rw http.ResponseWriter, r *http.Request) {
	// initialise a variable and assign empty array with type todomodel
	var todoListFromDB = []TodoModel{}
	// fetch all the todos stored in the databse collection
	filter := bson.D{}
	cursor, err := s.db.Collection(collectionName).Find(context.Background(), filter)
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

func (s Service) createTodo(rw http.ResponseWriter, r *http.Request) {
	var todoReq CreateTodo
	// process the client input, if it returns an error , send a http response of bad request
	if err := json.NewDecoder(r.Body).Decode(&todoReq); err != nil {
		log.Printf("failed to decode json data: %v\n", err.Error())
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "could not decode data",
		})
		return
	}
	// check if the title is empty, return a response error message of title required
	if todoReq.Title == "" {
		log.Println("no title added to response body")
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "please add a title",
		})
		return
	}
	// create a todomodel for adding a todo to the database
	todoModel := TodoModel{
		ID:        primitive.NewObjectID(),
		Title:     todoReq.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	// add the todo to the database
	data, err := s.db.Collection(collectionName).InsertOne(r.Context(), todoModel)
	// return http status response if todo failed to save to the database
	if err != nil {
		log.Printf("failed to insert data into the database: %v\n", err.Error())
		rnd.JSON(rw, http.StatusInternalServerError, renderer.M{
			"message": "Failed to insert data into the database",
			"error":   err.Error(),
		})
		return
	}
	// if successfull, return a http status response success with the id.
	rnd.JSON(rw, http.StatusCreated, renderer.M{
		"message": "Todo created successfully",
		"ID":      data.InsertedID,
	})
}

func (s Service) updateTodo(rw http.ResponseWriter, r *http.Request) {
	// get the id from the url params
	log.Printf("url value: %v\n", r.URL)
	id := chi.URLParam(r, "id")
	log.Printf("URL Param 'id': %s", id)
	// check if the id is a hex value because we stored it as a hex value, if error return a message with id invalid
	res, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Printf("the id param: %v\n", err.Error())
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "The id is invalid",
		})
		return
	}

	var updateTodoReq UpdateTodo
	if err := json.NewDecoder(r.Body).Decode(&updateTodoReq); err != nil {
		log.Printf("failed to decode the json response body data: %v\n", err.Error())
		rnd.JSON(rw, http.StatusInternalServerError, err.Error())
		return
	}

	if updateTodoReq.Title == "" {
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "please add a title",
		})
		return
	}
	// update the todo in the database
	filter := bson.M{"id": res}
	update := bson.M{"$set": bson.M{"title": updateTodoReq.Title, "completed": updateTodoReq.Completed}}
	data, err := s.db.Collection(collectionName).UpdateOne(r.Context(), filter, update)

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

func (s Service) deleteTodo(rw http.ResponseWriter, r *http.Request) {
	fmt.Printf("got here")
	id := chi.URLParam(r, "id")
	res, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("invalid id: %v\n", err.Error())
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "could not decode data",
		})
		return
	}

	filter := bson.M{"id": res}
	// delete that todo entry in the database
	if data, err := s.db.Collection(collectionName).DeleteOne(r.Context(), filter); err != nil {
		log.Printf("could not delete item from database: %v\n", err.Error())
		rnd.JSON(rw, http.StatusInternalServerError, renderer.M{
			"message": "an error eccoured while deleting todo item",
		})
	} else {
		rnd.JSON(rw, http.StatusOK, renderer.M{
			"message": "Todo deleted successfully",
			"data":    data.DeletedCount,
		})
	}

}

func main() {
	client, db := initializeDB("mongodb://localhost:27017")

	service := Service{
		db: db,
	}
	// step 8 create router and route handlers for home
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// publish the css file so the html file can use the styles
	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	router.Get("/", homeHandler)
	// update the todoHandlers to reciever methods in Go
	router.Mount("/todo", todoHandlers(service)) // Mount attaches another http.Handler along ./pattern/*

	server := &http.Server{
		Addr:         ":9000",
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	// start the server in a seperate go routine
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
// update
func todoHandlers(service Service) http.Handler {
	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		r.Get("/", service.getTodos)
		r.Post("/", service.createTodo)
		r.Put("/{id}", service.updateTodo)
		r.Delete("/{id}", service.deleteTodo)
	})
	return router
}

// step 5 define checkErroro function
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
