package api

import (
	"github.com/graphitehealth/operator-sdk/internal/olm/client"
	"github.com/graphitehealth/operator-sdk/internal/olm/installer"
	"k8s.io/client-go/rest"
)

var ErrOLMNotInstalled = client.ErrOLMNotInstalled

func Install(restConfig *rest.Config, version string) error {
	m := installer.Manager{}
	_, err := m.InstallWithRestConfig(restConfig, version)
	if err != nil {
		return err
	}

	return nil
}

func Uninstall(restConfig *rest.Config, version string) error {
	m := installer.Manager{}
	err := m.UninstallWithRestConfig(restConfig, version)
	if err != nil {
		return err
	}

	return nil
}

func GetInstalledVersion(restConfig *rest.Config) (string, error) {
	m := installer.Manager{}
	return m.GetInstalledVersionWithRestConfig(restConfig)
}
