package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsync/internal/envfile"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: envsync <command> [options]\ncommands: diff, sync, validate")
	}
	switch args[0] {
	case "diff":
		return runDiff(args[1:])
	case "sync":
		return runSync(args[1:])
	case "validate":
		return runValidate(args[1:])
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func runDiff(args []string) error {
	fs := flag.NewFlagSet("diff", flag.ContinueOnError)
	maskSecrets := fs.Bool("mask", true, "mask secret values in output")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 2 {
		return fmt.Errorf("diff requires two file arguments")
	}
	src, err := envfile.Parse(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("parsing source: %w", err)
	}
	dst, err := envfile.Parse(fs.Arg(1))
	if err != nil {
		return fmt.Errorf("parsing destination: %w", err)
	}
	result := envfile.Diff(src, dst)
	masker := envfile.NewMasker(nil)
	fmt.Print(envfile.FormatDiff(result, *maskSecrets, masker))
	return nil
}

func runSync(args []string) error {
	fs := flag.NewFlagSet("sync", flag.ContinueOnError)
	overwrite := fs.Bool("overwrite", false, "overwrite existing keys in destination")
	dryRun := fs.Bool("dry-run", false, "print changes without writing")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 2 {
		return fmt.Errorf("sync requires source and destination file arguments")
	}
	src, err := envfile.Parse(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("parsing source: %w", err)
	}
	dst, err := envfile.Parse(fs.Arg(1))
	if err != nil {
		return fmt.Errorf("parsing destination: %w", err)
	}
	opts := envfile.SyncOptions{Overwrite: *overwrite, DryRun: *dryRun}
	updated, err := envfile.Sync(src, dst, fs.Arg(1), opts)
	if err != nil {
		return fmt.Errorf("syncing: %w", err)
	}
	if *dryRun {
		masker := envfile.NewMasker(nil)
		fmt.Print(envfile.Format(updated, false, masker))
	}
	return nil
}

func runValidate(args []string) error {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	required := fs.String("required", "", "comma-separated list of required keys")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("validate requires a file argument")
	}
	entries, err := envfile.Parse(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("parsing file: %w", err)
	}
	result := envfile.Validate(entries)
	if *required != "" {
		keys := strings.Split(*required, ",")
		for i, k := range keys {
			keys[i] = strings.TrimSpace(k)
		}
		reqResult := envfile.ValidateRequiredKeys(entries, keys)
		result.Errors = append(result.Errors, reqResult.Errors...)
	}
	if !result.Valid() {
		for _, e := range result.Errors {
			fmt.Fprintf(os.Stderr, "validation error: %s\n", e.Error())
		}
		return fmt.Errorf("%d validation error(s) found", len(result.Errors))
	}
	fmt.Println("validation passed")
	return nil
}
