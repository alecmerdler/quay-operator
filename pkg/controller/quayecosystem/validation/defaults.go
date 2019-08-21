package validation

import (
	redhatcopv1alpha1 "github.com/redhat-cop/quay-operator/pkg/apis/redhatcop/v1alpha1"
	"github.com/redhat-cop/quay-operator/pkg/controller/quayecosystem/constants"
	"github.com/redhat-cop/quay-operator/pkg/controller/quayecosystem/resources"
	"github.com/redhat-cop/quay-operator/pkg/controller/quayecosystem/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func SetDefaults(client client.Client, quayConfiguration *resources.QuayConfiguration) bool {

	changed := false

	// Initialize Base variables and objects
	quayConfiguration.QuayConfigUsername = constants.QuayConfigUsername
	quayConfiguration.QuayConfigPassword = constants.QuayConfigDefaultPasswordValue
	quayConfiguration.QuaySuperuserUsername = constants.QuaySuperuserDefaultUsername
	quayConfiguration.QuaySuperuserPassword = constants.QuaySuperuserDefaultPassword
	quayConfiguration.QuaySuperuserEmail = constants.QuaySuperuserDefaultEmail
	quayConfiguration.QuayConfigPasswordSecret = resources.GetQuayConfigResourcesName(quayConfiguration.QuayEcosystem)
	quayConfiguration.QuayDatabase.Username = constants.QuayDatabaseCredentialsDefaultUsername
	quayConfiguration.QuayDatabase.Password = constants.QuayDatabaseCredentialsDefaultPassword
	quayConfiguration.QuayDatabase.Database = constants.QuayDatabaseCredentialsDefaultDatabaseName
	quayConfiguration.QuayDatabase.RootPassword = constants.QuayDatabaseCredentialsDefaultRootPassword
	quayConfiguration.QuayDatabase.Server = resources.GetDatabaseResourceName(quayConfiguration.QuayEcosystem, constants.DatabaseComponentQuay)
	quayConfiguration.ClairDatabase.Username = constants.ClairDatabaseCredentialsDefaultUsername
	quayConfiguration.ClairDatabase.Password = constants.ClairDatabaseCredentialsDefaultPassword
	quayConfiguration.ClairDatabase.Server = resources.GetDatabaseResourceName(quayConfiguration.QuayEcosystem, constants.DatabaseComponentClair)
	quayConfiguration.ClairDatabase.Database = constants.ClairDatabaseCredentialsDefaultDatabaseName
	quayConfiguration.ClairDatabase.RootPassword = constants.ClairDatabaseCredentialsDefaultRootPassword
	quayConfiguration.ClairUpdateInterval = constants.ClairDefaultUpdateInterval

	if quayConfiguration.QuayEcosystem.Spec.Quay == nil {
		quayConfiguration.QuayEcosystem.Spec.Quay = &redhatcopv1alpha1.Quay{}
		changed = true
	}

	if quayConfiguration.QuayEcosystem.Spec.Quay.Database == nil {
		quayConfiguration.QuayEcosystem.Spec.Quay.Database = &redhatcopv1alpha1.Database{}
		changed = true
	}

	if quayConfiguration.QuayEcosystem.Spec.Redis == nil {
		quayConfiguration.QuayEcosystem.Spec.Redis = &redhatcopv1alpha1.Redis{}
		changed = true
	}

	// Core Quay Values
	if utils.IsZeroOfUnderlyingType(quayConfiguration.QuayEcosystem.Spec.Quay.Image) {
		changed = true
		quayConfiguration.QuayEcosystem.Spec.Quay.Image = constants.QuayImage
	}

	// Default Quay Config Route
	if utils.IsZeroOfUnderlyingType(quayConfiguration.QuayEcosystem.Spec.Quay.ConfigRouteHost) {
		quayConfiguration.QuayConfigHostname = resources.GetQuayConfigResourcesName(quayConfiguration.QuayEcosystem)
	} else {
		quayConfiguration.QuayConfigHostname = quayConfiguration.QuayEcosystem.Spec.Quay.ConfigRouteHost
	}

	// Apply default values for Redis
	if utils.IsZeroOfUnderlyingType(quayConfiguration.QuayEcosystem.Spec.Redis.Hostname) {

		if utils.IsZeroOfUnderlyingType(quayConfiguration.QuayEcosystem.Spec.Redis.Image) {
			changed = true
			quayConfiguration.QuayEcosystem.Spec.Redis.Image = constants.RedisImage
		}
	}

	// Set Redis Hostname
	if utils.IsZeroOfUnderlyingType(quayConfiguration.QuayEcosystem.Spec.Redis.Hostname) {
		quayConfiguration.RedisHostname = resources.GetRedisResourcesName(quayConfiguration.QuayEcosystem)
	} else {
		quayConfiguration.RedisHostname = quayConfiguration.QuayEcosystem.Spec.Redis.Hostname
	}

	// Set Redis Port
	if utils.IsZeroOfUnderlyingType(quayConfiguration.QuayEcosystem.Spec.Redis.Port) {
		quayConfiguration.RedisPort = &constants.RedisPort
	} else {
		quayConfiguration.RedisPort = quayConfiguration.QuayEcosystem.Spec.Redis.Port
	}

	// User would like to have a database automatically provisioned if server not provided
	if utils.IsZeroOfUnderlyingType(quayConfiguration.QuayEcosystem.Spec.Quay.Database.Server) {

		quayConfiguration.QuayDatabase.Server = resources.GetDatabaseResourceName(quayConfiguration.QuayEcosystem, constants.DatabaseComponentQuay)

		if utils.IsZeroOfUnderlyingType(quayConfiguration.QuayEcosystem.Spec.Quay.Database.Image) {
			changed = true
			quayConfiguration.QuayEcosystem.Spec.Quay.Database.Image = constants.PostgresqlImage
		}

	} else {
		quayConfiguration.QuayDatabase.Server = quayConfiguration.QuayEcosystem.Spec.Quay.Database.Server

	}

	// Clair Core Values
	if quayConfiguration.QuayEcosystem.Spec.Clair != nil && quayConfiguration.QuayEcosystem.Spec.Clair.Enabled {

		// Add Clair Service Account to List of SCC's
		constants.RequiredAnyUIDSccServiceAccounts = append(constants.RequiredAnyUIDSccServiceAccounts, constants.ClairServiceAccount)

		if utils.IsZeroOfUnderlyingType(quayConfiguration.QuayEcosystem.Spec.Clair.Image) {
			changed = true
			quayConfiguration.QuayEcosystem.Spec.Clair.Image = constants.ClairImage
		}
		if quayConfiguration.QuayEcosystem.Spec.Clair.Database == nil {
			quayConfiguration.QuayEcosystem.Spec.Clair.Database = &redhatcopv1alpha1.Database{}
			changed = true
		}

		// User would like to have a database automatically provisioned if server not provided
		if utils.IsZeroOfUnderlyingType(quayConfiguration.QuayEcosystem.Spec.Clair.Database.Server) {

			quayConfiguration.ClairDatabase.Server = resources.GetDatabaseResourceName(quayConfiguration.QuayEcosystem, constants.DatabaseComponentClair)

			if utils.IsZeroOfUnderlyingType(quayConfiguration.QuayEcosystem.Spec.Clair.Database.Image) {
				changed = true
				quayConfiguration.QuayEcosystem.Spec.Clair.Database.Image = constants.PostgresqlImage
			}

		} else {
			quayConfiguration.ClairDatabase.Server = quayConfiguration.QuayEcosystem.Spec.Clair.Database.Server
		}

	}

	if !utils.IsZeroOfUnderlyingType(quayConfiguration.QuayEcosystem.Spec.Quay.KeepConfigDeployment) && quayConfiguration.QuayEcosystem.Spec.Quay.KeepConfigDeployment {
		quayConfiguration.DeployQuayConfiguration = true
	}

	if !quayConfiguration.QuayEcosystem.Status.SetupComplete {
		quayConfiguration.DeployQuayConfiguration = true
	}

	if !utils.IsZeroOfUnderlyingType(quayConfiguration.QuayEcosystem.Spec.Quay.RegistryStorage) {

		if utils.IsZeroOfUnderlyingType(quayConfiguration.QuayEcosystem.Spec.Quay.RegistryStorage.PersistentVolumeAccessModes) {
			quayConfiguration.QuayEcosystem.Spec.Quay.RegistryStorage.PersistentVolumeAccessModes = constants.QuayRegistryStoragePersistentVolumeAccessModes
			changed = true
		}

		if utils.IsZeroOfUnderlyingType(quayConfiguration.QuayEcosystem.Spec.Quay.RegistryStorage.PersistentVolumeSize) {
			quayConfiguration.QuayEcosystem.Spec.Quay.RegistryStorage.PersistentVolumeSize = constants.QuayRegistryStoragePersistentVolumeStoreSize
			changed = true
		}
	}

	if utils.IsZeroOfUnderlyingType(quayConfiguration.QuayEcosystem.Spec.Quay.RegistryBackends) {
		// Generate Default Local Storage
		quayConfiguration.QuayEcosystem.Spec.Quay.RegistryBackends = []redhatcopv1alpha1.RegistryBackend{
			redhatcopv1alpha1.RegistryBackend{
				Name: constants.RegistryStorageDefaultName,
				RegistryBackendSource: redhatcopv1alpha1.RegistryBackendSource{
					Local: &redhatcopv1alpha1.LocalRegistryBackendSource{
						StoragePath: constants.QuayRegistryStoragePath,
					},
				},
			},
		}

		changed = true
	}

	return changed
}