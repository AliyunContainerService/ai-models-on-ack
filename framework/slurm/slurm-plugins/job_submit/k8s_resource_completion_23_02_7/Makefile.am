# Makefile for job_submit/k8s_resource_completion plugin

AUTOMAKE_OPTIONS = foreign

PLUGIN_FLAGS = -module -avoid-version --export-dynamic

AM_CPPFLAGS = -DSLURM_PLUGIN_DEBUG -I$(top_srcdir)

pkglib_LTLIBRARIES = job_submit_k8s_resource_completion.la

# Job submit k8s_resource_completion plugin.
job_submit_k8s_resource_completion_la_SOURCES = job_submit_k8s_resource_completion.c
job_submit_k8s_resource_completion_la_LDFLAGS = $(PLUGIN_FLAGS)
