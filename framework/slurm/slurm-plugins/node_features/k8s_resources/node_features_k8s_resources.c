/*****************************************************************************\
 *  node_features_k8s_resources.c - Plugin for supporting arbitrary node features
 *  using external helper binaries
 *****************************************************************************
 *  Copyright (C) 2021 NVIDIA CORPORATION. All rights reserved.
 *  Written by NVIDIA CORPORATION.
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

#define _GNU_SOURCE

#include <ctype.h>
#include <stdio.h>

#include "slurm/slurm_errno.h"

#include "src/common/list.h"
#include "src/common/job_features.h"
#include "src/common/node_conf.h"
#include "src/common/read_config.h"
#include "src/common/run_command.h"
#include "src/common/uid.h"
#include "src/common/xmalloc.h"
#include "src/common/xstring.h"

#include "src/slurmd/slurmd/slurmd.h"

/*
 * These are defined here so when we link with something other than
 * the slurmctld we will have these symbols defined.  They will get
 * overwritten when linking with the slurmctld.
 */
#if defined (__APPLE__)
extern slurmd_conf_t *conf __attribute__((weak_import));
#else
slurmd_conf_t *conf = NULL;
#endif

const char plugin_name[] = "node_features k8s_resources plugin";
const char plugin_type[] = "node_features/k8s_resources";
const uint32_t plugin_version = SLURM_VERSION_NUMBER;

static uint32_t boot_time = (5 * 60);
static uint32_t exec_time = 10;

static s_p_options_t conf_options[] = {};

extern int init(void)
{
	/* Executed on slurmctld and not used by this plugin */
	return SLURM_SUCCESS;
}

extern int fini(void)
{
	/* Executed on slurmctld and not used by this plugin */
	return SLURM_SUCCESS;
}

extern bool node_features_p_changeable_feature(char *input)
{
	/* Executed on slurmctld and not used by this plugin */
	return true;
}

extern int node_features_p_job_valid(char *job_features, list_t *feature_list)
{
	/* Executed on slurmctld and not used by this plugin */
	return SLURM_SUCCESS;
}

extern int node_features_p_node_set(char *active_features)
{
	/* Executed on slurmctld and not used by this plugin */
	info("original: active_features=%s",active_features);

	if (active_features && !xstrstr(active_features, "k8s")) {	/* Append for multiple node_features plugins */
	    char * avail_sep;
		if (active_features[0])
			avail_sep = ",";
		else
			avail_sep = "";
		xstrfmtcat(active_features, "%sk8s", avail_sep);
	} else {
		active_features = "k8s";
	}

	info("new: active_features=%s",active_features);
	return SLURM_SUCCESS;
}

extern void node_features_p_node_state(char **avail_modes, char **current_mode)
{
	info("original: avail=%s current=%s",*avail_modes, *current_mode);

	if (*avail_modes) {	/* Append for multiple node_features plugins */
	    char * avail_sep;
		if (*avail_modes[0])
			avail_sep = ",";
		else
			avail_sep = "";
		xstrfmtcat(*avail_modes, "%sk8s", avail_sep);
	} else {
		*avail_modes = "k8s";
	}

	if (*current_mode) {	/* Append for multiple node_features plugins */
	    char * cur_sep;
		if (*current_mode[0])
			cur_sep = ",";
		else
			cur_sep = "";
		xstrfmtcat(*current_mode, "%sk8s", cur_sep);
	} else {
		*current_mode = "k8s";
	}

	info("new: avail=%s current=%s",*avail_modes, *current_mode);
}

extern char *node_features_p_node_xlate(char *new_features, char *orig_features,
					char *avail_features, int node_inx)
{
	info("node_features_p_node_xlate is called");
	/* Executed on slurmctld and not used by this plugin */
	char *ret_features = xstrdup(orig_features);
	node_record_t *node_ptr;
	bitstr_t *node_bitmap = node_conf_get_active_bitmap();

	node_ptr = next_node_bitmap(node_bitmap, &node_inx);
	if (!node_ptr){
		info("node_ptr is NULL: %d", node_inx);
		return NULL;
	}
	if (!node_ptr->gres){
		node_ptr->gres =
			xstrdup(node_ptr->config_ptr->gres);
	}
	gres_node_feature(node_ptr->name, "k8scpu",
				node_ptr->cpus > 0 ? node_ptr->cpus : 0, &node_ptr->gres,
				&node_ptr->gres_list);
	gres_node_feature(node_ptr->name, "k8smemory",
		node_ptr->real_memory > 0 ? node_ptr->real_memory : 0, &node_ptr->gres,
				&node_ptr->gres_list);

	info("set k8scpu and k8smemory features for node %d, node_ptr->cpus %d, node_ptr->real_memory %d. set node features, \
original: orig_features=%s, avail_features=%s, new_features=%s", 
		node_inx, node_ptr->cpus, node_ptr->real_memory, orig_features, avail_features, new_features);
	/* Executed on slurmctld and not used by this plugin */
	if (!xstrstr(orig_features, "k8s") && !xstrcmp(new_features, "k8s") ) {	/* Append for multiple node_features plugins */
		error("new: orig_features=%s",orig_features);
	    char * avail_sep;
		if (orig_features && orig_features[0])
			avail_sep = ",";
		else
			avail_sep = "";
		xstrfmtcat(ret_features, "%sk8s", avail_sep);
	} 
	info("new: ret_features=%s",ret_features);
	return ret_features;
}

extern char *node_features_p_job_xlate(char *job_features,
				       list_t *feature_list,
				       bitstr_t *job_node_bitmap)
{
	/* Executed on slurmctld and not used by this plugin */
	return NULL;
}

/* Return true if the plugin requires PowerSave mode for booting nodes */
extern bool node_features_p_node_power(void)
{
	/* Executed on slurmctld and not used by this plugin */
	return false;
}

/* Get node features plugin configuration */
extern void node_features_p_get_config(config_plugin_params_t *p)
{
	/* Executed on slurmctld and not used by this plugin */
}

extern bitstr_t *node_features_p_get_node_bitmap(void)
{
	/* Executed on slurmctld and not used by this plugin */
	return node_conf_get_active_bitmap();
}

extern char *node_features_p_node_xlate2(char *new_features)
{
	/* Executed on slurmctld and not used by this plugin */
	return xstrdup(new_features);
}

extern uint32_t node_features_p_boot_time(void)
{
	/* Executed on slurmctld and not used by this plugin */
	return boot_time;
}

extern int node_features_p_reconfig(void)
{
	/* Executed on slurmctld and not used by this plugin */
	return SLURM_SUCCESS;
}

extern bool node_features_p_user_update(uid_t uid)
{
	/* Executed on slurmctld and not used by this plugin */
	return true;
}

extern void node_features_p_step_config(bool mem_sort, bitstr_t *numa_bitmap)
{
	/* Executed on slurmctld and not used by this plugin */
	return;
}

extern int node_features_p_overlap(bitstr_t *active_bitmap)
{
	/* Executed on slurmctld and not used by this plugin */
	return bit_set_count(active_bitmap);
}

extern int node_features_p_get_node(char *node_list)
{
	/* Executed on slurmctld and not used by this plugin */
	return SLURM_SUCCESS;
}

/*
 * Note the active features associated with a set of nodes have been updated.
 * Specifically update the node's "hbm" GRES and "CpuBind" values as needed.
 * IN active_features - New active features
 * IN node_bitmap - bitmap of nodes changed
 * RET error code
 */
extern int node_features_p_node_update(char *active_features,
				       bitstr_t *node_bitmap)
{
	/* Executed on slurmctld and not used by this plugin */
	info("node_features_p_node_update is called");
	int i;
	node_record_t *node_ptr;
	for (i = 0; (node_ptr = next_node_bitmap(node_bitmap, &i)); i++) {
			if (!node_ptr->gres)
				node_ptr->gres =
					xstrdup(node_ptr->config_ptr->gres);
			gres_node_feature(node_ptr->name, "k8scpu",
					  node_ptr->cpus > 0 ? node_ptr->cpus : 0, &node_ptr->gres,
					  &node_ptr->gres_list);
			gres_node_feature(node_ptr->name, "k8smemory",
					  node_ptr->real_memory > 0 ? node_ptr->real_memory : 0, &node_ptr->gres,
					  &node_ptr->gres_list);

			info("set k8scpu and k8smemory features for node, node_ptr->cpus %d, node_ptr->real_memory %d.", 
				node_ptr->cpus, node_ptr->real_memory);
	}
	return SLURM_SUCCESS;
}

extern bool node_features_p_node_update_valid(void *node_ptr,
					      update_node_msg_t *update_node_msg)
{
	/* Executed on slurmctld and not used by this plugin */
	return true;
}
