// import package main
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client
var dbUrl string

// this creates an in-memeory database that you can use for tetsing
// docker_test

// testing parts
// step 1. create the test funtion
// step 2. setup an im-memory dockertest database
//step 3. create an initiolizeDB function or method that takes a uri paramater. you can also initilaize the db to take a client paramater instead
// step 3.5. still on the function, it should return the client and db, again see if it can take in a clinet instead becasue the mongodb one returnd a client. this is just for testing purposes
// update the dockertest so you can get the url in applyuri
// you also did some installation, for the docker test, testify
// step 4. update the todoHandlers to reciever methods
// step 5 update the todohandler function route to be a reciever method
// step I missed, before updatinmg the reciewver method
// create a service struct type and add the db as as fiekd, to make the db variable available to handlers
// update the handlers to use the service struct type as a recievcer method
// read up on reciever method

func TestMain(m *testing.M) {
	fmt.Println("test main running")
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pull mongodb docker image for version 5.0
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "5.0",
		Env: []string{
			// username and password for mongodb superuser
			"MONGO_INITDB_ROOT_USERNAME=root",
			"MONGO_INITDB_ROOT_PASSWORD=password",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	err = pool.Retry(func() error {
		var err error
		// save global test variable to be used above
		dbUrl = fmt.Sprintf("mongodb://root:password@localhost:%s", resource.GetPort("27017/tcp"))
		dbClient, err = mongo.Connect(
			context.TODO(),
			options.Client().ApplyURI(
				dbUrl,
			),
		)
		if err != nil {
			return err
		}
		return dbClient.Ping(context.TODO(), nil)
	})

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// run tests
	code := m.Run()

	// When you're done, kill and remove the container
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	// disconnect mongodb client
	if err = dbClient.Disconnect(context.TODO()); err != nil {
		panic(err)
	}

	os.Exit(code)

}

// import testify, to use for test assertion
// if you are familair with chai in react testing, this will seem familiar

// create Todo test function
// in most real life testing situation, you will need a tetsing libabrary for assertion and mocks. the msot popular one used in golanbg is testify
func TestCreateTodo(t *testing.T) {
	// thinsg to test
	// 1. that the response body was decoded successfully
	// 2. that the todo title whne empty returns a response
	// 3. a todo was succesfully added to the DB
	// generally, tests are for testing the response bodies
	// any method that accepts responsewriter also accepts response recorder
	// in the handlerFunc, I was confused if we are passing in the type or the function

	// user input is a json string
	jsonStr := []byte(`{"title": "go to the beach"}`)

	req, err := http.NewRequest("POST", "/todo", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	// we use ther initially created in memmory test dbUrl created in TestMain to initialize a db connection
	_, db := initializeDB(dbUrl)
	service := Service{
		db: db,
	}
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(service.createTodo)
	handler.ServeHTTP(recorder, req)

	// use tesitify assert to check the json respobnse returned is what we expected
	// why did we use require for one and assert for one
	// require because it will stop the tes# from running and if the sttus code is not 201 the no need to continue
	require.Equal(t, http.StatusCreated, recorder.Code)

	// this is called an anunimous struct. so we used it to create a struct for the json reponse body retuned.
	// we can then use the mesasge struct to assert that the message returned is what we expected
	result := struct {
		Message string `json:"message"`
		ID      string `json:"id"`
	}{}
	err = json.Unmarshal(recorder.Body.Bytes(), &result)
	assert.Nil(t, err)
	assert.Equal(t, "Todo created successfully", result.Message)

	// if status := recorder.Code; status != http.StatusCreated {
	// 	t.Errorf("createTodo returned wrong status code. Got: %v Expected %v", status, http.StatusCreated)
	// }

	// expected should be what you set the response body to return in the functiuon handler
	// expected := `{"message": "Todo created successfully"}`

	// if recorder.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v expected: %v", recorder.Body.String(), expected)
	// }

	// to test if this is working, make sure to have your DB is up and running.
	// I stopped at the part were the test was failing due to a different error response its getting
	// when you get back to this, run the code to see the failing test result
	// check chat gpt for the sample testing code it created

}
