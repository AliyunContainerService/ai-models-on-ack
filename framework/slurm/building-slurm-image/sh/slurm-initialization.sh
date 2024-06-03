#!/bin/sh

set -e

mkdir -p /var/spool/slurmctld
chown slurm:slurm /var/spool/slurmctld
chmod 755 /var/spool/slurmctld
touch /var/log/slurm_jobacct.log /var/log/slurm_jobcomp.log
chown slurm:slurm /var/log/slurm_jobacct.log /var/log/slurm_jobcomp.log

echo 'SLURMD_OPTIONS="-Z"' >> /etc/default/slurmd
