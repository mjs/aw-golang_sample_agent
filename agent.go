package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const listenPort = 5000

type AgentHandler func(json.RawMessage) (interface{}, error)

type ReportMessage struct {
	Total int
}

type Memory struct {
	Checks   int
	Receives int
}

// register is called by ActiveWorkflow to determine metadata and options for
// the agent.
func register(json.RawMessage) (interface{}, error) {
	log.Print("Register")
	return map[string]interface{}{
		"name":            "GoCounterAgent",
		"display_name":    "Go Test Counter Agent",
		"description":     "Description goes here",
		"default_options": map[string]interface{}{},
	}, nil
}

// check is called by ActiveWorkflow on a regular basis if the user has
// configured a schedule for the agent. It increments the counter for the
// number of checks seen in the agent's memory and returns a message containing
// the new total count.
func check(rawParams json.RawMessage) (interface{}, error) {
	log.Print("check", string(rawParams))
	var params struct {
		Memory Memory
	}
	err := json.Unmarshal(rawParams, &params)
	if err != nil {
		return nil, err
	}

	memory := params.Memory
	memory.Checks += 1

	return map[string]interface{}{
		"logs":     []string{"Check done"},
		"errors":   []string{},
		"memory":   memory,
		"messages": []ReportMessage{{Total: memory.Checks + memory.Receives}},
	}, nil
}

// receive is called whenever the agent receives a message from another agent.
// The counter for the number of messages seen is updated in the agent state
// and a new message containing the new total count is returned.
func receive(rawParams json.RawMessage) (interface{}, error) {
	fmt.Println("Receive", string(rawParams))
	var params struct {
		Memory  Memory
		Message struct {
			Id int
			// This agent doesn't care about anything but the message id so no
			// other fields are included for unmarshalling.
		}
	}
	err := json.Unmarshal(rawParams, &params)
	if err != nil {
		return nil, err
	}

	memory := params.Memory
	memory.Receives += 1

	return map[string]interface{}{
		"logs":     []string{fmt.Sprintf("Received message %d", params.Message.Id)},
		"errors":   []string{},
		"memory":   memory,
		"messages": []ReportMessage{{Total: memory.Checks + memory.Receives}},
	}, nil
}

func handle(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Unmarshal the top layer of JSON to determine which remote API method is
	// being called.
	var input struct {
		Method string
		Params json.RawMessage
	}
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "could not parse input: %v", err)
		return
	}

	// Find a handler for the type of request.
	handler := lookupHandler(input.Method)
	if handler == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "unknown method: %q", input.Method)
		return
	}

	// Call the handler.
	out, err := handler(input.Params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	// Put the result in a "result" field and convert it to JSON.
	err = json.NewEncoder(w).Encode(map[string]interface{}{"result": out})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
}

func lookupHandler(method string) AgentHandler {
	switch method {
	case "register":
		return register
	case "check":
		return check
	case "receive":
		return receive
	default:
		return nil
	}
}

func main() {
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", listenPort), nil))
}
