#!/bin/sh

./exit.sh
export EDV_PORT=5000
echo Building EDV server...
go build
echo Running EDV server...
nohup ./edv-server-go > out.txt 2>&1 &
echo $! > pid.txt
echo EDV server successfully running on port $EDV_PORT!
