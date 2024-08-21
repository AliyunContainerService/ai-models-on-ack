/*****************************************************************************\
 *  job_submit_k8s_resource_completion.c - Set k8s_resource_completion in job submit request specifications.
 *****************************************************************************
 *  Copyright (C) 2010 Lawrence Livermore National Security.
 *  Produced at Lawrence Livermore National Laboratory (cf, DISCLAIMER).
 *  Written by Morris Jette <jette1@llnl.gov>
 *  CODE-OCEC-09-009. All rights reserved.
 *
 *  This file is part of Slurm, a resource management program.
 *  For details, see <https://slurm.schedmd.com/>.
 *  Please also read the included file: DISCLAIMER.
 *
 *  Slurm is free software; you can redistribute it and/or modify it under
 *  the terms of the GNU General Public License as published by the Free
 *  Software Foundation; either version 2 of the License, or (at your option)
 *  any later version.
 *
 *  In addition, as a special exception, the copyright holders give permission
 *  to link the code of portions of this program with the OpenSSL library under
 *  certain conditions as described in each individual source file, and
 *  distribute linked combinations including the two. You must obey the GNU
 *  General Public License in all respects for all of the code used other than
 *  OpenSSL. If you modify file(s) with this exception, you may extend this
 *  exception to your version of the file(s), but you are not obligated to do
 *  so. If you do not wish to do so, delete this exception statement from your
 *  version.  If you delete this exception statement from all source files in
 *  the program, then also delete it here.
 *
 *  Slurm is distributed in the hope that it will be useful, but WITHOUT ANY
 *  WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
 *  FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more
 *  details.
 *
 *  You should have received a copy of the GNU General Public License along
 *  with Slurm; if not, write to the Free Software Foundation, Inc.,
 *  51 Franklin Street, Fifth Floor, Boston, MA 02110-1301  USA.
\*****************************************************************************/

#include <inttypes.h>
#include <stdio.h>
#include <string.h>
#include <sys/types.h>
#include <unistd.h>

#include "slurm/slurm_errno.h"
#include "src/common/slurm_xlator.h"
#include "src/slurmctld/slurmctld.h"
#include "src/common/log.h"
#include "src/common/xstring.h"

#define MAX_ACCTG_FREQUENCY 30

/*
 * These variables are required by the generic plugin interface.  If they
 * are not found in the plugin, the plugin loader will ignore it.
 *
 * plugin_name - a string giving a human-readable description of the
 * plugin.  There is no maximum length, but the symbol must refer to
 * a valid string.
 *
 * plugin_type - a string suggesting the type of the plugin or its
 * applicability to a particular form of data or method of data handling.
 * If the low-level plugin API is used, the contents of this string are
 * unimportant and may be anything.  Slurm uses the higher-level plugin
 * interface which requires this string to be of the form
 *
 *	<application>/<method>
 *
 * where <application> is a description of the intended application of
 * the plugin (e.g., "auth" for Slurm authentication) and <method> is a
 * description of how this plugin satisfies that application.  Slurm will
 * only load authentication plugins if the plugin_type string has a prefix
 * of "auth/".
 *
 * plugin_version - an unsigned 32-bit integer containing the Slurm version
 * (major.minor.micro combined into a single number).
 */
const char plugin_name[]       	= "Job submit k8s_resource_completion plugin";
const char plugin_type[]       	= "job_submit/k8s_resource_completion";
const uint32_t plugin_version   = SLURM_VERSION_NUMBER;

/*****************************************************************************\
 * We've provided a simple example of the type of things you can do with this
 * plugin. If you develop another plugin that may be of interest to others
 * please post it to slurm-dev@schedmd.com  Thanks!
\*****************************************************************************/
uint64_t get_total_gpu_count(char * IN_VAL) {
	char *tres_type = NULL, *name = NULL, *type = NULL, *save_ptr = NULL;
	int rc = true;
	// total cnt is the sum of all gpu types.
	uint64_t cnt = 0, totalCnt = 0;
	// in 23.02.7, the first parameter of slurm_get_next_tres is char*
	// in 23.11, the first parameter of slurm_get_next_tres is char**
	while (((rc = slurm_get_next_tres(tres_type,
				IN_VAL,
				&name, &type,
				&cnt, &save_ptr)) == SLURM_SUCCESS) &&
				save_ptr) {
		if (xstrcmp(tres_type, "gres")) {
			continue;
		}
		if (xstrcmp(name, "gpu")) {
			continue;
		}
		totalCnt += cnt;
	}
	return totalCnt;
}

int _job_set(job_desc_msg_t *job_desc, char **err_msg) {
	// --cpus-per-task is not compatible with --cpus-per-gpu
	if (job_desc->cpus_per_task) {
		if (job_desc->tres_per_task) {
			xstrfmtcat(job_desc->tres_per_task, ",gres:k8scpu:%d", job_desc->cpus_per_task);
		} else {
			job_desc->tres_per_task = xstrdup_printf("gres:k8scpu:%d", job_desc->cpus_per_task);
		}
		info("set job's tres_per_task as gres:k8scpu:%d", job_desc->cpus_per_task);
	}
	if (job_desc->cpus_per_tres) {
		char *tres_type = "gres", *name = NULL, *type = NULL, *save_ptr = NULL;
		int rc = true;
		uint64_t cpus_per_tres = 0;
		// now only one term will be in cpus_per_tres, so we only need to run once.
		// in 23.02.7, the first parameter of slurm_get_next_tres is char*
		// in 23.11, the first parameter of slurm_get_next_tres is char**
		rc = slurm_get_next_tres(tres_type, job_desc->tres_per_node, &name, &type, &cpus_per_tres, &save_ptr);
		if (rc != SLURM_SUCCESS) {
			info("job %d request cpus_per_tres %s", job_desc->job_id, job_desc->cpus_per_tres);
			if (err_msg) {
				*err_msg = xstrdup("failed to run slurm_get_next_tres to get cnt for cpus_per_tres when try to modify job");
			}
			return SLURM_ERROR;
		}
		if (job_desc->tres_per_node) {
			uint64_t totalCnt = get_total_gpu_count(job_desc->tres_per_node);
			xstrfmtcat(job_desc->tres_per_node, ",gres:k8scpu:%lu", totalCnt * cpus_per_tres);
			info("set job's tres_per_node as gres:k8scpu:%lu", totalCnt * cpus_per_tres);
		} else if (job_desc->tres_per_task) {
			uint64_t totalCnt = get_total_gpu_count(job_desc->tres_per_task);
			xstrfmtcat(job_desc->tres_per_task, ",gres:k8scpu:%lu", totalCnt * cpus_per_tres);
			info("set job's tres_per_task as gres:k8scpu:%lu", totalCnt * cpus_per_tres);
		} else {
			if (err_msg) {
				*err_msg = xstrdup("cpus_per_tres should be used with tres_per_node or tres_per_task");
			}
			return SLURM_ERROR;
		}
	}
	// --mem is not compatible with --mem-per-gpu
	if (job_desc->pn_min_memory) {
		if (job_desc->tres_per_node) {
			xstrfmtcat(job_desc->tres_per_node, ",gres:k8smemory:%lu", job_desc->pn_min_memory);
		} else {
			job_desc->tres_per_node = xstrdup_printf("gres:k8smemory:%lu", job_desc->pn_min_memory);
		}
		info("set job's tres_per_node as gres:k8smemory:%lu", job_desc->pn_min_memory);
	}
	if (job_desc->mem_per_tres) {
		char *tres_type = "gres", *name = NULL, *type = NULL, *save_ptr = NULL;
		int rc = true;
		uint64_t mem_per_tres = 0;
		// now only one term will be in mem_per_tres, so we only need to run once.
		// in 23.02.7, the first parameter of slurm_get_next_tres is char*
		// in 23.11, the first parameter of slurm_get_next_tres is char**
		rc = slurm_get_next_tres(tres_type, job_desc->tres_per_node, &name, &type, &mem_per_tres, &save_ptr);
		if (rc != SLURM_SUCCESS) {
			info("job %d request mem_per_tres %s", job_desc->job_id, job_desc->mem_per_tres);
			if (err_msg) {
				*err_msg = xstrdup("failed to run slurm_get_next_tres to get cnt for mem_per_tres when try to modify job");
			}
			return SLURM_ERROR;
		}
		if (job_desc->tres_per_node) {
			uint64_t totalCnt = get_total_gpu_count(job_desc->tres_per_node);
			xstrfmtcat(job_desc->tres_per_node, ",gres:k8smemory:%lu", totalCnt * mem_per_tres);
			info("set job's tres_per_node as gres:k8smemory:%lu", totalCnt * mem_per_tres);
		} else if (job_desc->tres_per_task) {
			uint64_t totalCnt = get_total_gpu_count(job_desc->tres_per_task);
			xstrfmtcat(job_desc->tres_per_task, ",gres:k8smemory:%lu", totalCnt * mem_per_tres);
			info("set job's tres_per_task as gres:k8smemory:%lu", totalCnt * mem_per_tres);
		} else {
			if (err_msg) {
				*err_msg = xstrdup("mem_per_tres should be used with tres_per_node or tres_per_task");
			}
			return SLURM_ERROR;
		}
	}

	info("job %u submit success", job_desc->job_id);
	return SLURM_SUCCESS;
}

extern int job_submit(job_desc_msg_t *job_desc, uint32_t submit_uid,
		      char **err_msg)
{
#if 0
	uint16_t acctg_freq = 0;
	if (job_desc->acctg_freq)
		acctg_freq = atoi(job_desc->acctg_freq);
	/* This example code will prevent users from setting an accounting
	 * frequency of less than 30 seconds in order to ensure more precise
	 *  accounting. Also remove any QOS value set by the user in order
	 * to use the default value from the database. */
	if (acctg_freq < MIN_ACCTG_FREQUENCY) {
		info("Changing accounting frequency of submitted job "
		     "from %u to %u",
		     acctg_freq, MIN_ACCTG_FREQUENCY);
		job_desc->acctg_freq = xstrdup_printf(
			"%d", MIN_ACCTG_FREQUENCY);
		if (err_msg)
			*err_msg = xstrdup("Changed job frequency");
	}

	if (job_desc->qos) {
		info("Clearing QOS (%s) from submitted job", job_desc->qos);
		xfree(job_desc->qos);
	}
#endif
    if (!(job_desc->bitflags & JOB_NTASKS_SET)) {
		info("job %d do not set --ntasks", job_desc->job_id);
		if (err_msg) {
			*err_msg = xstrdup("job must set --ntasks");
		}
		return SLURM_ERROR;
	}
	if (job_desc->tres_per_job) {
		info("job %d request tres_per_job %s", job_desc->job_id, job_desc->tres_per_job);
		if (err_msg) {
			*err_msg = xstrdup("tres_per_job is not supported for now, do not use --gpus, use \
				cpus_per_task with gpus_per_task orcpus_per_gpu with gpus_per_node instead.");
		}
		return SLURM_ERROR;
	}
	if (job_desc->tres_per_socket) {
		info("job %d request tres_per_socket %s", job_desc->job_id, job_desc->tres_per_socket);
		if (err_msg) {
			*err_msg = xstrdup("tres_per_socket is not supported for now, do not use --gpus-per-socket, use \
				cpus_per_task with gpus_per_task orcpus_per_gpu with gpus_per_node instead.");
		}
		return SLURM_ERROR;
	}
	return _job_set(job_desc, err_msg);
}

extern int job_modify(job_desc_msg_t *job_desc, job_record_t *job_ptr,
		      uint32_t submit_uid, char **err_msg)
{
#if 0
	uint16_t acctg_freq = 0;
	if (job_desc->acctg_freq)
		acctg_freq = atoi(job_desc->acctg_freq);
	/* This example code will prevent users from setting an accounting
	 * frequency of less than 30 seconds in order to ensure more precise
	 *  accounting. Also remove any QOS value set by the user in order
	 * to use the default value from the database. */
	if (acctg_freq < MIN_ACCTG_FREQUENCY) {
		info("Changing accounting frequency of modify job %u "
		     "from %u to %u", job_ptr->job_id,
		     job_desc->acctg_freq, MIN_ACCTG_FREQUENCY);
		job_desc->acctg_freq = xstrdup_printf(
			"%d", MIN_ACCTG_FREQUENCY);
	}

	if (job_desc->qos) {
		info("Clearing QOS (%s) from modify of job %u",
		     job_desc->qos, job_ptr->job_id);
		xfree(job_desc->qos);
	}
#endif
	return _job_set(job_desc, err_msg);
}
