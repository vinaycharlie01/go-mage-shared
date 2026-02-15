package kox

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/vinaycharlie01/go-mage-shared/execx"
)

// KoRunner handles ko command execution with dependency injection
type KoRunner struct {
	executor execx.Executor
}

// NewKoRunner creates a new KoRunner with the default executor
func NewKoRunner() *KoRunner {
	return &KoRunner{
		executor: execx.NewExec(),
	}
}

// NewKoRunnerWithExecutor creates a new KoRunner with a custom executor
func NewKoRunnerWithExecutor(executor execx.Executor) *KoRunner {
	return &KoRunner{
		executor: executor,
	}
}

// BuildOptions contains options for ko build
type BuildOptions struct {
	ImportPath          string   // Go import path to build
	Tags                []string // Image tags
	Platform            []string // Target platforms (e.g., linux/amd64,linux/arm64)
	BaseImage           string   // Base image to use
	Bare                bool     // Whether to use a bare image
	Local               bool     // Build locally without pushing
	Push                bool     // Push to registry
	PreserveImportPaths bool     // Preserve import paths in image names
}

// Build builds a container image using ko
func (k *KoRunner) Build(opts BuildOptions) error {
	if opts.ImportPath == "" {
		return fmt.Errorf("import path is required")
	}

	slog.Info("🐳 Building container image with ko...",
		"importPath", opts.ImportPath,
		"local", opts.Local,
		"push", opts.Push,
	)

	start := time.Now()

	args := []string{"build", opts.ImportPath}

	for _, tag := range opts.Tags {
		args = append(args, "--tags", tag)
	}

	for _, platform := range opts.Platform {
		args = append(args, "--platform", platform)
	}

	if opts.BaseImage != "" {
		args = append(args, "--base-import-paths", opts.BaseImage)
	}

	if opts.Bare {
		args = append(args, "--bare")
	}

	if opts.Local {
		args = append(args, "--local")
	}

	if opts.Push {
		args = append(args, "--push")
	}

	if opts.PreserveImportPaths {
		args = append(args, "--preserve-import-paths")
	}

	if err := k.executor.Run(context.Background(), "ko", false, args...); err != nil {
		return err
	}

	slog.Info("✅ Container image built", "duration", time.Since(start))
	return nil
}

// ApplyOptions contains options for ko apply
type ApplyOptions struct {
	Filenames           []string // Kubernetes manifest files
	Recursive           bool     // Process directories recursively
	Selector            string   // Label selector
	BaseImage           string   // Base image to use
	Platform            []string // Target platforms
	Local               bool     // Build locally without pushing
	Bare                bool     // Use bare image
	PreserveImportPaths bool     // Preserve import paths
}

// Apply builds images and applies Kubernetes manifests
func (k *KoRunner) Apply(opts ApplyOptions) error {
	if len(opts.Filenames) == 0 {
		return fmt.Errorf("at least one filename is required")
	}

	slog.Info("🚀 Building and applying with ko...",
		"files", opts.Filenames,
		"local", opts.Local,
	)

	start := time.Now()

	args := []string{"apply"}

	for _, filename := range opts.Filenames {
		args = append(args, "-f", filename)
	}

	if opts.Recursive {
		args = append(args, "--recursive")
	}

	if opts.Selector != "" {
		args = append(args, "--selector", opts.Selector)
	}

	if opts.BaseImage != "" {
		args = append(args, "--base-import-paths", opts.BaseImage)
	}

	for _, platform := range opts.Platform {
		args = append(args, "--platform", platform)
	}

	if opts.Local {
		args = append(args, "--local")
	}

	if opts.Bare {
		args = append(args, "--bare")
	}

	if opts.PreserveImportPaths {
		args = append(args, "--preserve-import-paths")
	}

	if err := k.executor.Run(context.Background(), "ko", false, args...); err != nil {
		return err
	}

	slog.Info("✅ Images built and manifests applied", "duration", time.Since(start))
	return nil
}

// DeleteOptions contains options for ko delete
type DeleteOptions struct {
	Filenames []string // Kubernetes manifest files
	Recursive bool     // Process directories recursively
	Selector  string   // Label selector
}

// Delete deletes Kubernetes resources
func (k *KoRunner) Delete(opts DeleteOptions) error {
	if len(opts.Filenames) == 0 {
		return fmt.Errorf("at least one filename is required")
	}

	slog.Info("🗑️  Deleting resources with ko...", "files", opts.Filenames)

	start := time.Now()

	args := []string{"delete"}

	for _, filename := range opts.Filenames {
		args = append(args, "-f", filename)
	}

	if opts.Recursive {
		args = append(args, "--recursive")
	}

	if opts.Selector != "" {
		args = append(args, "--selector", opts.Selector)
	}

	if err := k.executor.Run(context.Background(), "ko", false, args...); err != nil {
		return err
	}

	slog.Info("✅ Resources deleted", "duration", time.Since(start))
	return nil
}

// Resolve resolves import paths to image references
func (k *KoRunner) Resolve(importPaths []string, args ...string) error {
	if len(importPaths) == 0 {
		return fmt.Errorf("at least one import path is required")
	}

	slog.Info("🔍 Resolving import paths...", "paths", importPaths)

	start := time.Now()

	cmdArgs := []string{"resolve"}
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, importPaths...)

	if err := k.executor.Run(context.Background(), "ko", false, cmdArgs...); err != nil {
		return err
	}

	slog.Info("✅ Import paths resolved", "duration", time.Since(start))
	return nil
}

// Publish publishes images for import paths
func (k *KoRunner) Publish(importPath string, args ...string) error {
	if importPath == "" {
		return fmt.Errorf("import path is required")
	}

	slog.Info("📤 Publishing image...", "importPath", importPath)

	start := time.Now()

	cmdArgs := []string{"publish", importPath}
	cmdArgs = append(cmdArgs, args...)

	if err := k.executor.Run(context.Background(), "ko", false, cmdArgs...); err != nil {
		return err
	}

	slog.Info("✅ Image published", "duration", time.Since(start))
	return nil
}

// Made with Bob
