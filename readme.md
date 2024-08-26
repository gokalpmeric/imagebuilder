**Docker Image Builder API**

This project provides a REST API to build Docker images from a given Git repository and Dockerfile. The API is built using Go and the Gin framework.
Prerequisites
Go (latest version)
Docker
Git
Getting Started
Clone the Repository
git clone <repository-url>
cd <repository-directory>

**Build and Run the Application**

go build -o main .
./main

The server will start at http://localhost:8080.  


**API Endpoints**


Build Docker Image
URL: /build
Method: POST
Content-Type: application/json
Request Body:
{
"repo_url": "https://github.com/user/repo.git",
"dockerfile_path": "path/to/Dockerfile"
}
Response:
Success: 200 OK
{
"message": "Docker image built successfully"
}
Error: 400 Bad Request or 500 Internal Server Error
{
"error": "Error message"
}


License
This project is licensed under the MIT License.