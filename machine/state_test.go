package machine

import (
	"context"
	"log"
	"testing"
)

func TestMachine(t *testing.T) {
	var writeInitialHandshake, readHandshakeResponse, writeRespone State
	counter := 0
	writeInitialHandshake = func(runtime *Runtime) {
		counter += 5
		runtime.NextState(runtime.Context(), readHandshakeResponse)
	}
	readHandshakeResponse = func(runtime *Runtime) {
		counter += 6
		runtime.NextState(runtime.Context(), writeRespone)
	}
	writeRespone = func(runtime *Runtime) {
		counter += 6
		runtime.NextState(runtime.Context(), nil)
	}
	runtime := Run(context.Background(), writeInitialHandshake)
	<-runtime.Context().Done()
	log.Println(counter)

}
