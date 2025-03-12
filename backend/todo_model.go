package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Todo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	IsCompleted bool   `json:"isCompleted"`
}

func GetTodos(db *sql.DB) ([]Todo, error) {
	rows, err := db.Query("SELECT id, title, is_completed FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var todo Todo
		var isCompletedInt int
		if err := rows.Scan(&todo.ID, &todo.Title, &isCompletedInt); err != nil {
			return nil, err
		}
		todo.IsCompleted = isCompletedInt == 1
		todos = append(todos, todo)
	}
	return todos, nil
}

func GetTodo(db *sql.DB, id string) (*Todo, error) {
	row := db.QueryRow("SELECT id, title, is_completed FROM todos WHERE id = ?", id)
	var todo Todo
	var isCompletedInt int
	err := row.Scan(&todo.ID, &todo.Title, &isCompletedInt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("todo not found: %s", id)
		}
		return nil, err
	}
	todo.IsCompleted = isCompletedInt == 1
	return &todo, nil
}

func CreateTodo(db *sql.DB, title string) (*Todo, error) {
	id := uuid.New().String()
	_, err := db.Exec("INSERT INTO todos (id, title) VALUES (?, ?)", id, title)
	if err != nil {
		return nil, err
	}
	return GetTodo(db, id) // Fetch the newly created todo to return
}

func UpdateTodo(db *sql.DB, id string, updates map[string]interface{}) (*Todo, error) {
	if len(updates) == 0 {
		return GetTodo(db, id) // No updates, just return current todo
	}

	var setClauses []string
	var values []interface{}
	for key, value := range updates {
		switch key {
		case "title":
			setClauses = append(setClauses, "title = ?")
			values = append(values, value)
		case "isCompleted":
			setClauses = append(setClauses, "is_completed = ?")
			if boolValue, ok := value.(bool); ok {
				if boolValue {
					values = append(values, 1)
				} else {
					values = append(values, 0)
				}
			} else {
				return nil, fmt.Errorf("invalid type for isCompleted, expected boolean")
			}
		default:
			return nil, fmt.Errorf("unknown field: %s", key)
		}
	}

	setQuery := ""
	for i, clause := range setClauses {
		setQuery += clause
		if i < len(setClauses)-1 {
			setQuery += ", "
		}
	}
	values = append(values, id) // Add ID for WHERE clause

	query := fmt.Sprintf("UPDATE todos SET %s WHERE id = ?", setQuery)
	_, err := db.Exec(query, values...)
	if err != nil {
		return nil, err
	}
	return GetTodo(db, id)
}

func DeleteTodo(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}

func MarshalTodosToJMAPResponse(todos []Todo) ([]byte, error) {
	methodResponses := [][]interface{}{
		{"Todo/query", map[string]interface{}{
			"accountId":           "primary",  // In real JMAP, accountId is important
			"queryState":          "initial",  // Simplified, not really used here
			"canCalculateChanges": false,      // Simplified
			"ids":                 []string{}, // We'll populate this below
		}, "r1"},
		{"Todo/get", map[string]interface{}{
			"accountId": "primary",
			"state":     "initial", // Simplified
			"notFound":  []string{},
			"list":      []Todo{}, // We'll populate this below
		}, "r2"},
	}

	ids := []string{}
	todoList := []Todo{}
	for _, todo := range todos {
		ids = append(ids, todo.ID)
		todoList = append(todoList, todo)
	}

	methodResponses[0][1].(map[string]interface{})["ids"] = ids
	methodResponses[1][1].(map[string]interface{})["list"] = todoList

	response := map[string]interface{}{
		"methodResponses": methodResponses,
		"sessionState":    "initial-session-state", // Simplified
	}
	return json.Marshal(response)
}

func MarshalTodoToJMAPSetResponse(todo *Todo, created bool) ([]byte, error) {
	methodResponses := [][]interface{}{
		{"Todo/set", map[string]interface{}{
			"accountId":    "primary",
			"oldState":     "initial",       // Simplified
			"newState":     "updated-state", // Simplified
			"created":      map[string]interface{}{},
			"updated":      map[string]interface{}{},
			"destroyed":    []string{},
			"notCreated":   map[string]interface{}{},
			"notUpdated":   map[string]interface{}{},
			"notDestroyed": []string{},
		}, "r1"},
	}

	if created {
		methodResponses[0][1].(map[string]interface{})["created"] = map[string]interface{}{
			todo.ID: todo,
		}
	} else {
		methodResponses[0][1].(map[string]interface{})["updated"] = map[string]interface{}{
			todo.ID: todo,
		}
	}

	response := map[string]interface{}{
		"methodResponses": methodResponses,
		"sessionState":    "updated-session-state", // Simplified
	}

	return json.Marshal(response)
}

func MarshalTodoDeletionJMAPResponse(todoID string) ([]byte, error) {
	methodResponses := [][]interface{}{
		{"Todo/set", map[string]interface{}{
			"accountId":    "primary",
			"oldState":     "initial",       // Simplified
			"newState":     "updated-state", // Simplified
			"created":      map[string]interface{}{},
			"updated":      map[string]interface{}{},
			"destroyed":    []string{todoID},
			"notCreated":   map[string]interface{}{},
			"notUpdated":   map[string]interface{}{},
			"notDestroyed": map[string]interface{}{},
		}, "r1"},
	}

	response := map[string]interface{}{
		"methodResponses": methodResponses,
		"sessionState":    "updated-session-state", // Simplified
	}
	return json.Marshal(response)
}
