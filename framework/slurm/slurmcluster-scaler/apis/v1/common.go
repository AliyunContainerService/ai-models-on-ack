package v1

// JobStatus is the Slurm Job Status.
type JobStatus string

const (
	JobStatusPending   JobStatus = "PENDING"
	JobStatusRunning   JobStatus = "RUNNING"
	JobStatusStopped   JobStatus = "STOPPED"
	JobStatusSucceeded JobStatus = "SUCCEEDED"
	JobStatusFailed    JobStatus = "FAILED"
)

// JobDeploymentStatus indicates SlurmJob status including SlurmCluster lifecycle management and Job submission
type JobDeploymentStatus string

const (
	JobDeploymentStatusInitializing                    JobDeploymentStatus = "Initializing"
	JobDeploymentStatusFailedToGetOrCreateSlurmCluster JobDeploymentStatus = "JobDeploymentStatusFailedToGetOrCreateSlurmCluster"
	JobDeploymentStatusWaitForDashboard                JobDeploymentStatus = "WaitForDashboard"
	JobDeploymentStatusWaitForDashboardReady           JobDeploymentStatus = "WaitForDashboardReady"
	JobDeploymentStatusWaitForK8sJob                   JobDeploymentStatus = "WaitForK8sJob"
	JobDeploymentStatusFailedJobDeploy                 JobDeploymentStatus = "FailedJobDeploy"
	JobDeploymentStatusRunning                         JobDeploymentStatus = "Running"
	JobDeploymentStatusFailedToGetJobStatus            JobDeploymentStatus = "FailedToGetJobStatus"
	JobDeploymentStatusComplete                        JobDeploymentStatus = "Complete"
	JobDeploymentStatusSuspended                       JobDeploymentStatus = "Suspended"
)

type ClusterNodeType string

const (
	HeadNode   ClusterNodeType = "head"
	WorkerNode ClusterNodeType = "worker"
)

type EventReason string

const (
	SlurmConfigError       EventReason = "SlurmConfigError"
	PodReconciliationError EventReason = "PodReconciliationError"
)
