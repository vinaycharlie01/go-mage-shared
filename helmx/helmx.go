package helmx

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/vinaycharlie01/go-mage-shared/execx"
)

// HelmRunner handles Helm command execution with dependency injection
type HelmRunner struct {
	executor execx.Executor
}

// NewHelmRunner creates a new HelmRunner with the default executor
func NewHelmRunner() *HelmRunner {
	return &HelmRunner{
		executor: execx.NewExec(),
	}
}

// NewHelmRunnerWithExecutor creates a new HelmRunner with a custom executor
func NewHelmRunnerWithExecutor(executor execx.Executor) *HelmRunner {
	return &HelmRunner{
		executor: executor,
	}
}

// InstallOptions contains options for helm install
type InstallOptions struct {
	ReleaseName     string
	Chart           string
	Namespace       string
	Values          []string // --values or -f flags
	Set             []string // --set flags
	CreateNamespace bool
	Wait            bool
	Timeout         string
}

// Install installs a Helm chart
func (h *HelmRunner) Install(opts InstallOptions) error {
	if opts.ReleaseName == "" {
		return fmt.Errorf("release name is required")
	}
	if opts.Chart == "" {
		return fmt.Errorf("chart is required")
	}

	slog.Info("📦 Installing Helm chart...",
		"release", opts.ReleaseName,
		"chart", opts.Chart,
		"namespace", opts.Namespace,
	)

	start := time.Now()

	args := []string{"install", opts.ReleaseName, opts.Chart}

	if opts.Namespace != "" {
		args = append(args, "--namespace", opts.Namespace)
	}

	if opts.CreateNamespace {
		args = append(args, "--create-namespace")
	}

	for _, valuesFile := range opts.Values {
		args = append(args, "--values", valuesFile)
	}

	for _, setValue := range opts.Set {
		args = append(args, "--set", setValue)
	}

	if opts.Wait {
		args = append(args, "--wait")
	}

	if opts.Timeout != "" {
		args = append(args, "--timeout", opts.Timeout)
	}

	if err := h.executor.Run(context.Background(), "helm", false, args...); err != nil {
		return err
	}

	slog.Info("✅ Helm chart installed", "duration", time.Since(start))
	return nil
}

// UpgradeOptions contains options for helm upgrade
type UpgradeOptions struct {
	ReleaseName string
	Chart       string
	Namespace   string
	Values      []string
	Set         []string
	Install     bool // --install flag
	Wait        bool
	Timeout     string
}

// Upgrade upgrades a Helm release
func (h *HelmRunner) Upgrade(opts UpgradeOptions) error {
	if opts.ReleaseName == "" {
		return fmt.Errorf("release name is required")
	}
	if opts.Chart == "" {
		return fmt.Errorf("chart is required")
	}

	slog.Info("🔄 Upgrading Helm release...",
		"release", opts.ReleaseName,
		"chart", opts.Chart,
		"namespace", opts.Namespace,
	)

	start := time.Now()

	args := []string{"upgrade", opts.ReleaseName, opts.Chart}

	if opts.Namespace != "" {
		args = append(args, "--namespace", opts.Namespace)
	}

	if opts.Install {
		args = append(args, "--install")
	}

	for _, valuesFile := range opts.Values {
		args = append(args, "--values", valuesFile)
	}

	for _, setValue := range opts.Set {
		args = append(args, "--set", setValue)
	}

	if opts.Wait {
		args = append(args, "--wait")
	}

	if opts.Timeout != "" {
		args = append(args, "--timeout", opts.Timeout)
	}

	if err := h.executor.Run(context.Background(), "helm", false, args...); err != nil {
		return err
	}

	slog.Info("✅ Helm release upgraded", "duration", time.Since(start))
	return nil
}

// Uninstall uninstalls a Helm release
func (h *HelmRunner) Uninstall(releaseName, namespace string, args ...string) error {
	if releaseName == "" {
		return fmt.Errorf("release name is required")
	}

	slog.Info("🗑️  Uninstalling Helm release...",
		"release", releaseName,
		"namespace", namespace,
	)

	start := time.Now()

	cmdArgs := []string{"uninstall", releaseName}

	if namespace != "" {
		cmdArgs = append(cmdArgs, "--namespace", namespace)
	}

	cmdArgs = append(cmdArgs, args...)

	if err := h.executor.Run(context.Background(), "helm", false, cmdArgs...); err != nil {
		return err
	}

	slog.Info("✅ Helm release uninstalled", "duration", time.Since(start))
	return nil
}

// List lists Helm releases
func (h *HelmRunner) List(namespace string, args ...string) error {
	slog.Info("📋 Listing Helm releases...", "namespace", namespace)

	start := time.Now()

	cmdArgs := []string{"list"}

	if namespace != "" {
		cmdArgs = append(cmdArgs, "--namespace", namespace)
	} else {
		cmdArgs = append(cmdArgs, "--all-namespaces")
	}

	cmdArgs = append(cmdArgs, args...)

	if err := h.executor.Run(context.Background(), "helm", false, cmdArgs...); err != nil {
		return err
	}

	slog.Info("✅ Helm releases listed", "duration", time.Since(start))
	return nil
}

// Status shows the status of a Helm release
func (h *HelmRunner) Status(releaseName, namespace string, args ...string) error {
	if releaseName == "" {
		return fmt.Errorf("release name is required")
	}

	slog.Info("📊 Getting Helm release status...",
		"release", releaseName,
		"namespace", namespace,
	)

	start := time.Now()

	cmdArgs := []string{"status", releaseName}

	if namespace != "" {
		cmdArgs = append(cmdArgs, "--namespace", namespace)
	}

	cmdArgs = append(cmdArgs, args...)

	if err := h.executor.Run(context.Background(), "helm", false, cmdArgs...); err != nil {
		return err
	}

	slog.Info("✅ Helm release status retrieved", "duration", time.Since(start))
	return nil
}

// Template renders chart templates locally
func (h *HelmRunner) Template(releaseName, chart string, args ...string) error {
	if releaseName == "" {
		return fmt.Errorf("release name is required")
	}
	if chart == "" {
		return fmt.Errorf("chart is required")
	}

	slog.Info("📝 Rendering Helm templates...",
		"release", releaseName,
		"chart", chart,
	)

	start := time.Now()

	cmdArgs := []string{"template", releaseName, chart}
	cmdArgs = append(cmdArgs, args...)

	if err := h.executor.Run(context.Background(), "helm", false, cmdArgs...); err != nil {
		return err
	}

	slog.Info("✅ Helm templates rendered", "duration", time.Since(start))
	return nil
}

// Lint runs helm lint on a chart
func (h *HelmRunner) Lint(chart string, args ...string) error {
	if chart == "" {
		return fmt.Errorf("chart path is required")
	}

	slog.Info("🔍 Linting Helm chart...", "chart", chart)

	start := time.Now()

	cmdArgs := []string{"lint", chart}
	cmdArgs = append(cmdArgs, args...)

	if err := h.executor.Run(context.Background(), "helm", false, cmdArgs...); err != nil {
		return err
	}

	slog.Info("✅ Helm chart linted", "duration", time.Since(start))
	return nil
}

// Package packages a chart directory into a chart archive
func (h *HelmRunner) Package(chart string, args ...string) error {
	if chart == "" {
		return fmt.Errorf("chart path is required")
	}

	slog.Info("📦 Packaging Helm chart...", "chart", chart)

	start := time.Now()

	cmdArgs := []string{"package", chart}
	cmdArgs = append(cmdArgs, args...)

	if err := h.executor.Run(context.Background(), "helm", false, cmdArgs...); err != nil {
		return err
	}

	slog.Info("✅ Helm chart packaged", "duration", time.Since(start))
	return nil
}

// RepoAdd adds a chart repository
func (h *HelmRunner) RepoAdd(name, url string, args ...string) error {
	if name == "" {
		return fmt.Errorf("repository name is required")
	}
	if url == "" {
		return fmt.Errorf("repository URL is required")
	}

	slog.Info("➕ Adding Helm repository...", "name", name, "url", url)

	start := time.Now()

	cmdArgs := []string{"repo", "add", name, url}
	cmdArgs = append(cmdArgs, args...)

	if err := h.executor.Run(context.Background(), "helm", false, cmdArgs...); err != nil {
		return err
	}

	slog.Info("✅ Helm repository added", "duration", time.Since(start))
	return nil
}

// RepoUpdate updates chart repositories
func (h *HelmRunner) RepoUpdate(args ...string) error {
	slog.Info("🔄 Updating Helm repositories...")

	start := time.Now()

	cmdArgs := []string{"repo", "update"}
	cmdArgs = append(cmdArgs, args...)

	if err := h.executor.Run(context.Background(), "helm", false, cmdArgs...); err != nil {
		return err
	}

	slog.Info("✅ Helm repositories updated", "duration", time.Since(start))
	return nil
}

// Made with Bob
