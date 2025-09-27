# User Appointment Management System

This is a full-stack web application for managing user appointments. It features a backend built with Go, a MongoDB database, and a vanilla JavaScript frontend, with the entire environment containerized by Docker for easy setup and execution. The system handles user authentication, profile management, and appointment scheduling with intelligent timezone and working-hour validation.

---

## ✨ Key Features

- **Secure JWT Authentication**: Simple and secure login using a username with JWT-based authentication. Sessions expire after 1 hour.
- **Dynamic Signup Flow**: If a user attempts to log in with a username that doesn't exist, they are seamlessly redirected to a dedicated page to complete their profile.
- **User Profile Management**: Logged-in users can view and update their profile information, including their full name and preferred timezone.
- **Intelligent Appointment Scheduling**:
  - Users can create appointments for themselves or invite other users.
  - The system validates that the proposed time falls within the working hours (08:00 - 17:00) of all **invited** participants, based on their individual timezones.
  - Appointment times are always displayed to the logged-in user in their own preferred timezone for maximum readability.
- **Fully Containerized**: The entire application (Go backend, MongoDB database) runs in Docker containers, managed by a single `docker-compose` command.

---

## 💻 Tech Stack

- **Backend**: Go (Golang) with the Gin web framework
- **Database**: MongoDB (NoSQL)
- **Frontend**: HTML5, CSS3, Vanilla JavaScript
- **Containerization**: Docker & Docker Compose

---

## ✅ Requirements

To run this project, you will only need two pieces of software installed on your machine:

- **Docker**
- **Docker Compose**

> **Note**: If you are on Windows or macOS, installing **Docker Desktop** will provide you with both Docker and Docker Compose.

---

## 🚀 Installation and Tutorial

Follow these simple steps to get the application running on your local machine.

### 1. Clone the Repository

First, clone this project's repository to your local computer.

```bash
git clone https://github.com/Ezekielna70/appointment-system.git
cd appointment-system

```

### 2. Create the Environment File

The application requires an environment file to store database credentials and secrets.

Create a new file named .env in the root of the appointment-system folder and paste the following content into it:

```bash
MONGO_URI=mongodb://root:password@mongo:27017/
MONGO_DB_NAME=appointmentDB
JWT_SECRET=your_super_secret_key
MONGO_USER=root
MONGO_PASS=password
```

### 3. Run the Application

With Docker Desktop running, execute the following command from the root of the appointment-system folder. This single command will build the Go application, start the MongoDB container, and connect everything together.

```bash

docker-compose up --build

```

The terminal will show logs from both the database and the backend. Wait until you see the message:
Server starting on port 8080...

### 4. Use the Application

Your User Appointment Management System is now live!

Access the Login Page: Open your web browser and navigate to http://localhost:8080.

Sign Up: Try to log in with a new username (e.g., "testuser"). The system will not find this user and will automatically redirect you to the "Complete Your Profile" page. Fill out the form to create your account.

Use the Dashboard: Once you are logged in, you will be on the dashboard where you can:

Update your name and timezone in the "My Profile" section.

Create new appointments for yourself or invite other pre-loaded users (alice, bob, charlie).

View your upcoming appointments.

📚 API Documentation
This project includes a detailed **OpenAPI 3.0 (Swagger)** specification located at `documentation/openapi.yml` that can be used with tools like Apidog or Postman. This file describes all available API endpoints, expected request bodies, and example responses, making it easy to test and understand the backend.
