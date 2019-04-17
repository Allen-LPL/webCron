#!/bin/bash
case $1 in 
	start)
		nohup /usr/local/goWorkSpace/bin/bee run /usr/local/goWorkSpace/src/webcron >/dev/null 2>&1 &
		echo "服务已启动..."
		sleep 1
	;;
	stop)
		ps -ef |grep webcron|grep -v grep|awk '{print $2}'|xargs kill -9
		echo "服务已停止..."
		sleep 1
	;;
	restart)
		ps -ef |grep webcron|grep -v grep|awk '{print $2}'|xargs kill -9
		sleep 1
		nohup /usr/local/goWorkSpace/bin/bee run /usr/local/goWorkSpace/src/webcron >/dev/null 2>&1 &
		echo "服务已重启..."
		sleep 1
	;;
	*) 
		echo "$0 {start|stop|restart}"
		exit 4
	;;
esac

