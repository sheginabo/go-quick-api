# {repoName}

Desc about feature 
1. {Feature 1: xxx}
    - balabala

2. {Feature 2: xxx}
    - balabala

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- [Go](https://golang.org/dl/) (Version 1.23 or later)
- [Docker](https://www.docker.com/get-started)

### Clone (make sure branch is quick-websocket-chat)
```sh
git clone -b quick-websocket-chat https://github.com/sheginabo/go-quick-api.git
```

### Running Locally

To run the project locally for debugging purposes, follow these steps:

1. **準備開發環境**
   use `docker-compose.yaml`：
   ```sh
   docker-compose -p go-quick-api up -d
   ```
   If you want a very clean build, use the following command
   ```bash
   docker-compose build --no-cache && docker-compose -p go-quick-api up --force-recreate -d
   ```

2. Download the Go modules required by the project:
    ```sh
    go mod download
    ```
   
3. .env
   ```sh
   cp app.env .env
   ```

4. Run the application:
    ```sh
    go run ./cmd/main.go
    ```

### DEBUG
1. **check MongoDB**
    ```sh
    mongosh -u root -p example --authenticationDatabase admin
    ```