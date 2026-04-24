// Package main is the entry point for the envsync CLI tool.
// It provides commands to diff and sync .env files across environments
// with optional secret masking support.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envsync/internal/envfile"
)

const usage = `envsync - diff and sync .env files across environments

Usage:
  envsync diff <source> <target>     Show differences between two .env files
  envsync sync <source> <target>     Sync keys from source into target

Flags:
`

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		fmt.Print(usage)
		flag.PrintDefaults()
		return nil
	}

	switch args[0] {
	case "diff":
		return runDiff(args[1:])
	case "sync":
		return runSync(args[1:])
	case "help", "--help", "-h":
		fmt.Print(usage)
		return nil
	default:
		return fmt.Errorf("unknown command %q — run 'envsync help' for usage", args[0])
	}
}

func runDiff(args []string) error {
	fs := flag.NewFlagSet("diff", flag.ContinueOnError)
	maskSecrets := fs.Bool("mask", true, "mask sensitive values in output")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 2 {
		return fmt.Errorf("diff requires two arguments: <source> <target>")
	}

	sourcePath := fs.Arg(0)
	targetPath := fs.Arg(1)

	source, err := envfile.Parse(sourcePath)
	if err != nil {
		return fmt.Errorf("parsing source %q: %w", sourcePath, err)
	}

	target, err := envfile.Parse(targetPath)
	if err != nil {
		return fmt.Errorf("parsing target %q: %w", targetPath, err)
	}

	result := envfile.Diff(source, target)
	if !envfile.HasChanges(result) {
		fmt.Println("No differences found.")
		return nil
	}

	masker := envfile.NewMasker(nil)
	fmt.Print(envfile.FormatDiff(result, *maskSecrets, masker))
	return nil
}

func runSync(args []string) error {
	fs := flag.NewFlagSet("sync", flag.ContinueOnError)
	overwrite := fs.Bool("overwrite", false, "overwrite existing keys in target")
	dryRun := fs.Bool("dry-run", false, "preview changes without writing to disk")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 2 {
		return fmt.Errorf("sync requires two arguments: <source> <target>")
	}

	sourcePath := fs.Arg(0)
	targetPath := fs.Arg(1)

	source, err := envfile.Parse(sourcePath)
	if err != nil {
		return fmt.Errorf("parsing source %q: %w", sourcePath, err)
	}

	target, err := envfile.Parse(targetPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("parsing target %q: %w", targetPath, err)
	}

	synced, err := envfile.Sync(source, target, targetPath, envfile.SyncOptions{
		Overwrite: *overwrite,
		DryRun:    *dryRun,
	})
	if err != nil {
		return fmt.Errorf("syncing files: %w", err)
	}

	if *dryRun {
		fmt.Println("[dry-run] The following changes would be applied:")
		masker := envfile.NewMasker(nil)
		fmt.Print(envfile.Format(synced, true, masker))
	} else {
		fmt.Printf("Synced %d entries from %q into %q\n", len(synced), sourcePath, targetPath)
	}
	return nil
}
