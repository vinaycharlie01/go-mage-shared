package golang

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vinaycharlie01/go-mage-shared/execx"
)

// GoRunner handles Go command execution with dependency injection
type GoRunner struct {
	executor execx.Executor
}

// NewGoRunner creates a new GoRunner with the default executor
func NewGoRunner() *GoRunner {
	return &GoRunner{
		executor: execx.NewExec(),
	}
}

// NewGoRunnerWithExecutor creates a new GoRunner with a custom executor
func NewGoRunnerWithExecutor(executor execx.Executor) *GoRunner {
	return &GoRunner{
		executor: executor,
	}
}

// RunTests runs Go tests with given arguments
func (g *GoRunner) RunTests(args ...string) error {
	slog.Info("🧪 Running Go Tests...")
	defaultArgs := []string{"test", "./..."}
	start := time.Now()
	if err := g.executor.Run(context.Background(), "go", false, append(defaultArgs, args...)...); err != nil {
		return err
	}
	slog.Info("✅ Tests passed", "duration", time.Since(start))
	return nil
}

// RunLint runs golangci-lint with given arguments
func (g *GoRunner) RunLint(args ...string) error {
	slog.Info("🔍 Running Go Linter...")
	defaultArgs := []string{"run", "--timeout=5m"}
	start := time.Now()
	if err := g.executor.Run(context.Background(), "golangci-lint", false, append(defaultArgs, args...)...); err != nil {
		return err
	}
	slog.Info("✅ Lint passed", "duration", time.Since(start))
	return nil
}

// RunInstall installs Go packages
func (g *GoRunner) RunInstall(pkgs []string, args ...string) error {
	if len(pkgs) == 0 {
		return fmt.Errorf("no package specified for installation")
	}

	slog.Info("📦 Installing Go packages individually...", "packages", pkgs)

	start := time.Now()
	for _, pkg := range pkgs {
		cmdArgs := append([]string{"install", pkg}, args...)
		if err := g.executor.Run(context.Background(), "go", false, cmdArgs...); err != nil {
			return fmt.Errorf("failed to install %s: %w", pkg, err)
		}
	}

	slog.Info("✅ Installation complete", "duration", time.Since(start))
	return nil
}

// RunModTasks runs `go mod tidy` and `go mod verify` sequentially
func (g *GoRunner) RunModTasks() error {
	slog.Info("📦 Running Go module maintenance (tidy & verify)...")

	start := time.Now()

	commands := [][]string{
		{"mod", "tidy"},
		{"mod", "verify"},
	}

	for _, args := range commands {
		slog.Info("🔧 Executing", "command", fmt.Sprintf("go %s", strings.Join(args, " ")))
		if err := g.executor.Run(context.Background(), "go", false, args...); err != nil {
			return fmt.Errorf("failed to run 'go %s': %w", strings.Join(args, " "), err)
		}
	}
	slog.Info("✅ Module maintenance completed successfully", "duration", time.Since(start))
	return nil
}

// Run runs go mod tidy
func (g *GoRunner) Run() error {
	slog.Info("🧪 Running Go Mod Tidy...")
	defaultArgs := []string{"mod", "tidy"}
	start := time.Now()
	if err := g.executor.Run(context.Background(), "go", false, defaultArgs...); err != nil {
		return err
	}
	slog.Info("✅ Tests passed", "duration", time.Since(start))
	return nil
}

// Package-level convenience functions for backward compatibility
var defaultRunner = NewGoRunner()

// RunTests runs Go tests with given arguments
func RunTests(args ...string) error {
	return defaultRunner.RunTests(args...)
}

// RunLint runs golangci-lint with given arguments
func RunLint(args ...string) error {
	return defaultRunner.RunLint(args...)
}

// RunInstall installs Go packages
func RunInstall(pkgs []string, args ...string) error {
	return defaultRunner.RunInstall(pkgs, args...)
}

// RunModTasks runs `go mod tidy` and `go mod verify` sequentially
func RunModTasks() error {
	return defaultRunner.RunModTasks()
}

// Run runs go mod tidy
func Run() error {
	return defaultRunner.Run()
}

type BuildOptions struct {
	Binary         string
	Version        string
	OS             string
	Arch           string
	Debug          bool
	Packages       []string
	DestinationDir string // NEW
}

// RunBuild builds a Go binary with the given options
func (g *GoRunner) RunBuild(opts BuildOptions) error {
	if opts.Binary == "" {
		return fmt.Errorf("binary name is required")
	}
	if len(opts.Packages) == 0 {
		opts.Packages = []string{"."}
	}

	destDir := opts.DestinationDir
	if destDir == "" {
		destDir = "dist/binaries"
	}

	slog.Info("🏗️ Building Go binary...",
		"binary", opts.Binary,
		"os", opts.OS,
		"arch", opts.Arch,
		"debug", opts.Debug,
	)

	start := time.Now()

	// ---- ldflags ----
	ldflags := fmt.Sprintf("-X main.version=%s", opts.Version)
	if !opts.Debug {
		ldflags += " -s -w"
	}

	// ---- output path ----
	outDir := filepath.Join(destDir, opts.OS+"_"+opts.Arch)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	outPath := filepath.Join(outDir, opts.Binary)

	// ---- go build args ----
	buildArgs := []string{
		"GOOS=" + opts.OS,
		"GOARCH=" + opts.Arch,
		"CGO_ENABLED=0",
		"go",
		"build",
		"-ldflags", ldflags,
		"-o", outPath,
	}
	buildArgs = append(buildArgs, opts.Packages...)

	// ---- runtime-only env execution ----
	if err := g.executor.Run(
		context.Background(),
		"env",
		false,
		buildArgs...,
	); err != nil {
		return err
	}

	slog.Info("✅ Build completed",
		"output", outPath,
		"duration", time.Since(start),
	)

	return nil
}

// RunTestsWithCoverage runs Go tests with coverage
func (g *GoRunner) RunTestsWithCoverage(args ...string) error {
	slog.Info("🧪 Running tests with coverage...")
	defaultArgs := []string{"test", "-cover", "-coverprofile=coverage.out", "./..."}
	start := time.Now()
	if err := g.executor.Run(context.Background(), "go", false, append(defaultArgs, args...)...); err != nil {
		return err
	}
	slog.Info("✅ Tests with coverage passed", "duration", time.Since(start))
	return nil
}

// RunVet runs go vet
func (g *GoRunner) RunVet(args ...string) error {
	slog.Info("🔍 Running go vet...")
	defaultArgs := []string{"vet", "./..."}
	start := time.Now()
	if err := g.executor.Run(context.Background(), "go", false, append(defaultArgs, args...)...); err != nil {
		return err
	}
	slog.Info("✅ Go vet passed", "duration", time.Since(start))
	return nil
}

// RunFormat formats Go files using gofmt
func (g *GoRunner) RunFormat(args ...string) error {
	slog.Info("✨ Formatting Go files...")
	defaultArgs := []string{"-w", "."}
	start := time.Now()
	if err := g.executor.Run(context.Background(), "gofmt", false, append(defaultArgs, args...)...); err != nil {
		return err
	}
	slog.Info("✅ Formatting complete", "duration", time.Since(start))
	return nil
}

// RunFormatImports formats Go imports using goimports
func (g *GoRunner) RunFormatImports(args ...string) error {
	slog.Info("✨ Formatting Go imports...")
	defaultArgs := []string{"-w", "."}
	start := time.Now()
	if err := g.executor.Run(context.Background(), "goimports", false, append(defaultArgs, args...)...); err != nil {
		return err
	}
	slog.Info("✅ Import formatting complete", "duration", time.Since(start))
	return nil
}

// RunBuild builds a Go binary with the given options (package-level convenience function)
func RunBuild(opts BuildOptions) error {
	return defaultRunner.RunBuild(opts)
}

// RunTestsWithCoverage runs Go tests with coverage (package-level convenience function)
func RunTestsWithCoverage(args ...string) error {
	return defaultRunner.RunTestsWithCoverage(args...)
}

// RunVet runs go vet (package-level convenience function)
func RunVet(args ...string) error {
	return defaultRunner.RunVet(args...)
}

// RunFormat formats Go files (package-level convenience function)
func RunFormat(args ...string) error {
	return defaultRunner.RunFormat(args...)
}

// RunFormatImports formats Go imports (package-level convenience function)
func RunFormatImports(args ...string) error {
	return defaultRunner.RunFormatImports(args...)
}
