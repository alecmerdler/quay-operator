package context

// QuayRegistryContext contains additional information accrued and consumed during the reconcile loop of a `QuayRegistry`.
type QuayRegistryContext struct {
	SupportsRoutes           bool
	ClusterHostname          string
	SupportsObjectStorage    bool
	ObjectStorageInitialized bool
	StorageHostname          string
	StorageBucketName        string
	StorageAccessKey         string
	StorageSecretKey         string
}
