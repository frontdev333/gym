#!/bin/bash
set -e

echo "1. Creating a User..."
USER_RESP=$(curl -s -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Ivan"}')
echo "User response: $USER_RESP"
USER_ID=$(echo $USER_RESP | grep -o '"id":"[^"]*' | grep -o '[^"]*$')
echo "User ID: $USER_ID"

echo -e "\n2. Creating Exercises..."
EX_RESP1=$(curl -s -X POST http://localhost:8080/api/v1/exercises \
  -H "Content-Type: application/json" \
  -d '{"title": "Bench Press"}')
echo "Exercise 1: $EX_RESP1"
EX_ID1=$(echo $EX_RESP1 | grep -o '"id":"[^"]*' | grep -o '[^"]*$')

EX_RESP2=$(curl -s -X POST http://localhost:8080/api/v1/exercises \
  -H "Content-Type: application/json" \
  -d '{"title": "Squat"}')
echo "Exercise 2: $EX_RESP2"
EX_ID2=$(echo $EX_RESP2 | grep -o '"id":"[^"]*' | grep -o '[^"]*$')

echo -e "\n3. Logging Workouts..."
curl -s -X POST http://localhost:8080/api/v1/users/${USER_ID}/workouts \
  -H "Content-Type: application/json" \
  -d "{\"exercise_id\": \"${EX_ID1}\", \"amount\": 10}"
echo -e "\nLogged workout 1"

curl -s -X POST http://localhost:8080/api/v1/users/${USER_ID}/workouts \
  -H "Content-Type: application/json" \
  -d "{\"exercise_id\": \"${EX_ID2}\", \"amount\": 15}"
echo -e "\nLogged workout 2"

# add a workout from yesterday
YESTERDAY=$(date -v-1d -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -d "yesterday" -u +"%Y-%m-%dT%H:%M:%SZ")
curl -s -X POST http://localhost:8080/api/v1/users/${USER_ID}/workouts \
  -H "Content-Type: application/json" \
  -d "{\"exercise_id\": \"${EX_ID1}\", \"amount\": 5, \"performed_at\": \"${YESTERDAY}\"}"
echo -e "\nLogged workout from yesterday"

echo -e "\n4. Fetching Statistics..."
curl -s http://localhost:8080/api/v1/users/${USER_ID}/statistics | jq .

