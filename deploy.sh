#!/bin/bash

# build new program
export PATH=$PATH:/usr/local/go/bin
go build main.go

# kill old program
string=`ps -ef | grep "./main" | grep -v "grep" | tr -s " " | cut -d " " -f 2`
array=(${string//,/ })

for var in ${array[@]}
do
   sudo kill -9 $var
done

# run new program
ulimit -n 100000

nohup sudo ./main -p 443 &
echo $! > pid.txt
