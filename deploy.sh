#!/bin/bash
workDir=$(cd `dirname $0`; pwd)
port=1209
cd ${workDir}
tarname=${PWD##*/}\.tar\.gz
projectname=${PWD##*/}
go get
bee pack
[ $? -ne 0 ] && exit 1
tar zxf ${tarname} && rm -f ${tarname}
chmod +x ${projectname}
pid=`netstat -lntp|grep -w ${port}|awk '{print $NF}'|awk -F "/" '{print $1}' `
[ -n  "$pid" ] && kill -HUP ${pid}||exit 0