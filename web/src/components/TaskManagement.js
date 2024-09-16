import React, { useState, useEffect } from 'react';
import { TextField, Button, Container, Typography, List, ListItem, ListItemText, Box } from '@mui/material';
import { bulkCreateTasks, createTask, getAllTasks, getTaskById } from '../tasksapi';

const TaskManagement = () => {
    const [tasks, setTasks] = useState([]);
    const [newTask, setNewTask] = useState({ title: '', description: '' });
    const [taskId, setTaskId] = useState('');
    const [singleTask, setSingleTask] = useState(null);

    // Fetch all tasks when the component mounts
    useEffect(() => {
        fetchTasks();
    }, []);

    // Fetch all tasks from Go API
    const fetchTasks = async () => {
        try {
            const response = await getAllTasks();
            setTasks(response.data.tasks); // Assume 'tasks' is the field in the response
        } catch (error) {
            console.error('Error fetching tasks:', error);
        }
    };

    // Handle bulk create using Go API
    const handleBulkCreate = async () => {
        const bulkTasks = [
            { title: 'Task 1', description: 'Bulk task 1' },
            { title: 'Task 2', description: 'Bulk task 2' },
            // Add more tasks if needed
        ];
        try {
            await bulkCreateTasks(bulkTasks);
            fetchTasks(); // Refresh the tasks list
        } catch (error) {
            console.error('Error creating tasks in bulk:', error);
        }
    };

    // Handle create single task using Node.js API
    const handleCreateTask = async () => {
        try {
            await createTask(newTask);
            setNewTask({ title: '', description: '' }); // Reset input fields
            fetchTasks(); // Refresh the tasks list
        } catch (error) {
            console.error('Error creating task:', error);
        }
    };

    // Handle fetching a single task by ID using Node.js API
    const handleGetTaskById = async () => {
        try {
            const response = await getTaskById(taskId);
            setSingleTask(response.data); // Assume task data is returned
        } catch (error) {
            console.error('Error fetching task:', error);
        }
    };

    return (
        <Container>
            <Typography variant="h4" gutterBottom>Task Management</Typography>

            {/* Bulk Create Tasks */}
            <Box my={2}>
                <Button variant="contained" color="primary" onClick={handleBulkCreate}>
                    Bulk Create Tasks (Go API)
                </Button>
            </Box>

            {/* Create Single Task */}
            <Box my={2}>
                <Typography variant="h6">Create a New Task (Node.js API)</Typography>
                <TextField
                    label="Title"
                    value={newTask.title}
                    onChange={(e) => setNewTask({ ...newTask, title: e.target.value })}
                    fullWidth
                    margin="normal"
                />
                <TextField
                    label="Description"
                    value={newTask.description}
                    onChange={(e) => setNewTask({ ...newTask, description: e.target.value })}
                    fullWidth
                    margin="normal"
                />
                <Button variant="contained" color="primary" onClick={handleCreateTask}>
                    Create Task
                </Button>
            </Box>

            {/* Get All Tasks */}
            <Box my={2}>
                <Typography variant="h6">All Tasks (Go API)</Typography>
                <List>
                    {tasks?.map((task) => (
                        <ListItem key={task._id}>
                            <ListItemText primary={task.title} secondary={task.description} />
                        </ListItem>
                    ))}
                </List>
            </Box>

            {/* Get Single Task by ID */}
            <Box my={2}>
                <Typography variant="h6">Get Task by ID (Node.js API)</Typography>
                <TextField
                    label="Task ID"
                    value={taskId}
                    onChange={(e) => setTaskId(e.target.value)}
                    fullWidth
                    margin="normal"
                />
                <Button variant="contained" color="primary" onClick={handleGetTaskById}>
                    Get Task
                </Button>
                {singleTask && (
                    <Box my={2}>
                        <Typography variant="body1">Title: {singleTask.title}</Typography>
                        <Typography variant="body2">Description: {singleTask.description}</Typography>
                    </Box>
                )}
            </Box>
        </Container>
    );
};

export default TaskManagement;
