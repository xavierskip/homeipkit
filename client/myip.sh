#!/bin/bash
# url="127.0.0.1:9090/myip"
url=$1  # your post url
key=$2  # your key
auth="auth:auth"  # change it if you used
step=20  # as your server Settings

unixtime=`date +%s`
logtime=`date +"%Y/%m/%d %H:%M:%S"`
steptime=$[unixtime/step]
# echo ${key}${steptime}
# echo -n ${key}${steptime}|shasum -a 256|cut -c1-64
homekey=`echo -n ${key}${steptime}|shasum -a 256|cut -c1-64`
# -v
# --ignore-stdin for crontab
# http -a home:ipkit -f PUT ${url} homekey=${homekey}
# &> /dev/null
if http --check-status --ignore-stdin --timeout=5 -a ${auth} -f PUT ${url} homekey=${homekey} &> /dev/null; then
    echo "${logtime} OK"
else
    case $? in
        2) echo "${logtime} Request timed out" ;;
        3) echo "${logtime} Unexpected HTTP 3xx Redirectio" ;;
        4) echo "${logtime} HTTP 4xx Client Error! unixtime:${unixtime} steptime:${steptime}" ;;
        5) echo "${logtime} HTTP 5xx Server Erro" ;;
        6) echo "${logtime} Exceeded --max-redirects=<n> redirect" ;;
        *) echo "${logtime} Other Error" ;;
    esac
fi
