#!/bin/bash

case $1 in
    start)
        nohup /usr/bin/prometheus-slurm-exporter > /var/prometheus-slurm-exporter.log 2>&1 &
        echo $! > /var/prometheus-slurm-exporter.pid
        ;;
    stop)
        pid=$(cat /var/prometheus-slurm-exporter.pid)
        echo "killing $pid"
        kill -9 $pid
            rm /var/prometheus-slurm-exporter.pid
        ;;
    restart)
        cat /var/prometheus-slurm-exporter.pid | kill -9
        nohup /usr/bin/prometheus-slurm-exporter > /var/prometheus-slurm-exporter.log 2>&1 &
        ;;
    status)
        if [ -f /var/prometheus-slurm-exporter.pid ]; then
            pid=$(cat /var/prometheus-slurm-exporter.pid)
            if [[ $(ps -p $pid|sed 1d) == "" ]]; then
                echo "daemon is dead but pid file exists"
                exit 3
            fi
            echo "daemon is running"
            exit 0
        fi
        echo "daemon is not running"
        exit 3
        ;;
    *)
        echo "Usage: $0 {start|status|stop|restart}"
        exit 1
        ;;
esac

exit 0