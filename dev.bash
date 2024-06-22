if [ $1 = "build" ]; then
    echo "Building main.go..."
    go build -o bin/main main.go
elif [ $1 = "run" ]; then
    echo "Running main.go..."
    go run main.go
else 
    echo "Parameter error!"
fi