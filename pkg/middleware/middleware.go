package middleware

import (
	"strings"

	route "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	v1 "github.com/quay/quay-operator/apis/quay/v1"
	quaycontext "github.com/quay/quay-operator/pkg/context"
)

const configSecretPrefix = "quay-config-secret"

// Process applies any additional middleware steps to a managed k8s object that cannot be accomplished using
// the Kustomize toolchain.
func Process(ctx *quaycontext.QuayRegistryContext, quay *v1.QuayRegistry, obj client.Object) (client.Object, error) {
	objectMeta, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}

	quayComponentLabel := labels.Set(objectMeta.GetLabels()).Get("quay-component")

	// Flatten config bundle `Secret`
	if strings.Contains(objectMeta.GetName(), configSecretPrefix+"-") {
		configBundleSecret, err := FlattenSecret(obj.(*corev1.Secret))
		if err != nil {
			return nil, err
		}

		// NOTE: Remove user-provided TLS cert/key pair so it is not mounted to `/conf/stack`, otherwise it will affect Quay's generated NGINX config
		delete(configBundleSecret.Data, "ssl.cert")
		delete(configBundleSecret.Data, "ssl.key")

		return configBundleSecret, nil
	}

	// Add TLS cert/key pair to managed `Route`
	if routeObj, ok := obj.(*route.Route); ok &&
		!v1.ComponentIsManaged(quay.Spec.Components, v1.ComponentTLS) &&
		(quayComponentLabel == "quay-app-route" || quayComponentLabel == "quay-builder-route") {
		routeObj.Spec.TLS = &route.TLSConfig{
			Certificate: string(ctx.TLSCert),
			Key:         string(ctx.TLSKey),
		}

		return routeObj, nil
	}

	return obj, nil
}

// flattenSecret takes all Quay config fields in given secret and combines them under `config.yaml` key.
func FlattenSecret(configBundle *corev1.Secret) (*corev1.Secret, error) {
	flattenedSecret := configBundle.DeepCopy()

	var flattenedConfig map[string]interface{}
	if err := yaml.Unmarshal(configBundle.Data["config.yaml"], &flattenedConfig); err != nil {
		return nil, err
	}

	isConfigField := func(field string) bool {
		return strings.Contains(field, ".config.yaml")
	}

	for key, file := range configBundle.Data {
		if isConfigField(key) {
			var valueYAML map[string]interface{}
			if err := yaml.Unmarshal(file, &valueYAML); err != nil {
				return nil, err
			}

			for configKey, configValue := range valueYAML {
				flattenedConfig[configKey] = configValue
			}
			delete(flattenedSecret.Data, key)
		}
	}

	flattenedConfigYAML, err := yaml.Marshal(flattenedConfig)
	if err != nil {
		return nil, err
	}

	flattenedSecret.Data["config.yaml"] = []byte(flattenedConfigYAML)

	return flattenedSecret, nil
}
