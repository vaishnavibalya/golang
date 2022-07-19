package main

import (
	"testing"
)

type MockHttpCall struct {
}

var dummyhttp = MockHttpCall{}

var w = Worker{Http: &dummyhttp}

//ImageLoad function
func (h *MockHttpCall) Makerequest(url string) string {
	return "name"
}

func TestWork(t *testing.T) {
	ntasks := 15
	jobs := make(chan task, ntasks)
	var result = make(chan int, ntasks)
	jobs <- task{"1", "people"}
	jobs <- task{"2", "people"}
	close(jobs)
	w.Work(1, jobs, result)
	if len(jobs) == 0 {
		t.Logf("Success")
	} else {
		t.Errorf("failure")
	}
}
