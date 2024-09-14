#!/bin/bash

source /home/mohammad/Videos/go/proxy.env

# Define variables
LOGFILE="$HOME/Videos/go/Barcelona-watch/barcelona_watch_log.log"
SCRIPT_PATH="$HOME/Videos/go/Barcelona-watch" # Change this to your actual script directory

# Function to check proxy connection
check_proxy() {
    nc -zv "$PROXY_HOST" "$PROXY_PORT" >/dev/null 2>&1
    if [ $? -eq 0 ]; then
        return 0
    else
        return 1
    fi
}

# Function to check if Windscribe VPN is up
check_windscribe() {
    windscribe-cli status | grep -q "Connected"
    if [ $? -eq 0 ]; then
        return 0
    else
        return 1
    fi
}

# Function to check if the script has already run today
check_already_run() {
    TODAY=$(date '+%Y-%m-%d')
    echo "Checking if script has run today: $TODAY"

    if grep -q "$TODAY" "$LOGFILE"; then
        echo "Script has already run today ($TODAY), skipping."
        return 0 # Return true (0) if the script has run today
    else
        echo "Script has not run today yet."
        return 1 # Return false (1) if the script has not run today
    fi
}

# Function to log the run date
log_run_date() {
    TODAY=$(date '+%Y-%m-%d')
    echo "Logging the run date: $TODAY"

    # Append run date to the log file
    echo "Barcelona-watch ran on: $TODAY" >> "$LOGFILE"

    # Debug: check if the file was successfully written to
    if grep -q "$TODAY" "$LOGFILE"; then
        echo "Run date successfully logged."
    else
        echo "Failed to log run date. File contents:"
        cat "$LOGFILE"  # Output the current content of the file
    fi
}

# First, check if the script has already run today
if check_already_run; then
    exit 0  # Exit if the script has already run today
else
    echo "Proceeding with the script, as it has not run today."
fi

# Check Windscribe VPN first, then fallback to proxy logic
if check_windscribe; then
    echo "Windscribe VPN is up, running the barcelona-watch script without proxy..."
    cd "$SCRIPT_PATH" || { echo "Failed to change directory to $SCRIPT_PATH"; exit 1; }
    /home/mohammad/go/bin/barcelona
    GO_RUN_STATUS=$?

    if [ $GO_RUN_STATUS -eq 0 ]; then
        log_run_date
    else
        echo "Barcelona-watch script failed with status: $GO_RUN_STATUS"
    fi

elif check_proxy; then
    echo "Proxy is up, running the barcelona-watch script with proxy..."
    cd "$SCRIPT_PATH" || { echo "Failed to change directory to $SCRIPT_PATH"; exit 1; }
    /home/mohammad/go/bin/barcelona --proxy="$PROXY"
    GO_RUN_STATUS=$?

    if [ $GO_RUN_STATUS -eq 0 ]; then
        log_run_date
    else
        echo "Barcelona-watch script failed with status: $GO_RUN_STATUS"
    fi

else
    echo "Neither Windscribe VPN nor proxy is available. Skipping this attempt."
fi

