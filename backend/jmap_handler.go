package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type JMAPRequest struct {
	MethodCalls [][]interface{} `json:"methodCalls"`
}

type JMAPResponse struct {
	MethodResponses [][]interface{} `json:"methodResponses"`
	SessionState    string          `json:"sessionState"`
}

func JMAPHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Allow requests from your React app
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")        // Allow POST and OPTIONS
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")         // Allow Content-Type header

		// Handle preflight requests (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req JMAPRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		methodResponses := [][]interface{}{}
		for _, call := range req.MethodCalls {
			methodName := call[0].(string)
			methodArgs := call[1].(map[string]interface{})
			callID := call[2].(string)

			switch methodName {
			case "Core/getCapabilities":
				respArgs := handleGetCapabilities()
				methodResponses = append(methodResponses, []interface{}{methodName, respArgs, callID})
			case "Core/getSession":
				respArgs := handleGetSession()
				methodResponses = append(methodResponses, []interface{}{methodName, respArgs, callID})
			case "Todo/query":
				respArgs, err := handleTodoQuery(db, methodArgs)
				if err != nil {
					http.Error(w, fmt.Sprintf("Error handling Todo/query: %v", err), http.StatusInternalServerError)
					return
				}
				methodResponses = append(methodResponses, []interface{}{methodName, respArgs, callID})
			case "Todo/get":
				respArgs, err := handleTodoGet(db, methodArgs)
				if err != nil {
					http.Error(w, fmt.Sprintf("Error handling Todo/get: %v", err), http.StatusInternalServerError)
					return
				}
				methodResponses = append(methodResponses, []interface{}{methodName, respArgs, callID})
			case "Todo/set":
				respArgs, err := handleTodoSet(db, methodArgs)
				if err != nil {
					http.Error(w, fmt.Sprintf("Error handling Todo/set: %v", err), http.StatusInternalServerError)
					return
				}
				methodResponses = append(methodResponses, []interface{}{methodName, respArgs, callID})
			default:
				log.Printf("Unknown JMAP method: %s", methodName)
				methodResponses = append(methodResponses, []interface{}{"error", map[string]interface{}{
					"type": "unknownMethod",
				}, callID})
			}
		}

		resp := JMAPResponse{
			MethodResponses: methodResponses,
			SessionState:    "server-session-state-1", // Simplified session state
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(resp)
	}
}

func handleGetCapabilities() map[string]interface{} {
	return map[string]interface{}{
		"capabilities": map[string]interface{}{
			"urn:ietf:params:jmap:core": map[string]interface{}{
				"maxSizeRequest":        10000000,
				"maxConcurrentUpload":   4,
				"maxConcurrentDownload": 4,
				"maxObjectsInGet":       500,
				"maxObjectsInSet":       500,
				"getStringEncoding":     "UTF-8",
				"maxCallsInRequest":     10,
				"maxObjectsInQuery":     100,
				"versions":              []string{"1"},
			},
			"urn:example:params:jmap:todo": map[string]interface{}{ // Example Todo capability
				"maxObjectsInQuery": 100,
				"maxObjectsInSet":   100,
				"querySortOptions":  []string{"id", "title"},
			},
		},
	}
}

func handleGetSession() map[string]interface{} {
	// In real JMAP, this is more complex and involves authentication.
	// Here, we return a simplified session.
	return map[string]interface{}{
		"capabilities": map[string]interface{}{
			"urn:ietf:params:jmap:core": map[string]interface{}{
				"maxSizeRequest":        10000000,
				"maxConcurrentUpload":   4,
				"maxConcurrentDownload": 4,
				"maxObjectsInGet":       500,
				"maxObjectsInSet":       500,
				"getStringEncoding":     "UTF-8",
				"maxCallsInRequest":     10,
				"maxObjectsInQuery":     100,
				"versions":              []string{"1"},
			},
			"urn:example:params:jmap:todo": map[string]interface{}{ // Example Todo capability
				"maxObjectsInQuery": 100,
				"maxObjectsInSet":   100,
				"querySortOptions":  []string{"id", "title"},
			},
		},
		"accounts": map[string]interface{}{
			"primary": map[string]interface{}{ // Example account
				"name":       "Primary Account",
				"isPersonal": true,
				"accountCapabilities": map[string]interface{}{
					"urn:example:params:jmap:todo": map[string]interface{}{
						"maxObjectsInQuery": 100,
						"maxObjectsInSet":   100,
					},
				},
			},
		},
		"username":       "user@example.com",                  // Simplified username
		"apiUrl":         "http://localhost:8080/jmap",        // API endpoint
		"downloadUrl":    "http://localhost:8080/download",    // Example, not used
		"uploadUrl":      "http://localhost:8080/upload",      // Example, not used
		"eventSourceUrl": "http://localhost:8080/eventsource", // Example, not used
	}
}

func handleTodoQuery(db *sql.DB, args map[string]interface{}) (map[string]interface{}, error) {
	todos, err := GetTodos(db)
	if err != nil {
		return nil, err
	}

	respData, err := MarshalTodosToJMAPResponse(todos)
	if err != nil {
		return nil, err
	}

	var respMap map[string]interface{}
	if err := json.Unmarshal(respData, &respMap); err != nil {
		return nil, err
	}

	// Extract the "Todo/query" method response part
	for _, methodResp := range respMap["methodResponses"].([]interface{}) {
		if methodResp.([]interface{})[0] == "Todo/query" {
			return methodResp.([]interface{})[1].(map[string]interface{}), nil
		}
	}

	return nil, fmt.Errorf("todo/query response not found in marshaled data")
}

func handleTodoGet(db *sql.DB, args map[string]interface{}) (map[string]interface{}, error) {
	ids, ok := args["ids"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'ids' argument in Todo/get")
	}

	todoList := []Todo{}
	notFoundIDs := []string{}

	for _, idInterface := range ids {
		idStr, ok := idInterface.(string)
		if !ok {
			return nil, fmt.Errorf("invalid id type in Todo/get, expected string")
		}
		todo, err := GetTodo(db, idStr)
		if err != nil {
			notFoundIDs = append(notFoundIDs, idStr) // Assume not found for simplicity in this example
			continue
		}
		todoList = append(todoList, *todo)
	}

	respData, err := MarshalTodosToJMAPResponse(todoList) // Re-use marshaler for simplicity, adjust if needed
	if err != nil {
		return nil, err
	}

	var respMap map[string]interface{}
	if err := json.Unmarshal(respData, &respMap); err != nil {
		return nil, err
	}

	// Modify the response to fit Todo/get format (only list and notFound)
	getResponse := map[string]interface{}{
		"accountId": "primary",
		"state":     "initial", // Simplified
		"notFound":  notFoundIDs,
		"list":      todoList,
	}

	return getResponse, nil
}

func handleTodoSet(db *sql.DB, args map[string]interface{}) (map[string]interface{}, error) {
	createdMap, hasCreated := args["create"].(map[string]interface{})
	updatedMap, hasUpdated := args["update"].(map[string]interface{})
	destroyedIDs, hasDestroyed := args["destroy"].([]interface{})

	setResponse := map[string]interface{}{
		"accountId":    "primary",
		"oldState":     "initial",       // Simplified
		"newState":     "updated-state", // Simplified
		"created":      map[string]interface{}{},
		"updated":      map[string]interface{}{},
		"destroyed":    []string{},
		"notCreated":   map[string]interface{}{},
		"notUpdated":   map[string]interface{}{},
		"notDestroyed": map[string]interface{}{},
	}

	if hasCreated {
		for clientID, createData := range createdMap {
			todoData, ok := createData.(map[string]interface{})
			if !ok {
				setResponse["notCreated"].(map[string]interface{})[clientID] = map[string]interface{}{
					"type":       "invalidProperties", // Simplified error type
					"properties": []string{"*"},
				}
				continue
			}
			title, ok := todoData["title"].(string)
			if !ok || title == "" {
				setResponse["notCreated"].(map[string]interface{})[clientID] = map[string]interface{}{
					"type":       "invalidProperties",
					"properties": []string{"title"},
				}
				continue
			}

			newTodo, err := CreateTodo(db, title)
			if err != nil {
				log.Printf("Error creating todo: %v", err)
				setResponse["notCreated"].(map[string]interface{})[clientID] = map[string]interface{}{
					"type": "serverFail", // Generic server error
				}
				continue
			}
			setResponse["created"].(map[string]interface{})[clientID] = newTodo
		}
	}

	if hasUpdated {
		for todoID, updateData := range updatedMap {
			updates, ok := updateData.(map[string]interface{})
			if !ok {
				setResponse["notUpdated"].(map[string]interface{})[todoID] = map[string]interface{}{
					"type":       "invalidProperties", // Simplified error type
					"properties": []string{"*"},
				}
				continue
			}

			updatedTodo, err := UpdateTodo(db, todoID, updates)
			if err != nil {
				log.Printf("Error updating todo %s: %v", todoID, err)
				setResponse["notUpdated"].(map[string]interface{})[todoID] = map[string]interface{}{
					"type": "serverFail", // Generic server error
				}
				continue
			}
			setResponse["updated"].(map[string]interface{})[todoID] = updatedTodo
		}
	}

	if hasDestroyed {
		for _, idInterface := range destroyedIDs {
			todoID, ok := idInterface.(string)
			if !ok {
				setResponse["notDestroyed"].(map[string]interface{})[todoID] = map[string]interface{}{
					"type": "invalidId", // Simplified error type
				}
				continue
			}
			err := DeleteTodo(db, todoID)
			if err != nil {
				log.Printf("Error deleting todo %s: %v", todoID, err)
				setResponse["notDestroyed"].(map[string]interface{})[todoID] = map[string]interface{}{
					"type": "serverFail", // Generic server error
				}
				continue
			}
			setResponse["destroyed"] = append(setResponse["destroyed"].([]string), todoID)
		}
	}

	return setResponse, nil
}
