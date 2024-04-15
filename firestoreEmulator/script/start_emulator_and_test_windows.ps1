# PowerShell Script

# Set the environment variable for the Firestore Emulator
Write-Host "Setting up environment variable for Firestore Emulator..."
$env:FIRESTORE_EMULATOR_HOST = "localhost:8081"

# Start Firestore Emulator
Write-Host "Starting Firestore Emulator..."
$EmulatorProcess = Start-Process -FilePath "C:\Users\Mariu\GolandProjects\countries-dashboard-service\firestoreEmulator\emulatorFiles" -ArgumentList "emulators:start --only firestore --project countries-dashboard-service" -PassThru -NoNewWindow

# Wait for the emulator to fully start
Write-Host "Waiting for the Firestore Emulator to fully start..."
Start-Sleep -Seconds 5

# Populate emulator with data
Write-Host "Populating Firestore Emulator with test data..."
& go run "C:/Users/Mariu/GolandProjects/countries-dashboard-service/firestoreEmulator/populateEmulatedFirestore.go"

# Uncomment the following lines when you have tests to run
#Write-Host "Running tests..."
##& go test "./..."

# Kill the Firestore Emulator after tests are done
Write-Host "Killing the Firestore Emulator..."
$EmulatorProcess | Stop-Process

# Unset the environment variable
Write-Host "Cleaning up environment..."
Remove-Item Env:FIRESTORE_EMULATOR_HOST

Write-Host "Test environment cleaned up and Firestore Emulator shut down."