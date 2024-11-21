build_go_files() {
    echo "Building Go files..."
    go build -o server/server server/server.go
    echo "Go file built successfully."
}
run_server() {

    echo "Running Go server"
    ./server/server & # Run server in the background
    server_pid=$!
    echo "Server running with PID: $server_pid"
}

build_go_files
run_server