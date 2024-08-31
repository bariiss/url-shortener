# URL Shortener
This project is a URL shortener service built with Go and using Redis as an in-memory store. The application allows users to create short URLs for long URLs and redirect to the original URL when the short URL is accessed.

## Project Structure
The project consists of the following key components:
- Go application files (including `main.go` and `handlers.go`)
- `Dockerfile` for containerizing the application
- `docker-compose.yml` for defining and running multi-container Docker applications
- `.env` file for environment variables
- `.dockerignore` file to exclude unnecessary files from the Docker build context

## Prerequisites
- Docker
- Docker Compose

## Configuration
### Environment Variables
The project uses a `.env` file to manage environment variables. Make sure to set the following variables:
```
APP_PORT=8080
EXT_PORT=8080
REDIS_ADDR=redis:6379
```

### .dockerignore
To optimize the Docker build process, we use a `.dockerignore` file with the following content:
```
docker-compose.yml
Dockerfile
.gitignore
.dockerignore
.git
redisdata/
```
This prevents unnecessary files from being included in the Docker build context.

## Dockerfile
The `Dockerfile` is configured to create a minimal image for the Go application. It uses a multi-stage build process:
1. Builds the Go application
2. Installs and uses UPX for binary compression
3. Creates a minimal scratch image with only the necessary components

Key steps in the Dockerfile:
- Uses `golang:1.23-bookworm` as the base image for building
- Installs UPX for binary compression
- Copies and builds the Go application
- Creates a minimal scratch image with the compiled binary and necessary static files

## Docker Compose
The `docker-compose.yml` file defines two services:
1. `app`: The URL shortener application
2. `redis`: The Redis in-memory store

The compose file sets up the necessary network connections and volume mounts.

## Building and Running
To build and run the application, use the following command:
```bash
docker-compose up --force-recreate --remove-orphans --build
```
This command will:
- Build the Docker images
- Create and start the containers
- Force recreation of containers and remove any orphaned containers
- Display the logs from both services

## Accessing the Application
Once the containers are up and running, you can access the URL shortener service at:
```
http://localhost:8080
```

## Troubleshooting
If you encounter any issues during the build process, ensure that:
1. All necessary Go files are present in your project directory
2. The `go.mod` and `go.sum` files are up to date
3. The `Dockerfile` and `docker-compose.yml` files are correctly configured
4. The `.env` file contains the correct environment variables

## Contributing
Feel free to submit issues and pull requests for any improvements or bug fixes.

## License
This project is open-source and available under the [MIT License](https://opensource.org/licenses/MIT).

Copyright (c) [2024] [bariiss]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
