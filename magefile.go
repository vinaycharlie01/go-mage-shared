//go:build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/vinaycharlie01/go-mage-shared/helmmagex"
	"github.com/vinaycharlie01/go-mage-shared/helmx"
	"github.com/vinaycharlie01/go-mage-shared/komagex"
	"github.com/vinaycharlie01/go-mage-shared/kox"
)

// Helm namespace for Helm-related targets
type Helm mg.Namespace

// Install installs a Helm chart
func (Helm) Install() error {
	return helmmagex.Install(helmx.InstallOptions{
		ReleaseName:     "example",
		Chart:           "./charts/example",
		Namespace:       "default",
		CreateNamespace: true,
		Wait:            true,
	})
}

// Upgrade upgrades a Helm release
func (Helm) Upgrade() error {
	return helmmagex.Upgrade(helmx.UpgradeOptions{
		ReleaseName: "example",
		Chart:       "./charts/example",
		Namespace:   "default",
		Install:     true,
		Wait:        true,
	})
}

// Uninstall uninstalls a Helm release
func (Helm) Uninstall() error {
	return helmmagex.Uninstall("example", "default")
}

// List lists all Helm releases
func (Helm) List() error {
	return helmmagex.List("", "--all-namespaces")
}

// Lint lints a Helm chart
func (Helm) Lint() error {
	return helmmagex.Lint("./charts/example")
}

// RepoUpdate updates Helm repositories
func (Helm) RepoUpdate() error {
	return helmmagex.RepoUpdate()
}

// Ko namespace for Ko (container building) targets
type Ko mg.Namespace

// Build builds a container image with ko
func (Ko) Build() error {
	return komagex.Build(kox.BuildOptions{
		ImportPath: "/Users/vinaykumar/selfhosted/enlearn/operator-1/dist/darwin_arm64/gateway-controller-linux-amd64",
		Tags:       []string{"latest"},
		Platform:   []string{"linux/amd64"},
		Local:      true,
	})
}

// BuildMultiPlatform builds multi-platform container images
func (Ko) BuildMultiPlatform() error {
	return komagex.Build(kox.BuildOptions{
		ImportPath: "./cmd/app",
		Tags:       []string{"latest", "v1.0.0"},
		Platform:   []string{"linux/amd64", "linux/arm64"},
		Push:       true,
	})
}

// Apply builds images and applies Kubernetes manifests
func (Ko) Apply() error {
	return komagex.Apply(kox.ApplyOptions{
		Filenames: []string{"k8s/deployment.yaml"},
		Local:     false,
		Platform:  []string{"linux/amd64"},
	})
}

// ApplyLocal builds images locally and applies manifests
func (Ko) ApplyLocal() error {
	return komagex.Apply(kox.ApplyOptions{
		Filenames: []string{"k8s/deployment.yaml"},
		Local:     true,
		Platform:  []string{"linux/amd64"},
	})
}

// Delete deletes Kubernetes resources
func (Ko) Delete() error {
	return komagex.Delete(kox.DeleteOptions{
		Filenames: []string{"k8s/deployment.yaml"},
	})
}

// Publish publishes a container image
func (Ko) Publish() error {
	return komagex.Publish("./cmd/app")
}
