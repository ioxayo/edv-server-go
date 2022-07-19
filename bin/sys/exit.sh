#!/usr/bin/env bash

export PORT=5000

# Exit local server
exitLocal() {
  EDV=$(netstat -vanp tcp | grep $PORT | awk -F ' ' '{print $9}')
  echo Killing EDV server process $EDV...
  kill -9 $EDV
  echo Deleting EDV server metadata files...
  rm edv-server-go
  rm out.txt
  rm pid.txt
  echo EDV server successfully terminated!
}

# Exit Docker server
exitDocker() {
  docker compose stop edv
  docker compose kill edv
  docker compose rm -f edv
}

for i in "$@"; do
  case $i in
    -l|--local)
      exitLocal
      ;;
    -d|--docker)
      exitDocker
      ;;
    -*)
      echo "Invalid option: $i"
      exit 1
      ;;
    *)
      exitLocal
      ;;
  esac
done
