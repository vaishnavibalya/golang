package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Make a worker pool that processes the list of jobs concurrently.
// The worker should:
//   - print out "Worker <n> retreiving <type>/<id>",
//   - request the item from Star Wars API (https://swapi.dev/api/<type/<id>)
// 	 - after the req finishes, print out "Worker <n> found <type>/<id> <response.name>".

// Example Working output:
// Worker  3 retreiving people/1
// Worker  2 retreiving people/3
// Worker  1 retreiving people/2
// Worker 2 found people/3 R2-D2
// Worker  2 retreiving people/4
// Worker 1 found people/2 C-3PO
// Worker  1 retreiving people/5
// Worker 3 found people/1 Luke Skywalker
// Worker  3 retreiving planets/1
// Worker 1 found people/5 Leia Organa
// Worker  1 retreiving planets/2
// Worker 3 found planets/1 Tatooine
// Worker  3 retreiving planets/3
// Worker 2 found people/4 Darth Vader
// Worker  2 retreiving planets/4
// Worker 1 found planets/2 Alderaan
// Worker  1 retreiving planets/5
// Worker 3 found planets/3 Yavin IV
// Worker  3 retreiving starships/2
// Worker 2 found planets/4 Hoth
// Worker  2 retreiving starships/3
// Worker 1 found planets/5 Dagobah
// Worker  1 retreiving starships/5
// Worker 3 found starships/2 CR90 corvette
// Worker  3 retreiving starships/9
// Worker 2 found starships/3 Star Destroyer
// Worker  2 retreiving starships/10
// Worker 1 found starships/5 Sentinel-class landing craft
// Worker 3 found starships/9 Death Star
// Worker 2 found starships/10 Millennium Falcon

type task struct {
	ID   string
	Type string
}

func Makerequest(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	Body := make(map[string]interface{})
	err = json.Unmarshal(body, &Body)
	if err != nil {
		log.Fatal(err)
	}
	name := Body["name"].(string)
	return name
}

func Work(wnumber int, jobs <-chan task, result chan<- int) {
	for job := range jobs {
		fmt.Println("Worker", wnumber, "retreiving", job.Type+"/"+job.ID)
		url := "https://swapi.dev/api/" + job.Type + "/" + job.ID
		name := Makerequest(url)
		fmt.Println("Worker", wnumber, "found", job.Type+"/"+job.ID, name)
		result <- 1
	}
}

func main() {
	todo := []task{
		{"1", "people"},
		{"2", "people"},
		{"3", "people"},
		{"4", "people"},
		{"5", "people"},
		{"1", "planets"},
		{"2", "planets"},
		{"3", "planets"},
		{"4", "planets"},
		{"5", "planets"},
		{"2", "starships"},
		{"3", "starships"},
		{"5", "starships"},
		{"9", "starships"},
		{"10", "starships"},
	}

	// Process todos with workers.
	nworkers := 3
	ntasks := len(todo)
	result := make(chan int, ntasks)
	jobs := make(chan task, ntasks)

	for i := 1; i <= nworkers; i++ {
		go Work(i, jobs, result)
	}
	for _, job := range todo {
		jobs <- job
	}
	close(jobs)

	for i := 0; i < ntasks; i++ {
		<-result
	}
}
