package komagex

import (
	"github.com/vinaycharlie01/go-mage-shared/kox"
)

// Package-level convenience functions for mage targets
var defaultRunner = kox.NewKoRunner()

// Build builds a container image using ko
func Build(opts kox.BuildOptions) error {
	return defaultRunner.Build(opts)
}

// Apply builds images and applies Kubernetes manifests
func Apply(opts kox.ApplyOptions) error {
	return defaultRunner.Apply(opts)
}

// Delete deletes Kubernetes resources
func Delete(opts kox.DeleteOptions) error {
	return defaultRunner.Delete(opts)
}

// Resolve resolves import paths to image references
func Resolve(importPaths []string, args ...string) error {
	return defaultRunner.Resolve(importPaths, args...)
}

// Publish publishes images for import paths
func Publish(importPath string, args ...string) error {
	return defaultRunner.Publish(importPath, args...)
}

// Made with Bob
