#!/bin/bash
cd "$(dirname "$0")"
processname=homeip
# test $( pgrep -x homip | wc -l ) -gt 0 && echo "OK" || echo "Not"
if test $( pgrep -x ${processname} | wc -l ) -gt 0
then
    echo "${processname} already running"
    exit 2
else
    echo "${processname} start running..."
fi
./homeip 1>>access.log 2>error.log &
echo $! > homeipkit.pid
# pgrep -lnx ${processname}
if test $(cat error.log | wc -l) -eq 0;then
    exit
else
    cat error.log
    exit 1
fi