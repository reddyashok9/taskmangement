import axios from 'axios';

// API URLs (Replace these with your actual Go and Node.js server URLs)
const GO_API_URL = 'http://localhost:8080';
const NODE_API_URL = 'http://localhost:3001/api';

// Bulk create tasks using Go API
export const bulkCreateTasks = async (tasks) => {
    return axios.post(`${GO_API_URL}/tasks/bulk-create`, tasks);
};

// Create a single task using Node.js API
export const createTask = async (task) => {
    return axios.post(`${NODE_API_URL}/tasks`, task);
};

// Get all tasks using Go API
export const getAllTasks = async () => {
    return axios.get(`${GO_API_URL}/tasks`);
};

// Get a single task using Node.js API
export const getTaskById = async (id) => {
    return axios.get(`${NODE_API_URL}/tasks/${id}`);
};
