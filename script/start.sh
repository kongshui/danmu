#!/bin/bash

workDir="/data/web"
gamed="web"
timeStr=`date +%Y%m%d`
logFile="${workDir}/logs/${gamed}.${timeStr}.log"
execFile="${workDir}/${gamed}"
chmod +x $execFile
ePid=`/usr/bin/lsof web 2> /dev/null |grep -v "PID" |awk '{print $2}'`
if [[ ${ePid} != "" ]]
then
    kill -9 ${ePid}
fi
sleep 2
if [ -f ${logFile} ]
then
    mv ${logFile} ${logFile}.`date +%s`
fi

$execFile &> ${logFile}
