#!/bin/bash

edv=$(netstat -vanp tcp | grep $EDV_PORT | awk -F ' ' '{print $9}')
echo Killing EDV server process $edv...
kill -9 $edv
echo Deleting EDV server metadata files...
rm edv-server-go
rm out.txt
rm pid.txt
echo EDV server successfully terminated!
