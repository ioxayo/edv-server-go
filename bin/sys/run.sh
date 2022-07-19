#!/usr/bin/env bash

export PORT=5000

# Run server locally
runLocal() {
  ./bin/sys/exit.sh -l
  echo Building EDV server...
  go build
  echo Running EDV server...
  nohup ./edv-server-go > out.txt 2>&1 &
  echo $! > pid.txt
  echo EDV server successfully running on port $PORT!
}

# Run server in Docker
runDocker() {
  ./bin/sys/exit.sh -d
  docker compose build edv
  docker compose up -d edv
}

for i in "$@"; do
  case $i in
    -l|--local)
      runLocal
      ;;
    -d|--docker)
      runDocker
      ;;
    -*)
      echo "Invalid option: $i"
      exit 1
      ;;
    *)
      runLocal
      ;;
  esac
done
