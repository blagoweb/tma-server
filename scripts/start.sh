#!/bin/bash

# Wait for database to be ready
echo "Waiting for database to be ready..."
sleep 5

# Run database migrations (if any)
echo "Running database setup..."

# Start the application
echo "Starting Telegram Mini App API server..."
exec ./main
