const express = require('express');
const dotenv = require('dotenv');
const cors = require('cors');
const connectDB = require('./config/db');
const taskRoutes = require('./routes/taskRoutes');

// Load environment variables
dotenv.config();

// Connect to MongoDB
connectDB();

const app = express();

// Enable CORS for all routes
app.use(cors({
    origin: 'http://localhost:3002', // Replace with the origin of your frontend (React app)
    methods: ['GET', 'POST', 'PUT', 'DELETE'],
    allowedHeaders: ['Authorization', 'Content-Type'],
    credentials: true // Allow cookies or credentials if needed
  }));


app.use(express.json()); // To parse JSON bodies

// Task routes
app.use('/api', taskRoutes);

// Server listening
const PORT = process.env.PORT || 3001;
app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
