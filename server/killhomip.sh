#!/bin/bash
cd "$(dirname "$0")"
processname=homeip
pgrep -lnf ${processname}
if test $( pgrep -f ${processname} | wc -l ) -eq 0
then
    echo "${processname} is not running"
    exit
else
    printf "kill ${processname}...\nwait 2s\n"
fi
pkill -F homeipkit.pid
sleep 2s
rm homeipkit.pid
# pgrep -lnf ${processname}
