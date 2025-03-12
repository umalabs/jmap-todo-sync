import React, { useEffect, useState } from 'react';
import { todoQuery, todoSet } from './jmapClient';

function TodoList() {
    const [todos, setTodos] = useState([]);
    const [newTodoTitle, setNewTodoTitle] = useState('');

    useEffect(() => {
        fetchTodos();
    }, []);

    const fetchTodos = async () => {
        try {
            const response = await todoQuery();
            if (response.methodResponses && response.methodResponses.length > 0) {
                const getResponse = response.methodResponses.find(resp => resp[0] === 'Todo/get');
                if (getResponse && getResponse[1] && getResponse[1].list) {
                    setTodos(getResponse[1].list);
                }
            }
        } catch (error) {
            console.error("Failed to fetch todos:", error);
        }
    };

    const handleAddTodo = async () => {
        if (newTodoTitle.trim() === '') return;

        const createArgs = {
            [`${Date.now()}`]: { title: newTodoTitle } // Use timestamp as temporary client-side ID
        };

        try {
            await todoSet({ create: createArgs });
            setNewTodoTitle('');
            fetchTodos(); // Re-fetch todos after adding
        } catch (error) {
            console.error("Failed to add todo:", error);
        }
    };

    const handleToggleComplete = async (todoId, isCompleted) => {
        try {
            await todoSet({ update: { [todoId]: { isCompleted: !isCompleted } } });
            fetchTodos(); // Re-fetch todos after updating
        } catch (error) {
            console.error("Failed to toggle todo completion:", error);
        }
    };

    const handleDeleteTodo = async (todoId) => {
        try {
            await todoSet({ destroy: [todoId] });
            fetchTodos(); // Re-fetch todos after deleting
        } catch (error) {
            console.error("Failed to delete todo:", error);
        }
    };


    return (
        <div>
            <input
                type="text"
                placeholder="New todo title"
                value={newTodoTitle}
                onChange={(e) => setNewTodoTitle(e.target.value)}
            />
            <button onClick={handleAddTodo}>Add Todo</button>

            <ul>
                {todos.map(todo => (
                    <li key={todo.id}>
                        <input
                            type="checkbox"
                            checked={todo.isCompleted}
                            onChange={() => handleToggleComplete(todo.id, todo.isCompleted)}
                        />
                        <span style={{ textDecoration: todo.isCompleted ? 'line-through' : 'none' }}>
                            {todo.title}
                        </span>
                        <button onClick={() => handleDeleteTodo(todo.id)}>Delete</button>
                    </li>
                ))}
            </ul>
        </div>
    );
}

export default TodoList;