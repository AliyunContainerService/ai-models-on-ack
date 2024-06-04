Optional, you can define a special initialization script by setting the value of initBash in slurmCluster CR. 

We will, by default, set up an initBash for the Headpod as well as workpods. Concatenating the content set by users in .containers[].command and .containers[].args and append the content after the initBash. Last, we execute the `sleep inf` command.

The default initBash for head pod is
``` bash
mkdir /var/k8s-info
echo "$KUBERNETES_SERVICE_HOST" > /var/k8s-info/k8s_svc_host
echo "$KUBERNETES_SERVICE_PORT" > /var/k8s-info/k8s_svc_port
cat {{.MungeKeyPath}} > /etc/munge/munge.key
service munge start
cp -L -r {{.SlurmConfPath}}/* /etc/slurm
if [ -f "/etc/slurm/slurmdbd.conf" ]; then
    chmod 0600 "/etc/slurm/slurmdbd.conf"
    chown slurm: "/etc/slurm/slurmdbd.conf"
fi
if [ -f "/etc/init.d/prometheus-slurm-exporter" ]; then
    service prometheus-slurm-exporter start
fi
chmod 777 /var/log
slurm_conf="/etc/slurm/slurm.conf"
stateSaveLocation=$(grep "^StateSaveLocation=" "$slurm_conf" | awk -F "=" '{print $2}')
if [ -z "$stateSaveLocation" ]; then
    stateSaveLocation="/var/spool/slurm"
fi
slurmSpoolDir=$(grep "^SlurmdSpoolDir=" "$slurm_conf" | awk -F "=" '{print $2}')
if [ -z "$slurmSpoolDir" ]; then
    slurmSpoolDir="/var/spool/slurm/d"
fi
slurmctldLogFile=$(grep "^SlurmctldLogFile=" "$slurm_conf" | awk -F "=" '{print $2}')
if [ -z "$slurmctldLogFile" ]; then
    slurmctldLogFile="/var/log/slurmctld.log"
fi
slurmdLog=$(grep "^SlurmdLogFile=" "$slurm_conf" | awk -F "=" '{print $2}')
if [ -z "$slurmdLog" ]; then
    slurmdLog="/var/log/slurmd.log"
fi
mkdir -p "$stateSaveLocation"/cluster_state
chown -R slurm:slurm "$stateSaveLocation"
chmod 755 "$stateSaveLocation"
mkdir -p "$slurmSpoolDir"
chmod 755 "$slurmSpoolDir"
chown -R slurm: "$slurmSpoolDir"
touch "$slurmctldLogFile"
chown slurm:slurm "$slurmctldLogFile"
touch "$slurmdLog"
chown slurm: "$slurmdLog"
echo "stateSaveLocation=$stateSaveLocation"
echo "slurmSpoolDir=$slurmSpoolDir"
echo "slurmctldLogFile=$slurmctldLogFile"
echo "slurmdLog=$slurmdLog"
role="master"
if [ -f "/etc/slurm/slurmdbd.conf" ]; then
    service slurmdbd start
    sleep 5
fi
service sshd start
service slurmctld start
```

And the default initBash for worker pods is 
``` bash
cat {{.MungeKeyPath}} > /etc/munge/munge.key
service munge start
cp -L -r {{.SlurmConfPath}}/* /etc/slurm
if [ -f "/etc/slurm/slurmdbd.conf" ]; then
    chmod 0600 "/etc/slurm/slurmdbd.conf"
    chown slurm: "/etc/slurm/slurmdbd.conf"
fi

slurm_conf="/etc/slurm/slurm.conf"
stateSaveLocation=$(grep "^StateSaveLocation=" "$slurm_conf" | awk -F "=" '{print $2}')
if [ -z "$stateSaveLocation" ]; then
    stateSaveLocation="/var/spool/slurm"
fi
slurmSpoolDir=$(grep "^SlurmdSpoolDir=" "$slurm_conf" | awk -F "=" '{print $2}')
if [ -z "$slurmSpoolDir" ]; then
    slurmSpoolDir="/var/spool/slurm/d"
fi
slurmctldLogFile=$(grep "^SlurmctldLogFile=" "$slurm_conf" | awk -F "=" '{print $2}')
if [ -z "$slurmctldLogFile" ]; then
    slurmctldLogFile="/var/log/slurmctld.log"
fi
slurmdLog=$(grep "^SlurmdLogFile=" "$slurm_conf" | awk -F "=" '{print $2}')
if [ -z "$slurmdLog" ]; then
    slurmdLog="/var/log/slurmd.log"
fi
mkdir -p "$stateSaveLocation"/cluster_state
chown -R slurm:slurm "$stateSaveLocation"
chmod 755 "$stateSaveLocation"
mkdir -p "$slurmSpoolDir"
chmod 755 "$slurmSpoolDir"
chown -R slurm: "$slurmSpoolDir"
touch "$slurmctldLogFile"
chown slurm:slurm "$slurmctldLogFile"
touch "$slurmdLog"
chown slurm: "$slurmdLog"
echo "stateSaveLocation=$stateSaveLocation"
echo "slurmSpoolDir=$slurmSpoolDir"
echo "slurmctldLogFile=$slurmctldLogFile"
echo "slurmdLog=$slurmdLog"
CPUs=1
RealMemory=1024
if [ -f "/etc/slurmd-podinfo/cpu_limit" ]; then
    CPUs=$(cat /etc/slurmd-podinfo/cpu_limit)
fi
if [ -f "/etc/slurmd-podinfo/memory_limit" ]; then
    RealMemory=$(cat /etc/slurmd-podinfo/memory_limit)
fi
confOpt=""
nvidia_visible_devices=$NVIDIA_VISIBLE_DEVICES
if [ $# -ge 4 ]; then
    Partition="$4"
    default_option="SLURMD_OPTIONS='-Z -b --conf=\"Feature=$Partition CPUs=$CPUs RealMemory=$RealMemory\"'"
    confOpt="Feature=$Partition CPUs=$CPUs RealMemory=$RealMemory"
else
    default_option="SLURMD_OPTIONS='-Z -b --conf=\"CPUs=$CPUs RealMemory=$RealMemory\"'"
    confOpt="CPUs=$CPUs RealMemory=$RealMemory"
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
        default_option="SLURMD_OPTIONS='-Z --conf=\"Feature=$Partition Gres=$gpu_list CPUs=$CPUs RealMemory=$RealMemory\"'"
        confOpt="Feature=$Partition Gres=$gpu_list CPUs=$CPUs RealMemory=$RealMemory"
    fi
else
    echo "NVIDIA_VISIBLE_DEVICES is empty."
fi
echo $default_option > /etc/default/slurmd
# service slurmd start
slurmd -D -Z -vvv -b --conf "$confOpt"
```

For example, if you set the following command and args in head pod.
```yaml
spec:
  containers:
  - command:
    - bash
    - -c
    args:
    - echo "hello world again again"
```

We will build a head pod like this:
```yaml
spec:
    containers:
    - args:
        - |-
        mkdir /var/k8s-info
        echo "$KUBERNETES_SERVICE_HOST" > /var/k8s-info/k8s_svc_host
        echo "$KUBERNETES_SERVICE_PORT" > /var/k8s-info/k8s_svc_port
        cat /var/munge > /etc/munge/munge.key
        service munge start
        cp -L -r /var/slurm/* /etc/slurm
        if [ -f "/etc/slurm/slurmdbd.conf" ]; then
            chmod 0600 "/etc/slurm/slurmdbd.conf"
            chown slurm: "/etc/slurm/slurmdbd.conf"
        fi
        if [ -f "/etc/init.d/prometheus-slurm-exporter" ]; then
            service prometheus-slurm-exporter start
        fi
        chmod 777 /var/log
        slurm_conf="/etc/slurm/slurm.conf"
        stateSaveLocation=$(grep "^StateSaveLocation=" "$slurm_conf" | awk -F "=" '{print $2}')
        if [ -z "$stateSaveLocation" ]; then
            stateSaveLocation="/var/spool/slurm"
        fi
        slurmSpoolDir=$(grep "^SlurmdSpoolDir=" "$slurm_conf" | awk -F "=" '{print $2}')
        if [ -z "$slurmSpoolDir" ]; then
            slurmSpoolDir="/var/spool/slurm/d"
        fi
        slurmctldLogFile=$(grep "^SlurmctldLogFile=" "$slurm_conf" | awk -F "=" '{print $2}')
        if [ -z "$slurmctldLogFile" ]; then
            slurmctldLogFile="/var/log/slurmctld.log"
        fi
        slurmdLog=$(grep "^SlurmdLogFile=" "$slurm_conf" | awk -F "=" '{print $2}')
        if [ -z "$slurmdLog" ]; then
            slurmdLog="/var/log/slurmd.log"
        fi
        mkdir -p "$stateSaveLocation"/cluster_state
        chown -R slurm:slurm "$stateSaveLocation"
        chmod 755 "$stateSaveLocation"
        mkdir -p "$slurmSpoolDir"
        chmod 755 "$slurmSpoolDir"
        chown -R slurm: "$slurmSpoolDir"
        touch "$slurmctldLogFile"
        chown slurm:slurm "$slurmctldLogFile"
        touch "$slurmdLog"
        chown slurm: "$slurmdLog"
        echo "stateSaveLocation=$stateSaveLocation"
        echo "slurmSpoolDir=$slurmSpoolDir"
        echo "slurmctldLogFile=$slurmctldLogFile"
        echo "slurmdLog=$slurmdLog"
        role="master"
        if [ -f "/etc/slurm/slurmdbd.conf" ]; then
            service slurmdbd start
            sleep 5
        fi
        service sshd start
        service slurmctld start
        "bash" "-c" "echo \"hello world again again\""
        service prometheus-slurm-exporter start
        while true;
        do
            sleep 1000
        done
      command:
        - /bin/bash
        - -lc
        - --
```