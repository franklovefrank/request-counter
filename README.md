# Go HTTP Server with Moving Window Request Counter

This Go application implements a simple HTTP server that maintains a moving window request counter. The server responds to each request with the total number of requests it has received during the previous 60 seconds.

## Table of Contents

- [Features](#features)
- [Usage](#usage)
- [Testing](#testing)
- [Implementation Details](#implementation-details)

## Features

- **Moving Window Counter:** The server keeps track of the total number of requests in the last 60 seconds.

- **Persistence:** Data is persisted to a file (`request_counter.gob`) to ensure that the server can recover the correct count even after restarts.

## Usage

1. Run the server:

    ```bash
    go run main.go
    ```

2. Access the server at [http://localhost:8080](http://localhost:8080).

3. The server will respond with a message indicating the total number of requests in the last 60 seconds.

4. To gracefully shut down the server and save counter data, use `Ctrl+C`.

## Testing

The application includes unit tests to ensure the correctness of the moving window counter and data persistence functionalities.

Run tests using the following command:

```bash
go test
