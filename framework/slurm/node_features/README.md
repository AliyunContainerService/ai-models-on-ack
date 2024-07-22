## Why need k8s_resources?

Slurm-Copilot on ack will update GRES on slurm node to prevent over-subscription.

But slurmctld will check GRES when slurmd send a registration rpc request and mark
the node as INVAL if the GRES was updated by slurm-copilot.

So we need to set GRES updated by slurm-copilot as node_features so that slurmctld will
skip the check.

## How to use k8s_resources?

1. Compile k8s_resources plugin

Move the `k8s_resources` directory into `${SLURM_SOURCE_CODE}/src/plugins/node_features`.
Modify the `${SLURM_SOURCE_CODE}/src/plugins/node_features/Makefile` and 
`${SLURM_SOURCE_CODE}/src/plugins/node_features/Makefile.am` to add `k8s_resources` 
in `SUBDIRS`.

Run `make` in `${SLURM_SOURCE_CODE}` to build slurm or only run `make` in  
`${SLURM_SOURCE_CODE}/src/plugins/node_features/k8s_resources` to build the plugin only.

Move `node_features_k8s_resources.so`, `node_features_k8s_resources.la` and 
`node_features_k8s_resources.a` to slurm lib directory. (If you build the whole slurm, `make install`
is enough.)

2. Enable k8s_resources plugin in slurm.conf

Add `NodeFeaturesPlugin=node_features/k8s_resources` in slurm.conf.

## What does k8s_resources do?

k8s_resources will set `gres_ns->node_feature` to `true` when activeFeatures changed on slurm node
so that slurmctld will skip the check.