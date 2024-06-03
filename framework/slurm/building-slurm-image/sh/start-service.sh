#!/bin/bash

cat /etc/secret-volume/.munge-key > /etc/munge/munge.key
service munge start
cp -L -r /var/slurm/* /etc/slurm
if [ -f "/etc/slurm/slurmdbd.conf" ]; then
    chmod 0600 "/etc/slurm/slurmdbd.conf"
    chown slurm: "/etc/slurm/slurmdbd.conf"
fi

stateSaveLocation=${StateSaveLocation:-/var/spool/slurm}
mkdir -p "$stateSaveLocation"/cluster_state
chown -R slurm:slurm "$stateSaveLocation"
chmod 755 "$stateSaveLocation"

slurmSpoolDir=${SlurmdSpoolDir:-/var/spool/slurm/d}
mkdir -p "$slurmSpoolDir"
chmod 755 "$slurmSpoolDir"
chown -R slurm: "$slurmSpoolDir"
slurmctldLogFile=${SlurmctldLogFile:-/var/log/slurmctld.log}
touch "$slurmctldLogFile"
chown slurm:slurm "$slurmctldLogFile"
slurmdLog=${SlurmdLogFile:-/var/log/slurmd.log}
touch "$slurmdLog"
chown slurm: "$slurmdLog"

role=$1
if [ "$role" = "master" ]; then
    if [ -f "/etc/slurm/slurmdbd.conf" ]; then
        service slurmdbd start
        sleep 5
    fi
    service slurmctld start
    service prometheus-slurm-exporter start
fi
if [ "$role" = "worker" ]; then
    CPUs=1
    RealMemory=1024
    if [ -f "/etc/slurmd-podinfo/cpu_limit" ]; then
        CPUs=$(cat /etc/slurmd-podinfo/cpu_limit)
    fi
    if [ -f "/etc/slurmd-podinfo/memory_limit" ]; then
        RealMemory=$(cat /etc/slurmd-podinfo/memory_limit)
    fi


    nvidia_visible_devices=$NVIDIA_VISIBLE_DEVICES
    if [ $# -ge 4 ]; then
        Partition="$4"
        default_option="SLURMD_OPTIONS='-Z --conf=\"Feature=$Partition CPUs=$CPUs RealMemory=$RealMemory\"'"
    else
        default_option="SLURMD_OPTIONS='-Z --conf=\"CPUs=$CPUs RealMemory=$RealMemory\"'"
    fi
    if [ -n "$nvidia_visible_devices" ]; then
        if ! command -v nvidia-smi > /dev/null; then
          echo "Error: nvidia-smi command not found."
        else
          gpu_ids=$(echo "$nvidia_visible_devices" | tr ',' '\n')
    
          gpu_list=""
          for gpu_id in $gpu_ids; do
              gpu_name=$(nvidia-smi --id="$gpu_id" --query-gpu=name --format=csv,noheader | tr ' ' '_')
              modified_gpu_name="gpu:${gpu_name}:1"
              gpu_list="$gpu_list,$modified_gpu_name"
          done
    
          gpu_list=${gpu_list:1}  # Remove leading comma
          echo "GPU list: $gpu_list"
          default_option="SLURMD_OPTIONS='-Z --conf=\"Feature=$Partition Gres=$gpu_list CPUs=$2 RealMemory=$3\"'"
        fi
    else
        echo "NVIDIA_VISIBLE_DEVICES is empty."
    fi

    echo $default_option > /etc/default/slurmd
    service slurmd start
fi

while true;
do
        sleep 1000
done
