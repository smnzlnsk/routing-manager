package repository

type Repositories struct {
	// The alert repository is used to store the triggered routing alerts of the service instances.
	AlertRepository AlertRepository

	// The interest repository is used to store the currently active interests of the service instances.
	// TODO: This repository can be moved to jobs/jobs database table, as it does not contain any new information.
	// The interests can instead be taken from the interested_nodes field of each service.
	InterestRepository InterestRepository

	// The job repository is used to store the jobs of the service instances.
	// It is generally managed by the cluster's service-manager.
	// The routing-manager only reads the data through this repository.
	JobRepository JobRepository

	// The routing repository bases its information on the job repository.
	// It is generally managed by the cluster's service-manager.
	// The routing-manager only administers the routing priorities for each service instance through this repository.
	RoutingRepository RoutingRepository

	// TODO: Add other repositories here
}
