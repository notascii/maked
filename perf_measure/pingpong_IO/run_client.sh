build_go_files() {
    echo "Building Go files..."
    go build -o client/client client/client.go
    echo "Go file built successfully."
}
run_client_file() {

    echo "Running Go client"
    ./client/client $1 # Run client in the background
}

# Function to run the Python file
run_python_file() {
    echo "Running Python file..."
    python_file="./perf/perf_benchmarks/plot_graphs.py"
    python_file_table="./perf/perf_benchmarks/plot_tables.py"
    if [[ -f "$python_file" ]]; then
        python3 "$python_file"
	python3 "$python_file_table"
        echo "Python file executed successfully."
    else
        echo "Python file not found: $python_file"
        exit 1
    fi
}

build_go_files
run_client_file $1
run_python_file
