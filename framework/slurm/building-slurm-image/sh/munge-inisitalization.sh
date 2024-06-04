#!/bin/sh

mkdir -p /etc/munge/ /var/log/munge/ 
chown -R munge: /etc/munge/ /var/log/munge/ 
chmod 0700 /etc/munge/ /var/log/munge/
mkdir -p /run/munge
chmod 0755 /run/munge
chown -R munge:root /run/munge