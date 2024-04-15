#!/bin/bash

# Set the environment variable for the Firestore Emulator
echo "Setting up environment variable for Firestore Emulator..."
export FIRESTORE_EMULATOR_HOST="localhost:8081"

# Start Firestore Emulator
echo "Starting Firestore Emulator..."
firebase emulators:start --only firestore &
EMULATOR_PID=$!

# Wait for the emulator to fully start
echo "Waiting for the Firestore Emulator to fully start..."
sleep 5

# Populate emulator with data
echo "Populating Firestore Emulator with test data..."
go run ./C:/Users/Mariu/GolandProjects/countries-dashboard-service/firestoreEmulator/populateEmulatedFirestore.go

# Uncomment the following lines when you have tests to run
#echo "Running tests..."
#go test ./...

# Kill the Firestore Emulator after tests are done
echo "Killing the Firestore Emulator..."
kill $EMULATOR_PID

# Unset the environment variable
echo "Cleaning up environment..."
unset FIRESTORE_EMULATOR_HOST

echo "Test environment cleaned up and Firestore Emulator shut down."