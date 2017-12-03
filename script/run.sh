#!/bin/bash

RUNBIN="${BASH_SOURCE-$0}"
RUNBIN="$(dirname "${RUNBIN}")"
BINDIR="$(cd "${RUNBIN}"; pwd)"

# parse gateway config file
config_file="/opt/galaxy/galaxy-s3-gw/bin/galaxy-s3-gateway.cfg"
if [ ! -f "$config_file" ]; then
    echo 's3 gateway config file '/opt/galaxy/galaxy-s3-gw/bin/galaxy-s3-gateway.cfg' not exist'
    exit 1
fi

while read line; do
    if [[ ${line:0:1} != "#" ]]; then
        eval "$line"
    fi
done < $config_file

mkdir -p $log_dir

cd $BINDIR && nohup ./galaxy-s3-gateway -gfs_zk_addr=$zookeeper -log_dir=$log_dir -logtostderr=false -mongodb_addr=$mongodb_address -port=$listen_port > nohup.out 2>&1 &

PID=$!
sleep 1
echo "start galaxy-s3-gateway finished, PID=$PID"
echo "checking if $PID is running..."
sleep 2
kill -0 $PID > /dev/null 2>&1
if [ $? -eq 0 ]
then
	echo "$PID is running, start galaxy-s3-gateway success."
	exit 0
else
	echo "start galaxy-s3-gateway failed."
	exit 1
fi
