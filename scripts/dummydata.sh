#!/bin/bash
WAIT_TIME=0.01
endpoints=(
  "/api/books/"
  "/api/auth/loginStatus"
  "/api/auth/loginAccount"
  "/api/categories"
  "/api/formats"
  )

for i in {1..1000}; do
  for path in ${endpoints[@]}; do
    curl "http://localhost:8080$path" > /dev/null 2>&1;
  done
  sleep $WAIT_TIME;
done
