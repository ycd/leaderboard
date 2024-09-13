#!/bin/bash

set -e

# Function to generate a random username
generate_random_username() {
    echo "testuser-$(date +%s%N | cut -b1-13)"
}

# Function to create a user and return the user ID
create_user() {
    local username=$(generate_random_username)
    local response=$(curl -s -X POST 'http://localhost:8080/user/profile' \
        -H 'Content-Type: application/json' \
        -d "{\"username\": \"$username\", \"country\": \"US\"}")
    echo $response | jq -r '.id'
}

# Function to set user level
set_user_level() {
    local user_id=$1
    local level=$((RANDOM % 20 + 81))  # Random level between 81 and 100
    curl -s -X POST 'http://localhost:8080/user/level' \
        -H 'Content-Type: application/json' \
        -d "{\"user_id\": \"$user_id\", \"level\": $level}"
}

# Function to create an event
create_event() {
    local name="Event-$(date +%s)"
    local start_time=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    local end_time=$(date -u -v+30d +"%Y-%m-%dT%H:%M:%SZ")
    local response=$(curl -s -X POST 'http://localhost:8080/admin/event' \
        -H 'Content-Type: application/json' \
        -d "{\"name\": \"$name\", \"start_time\": \"$start_time\", \"end_time\": \"$end_time\"}")
    echo $response | jq -r '.id'
}

# Function to join an event
join_event() {
    local event_id=$1
    local user_id=$2
    curl -s -X POST 'http://localhost:8080/event/join' \
        -H 'Content-Type: application/json' \
        -d "{\"event_id\": \"$event_id\", \"user_id\": \"$user_id\"}"
}

# Function to update leaderboard progress
update_leaderboard() {
    local event_id=$1
    local user_id=$2
    local score=$((RANDOM % 10000 + 1))  # Random score between 1 and 10000
    curl -s -X POST 'http://localhost:8080/event/leaderboard/progress' \
        -H 'Content-Type: application/json' \
        -d "{\"event_id\": \"$event_id\", \"user_id\": \"$user_id\", \"score\": $score}"
}

# Main execution

echo "Creating events..."
event_ids=()
for i in {1..3}; do
    event_id=$(create_event)
    event_ids+=($event_id)
    echo "Created event with ID: $event_id"
done

echo "Creating users and joining events..."
user_count=${1:-5000}
user_ids=()
for i in $(seq 1 $user_count); do
    user_id=$(create_user)
    set_user_level $user_id
    user_ids+=($user_id)
    for event_id in "${event_ids[@]}"; do
        join_event $event_id $user_id
    done
    echo "Created and joined user $i/$user_count"
done

echo "Updating leaderboard..."
for i in {1..20000}; do
    random_user_index=$((RANDOM % ${#user_ids[@]}))
    random_event_index=$((RANDOM % ${#event_ids[@]}))
    update_leaderboard ${event_ids[$random_event_index]} ${user_ids[$random_user_index]}
    if [ $((i % 1000)) -eq 0 ]; then
        echo "Processed $i/20000 leaderboard updates"
    fi
done

echo "Integration test completed successfully!"
