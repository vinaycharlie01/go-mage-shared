package helmmagex

import (
	"github.com/vinaycharlie01/go-mage-shared/helmx"
)

// Package-level convenience functions for backward compatibility
var defaultRunner = helmx.NewHelmRunner()

// Install installs a Helm chart
func Install(opts helmx.InstallOptions) error {
	return defaultRunner.Install(opts)
}

// Upgrade upgrades a Helm release
func Upgrade(opts helmx.UpgradeOptions) error {
	return defaultRunner.Upgrade(opts)
}

// Uninstall uninstalls a Helm release
func Uninstall(releaseName, namespace string, args ...string) error {
	return defaultRunner.Uninstall(releaseName, namespace, args...)
}

// List lists Helm releases
func List(namespace string, args ...string) error {
	return defaultRunner.List(namespace, args...)
}

// Status shows the status of a Helm release
func Status(releaseName, namespace string, args ...string) error {
	return defaultRunner.Status(releaseName, namespace, args...)
}

// Template renders chart templates locally
func Template(releaseName, chart string, args ...string) error {
	return defaultRunner.Template(releaseName, chart, args...)
}

// Lint runs helm lint on a chart
func Lint(chart string, args ...string) error {
	return defaultRunner.Lint(chart, args...)
}

// Package packages a chart directory into a chart archive
func Package(chart string, args ...string) error {
	return defaultRunner.Package(chart, args...)
}

// RepoAdd adds a chart repository
func RepoAdd(name, url string, args ...string) error {
	return defaultRunner.RepoAdd(name, url, args...)
}

// RepoUpdate updates chart repositories
func RepoUpdate(args ...string) error {
	return defaultRunner.RepoUpdate(args...)
}
