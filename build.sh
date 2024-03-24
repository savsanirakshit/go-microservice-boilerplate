#!/bin/bash

# Run go get
echo "Running go get..."
go get

# Check if go get was successful
if [ $? -ne 0 ]; then
    echo "Error: go get failed. Exiting..."
    exit 1
fi

# Run go build
echo "Running go build..."
go build

# Check if go build was successful
if [ $? -ne 0 ]; then
    echo "Error: go build failed. Exiting..."
    exit 1
fi

# List of files to include in the zip
binary_file="golang-microservice-boilerplate"
environment_file="app.config"
folder_to_copy="db"
directory1="filedb"
directory2="logs"

# Create a temporary directory to store files
temp_dir="microservice-service"
mkdir -p "$temp_dir"

cp "$binary_file" "$temp_dir"
cp "$environment_file" "$temp_dir"
cp -r "$folder_to_copy" "$temp_dir"
mkdir -p "$temp_dir/$directory1"
mkdir -p "$temp_dir/$directory2"

# Create a zip file
echo "Creating service zip file..."
zip -r service.zip "$temp_dir"

# Check if zip creation was successful
if [ $? -ne 0 ]; then
    echo "Error: Failed to create zip file. Exiting..."
    exit 1
fi

rm -rf "$temp_dir"
echo "Build successful!"

