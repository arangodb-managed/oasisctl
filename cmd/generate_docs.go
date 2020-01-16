//
// DISCLAIMER
//
// Copyright 2020 ArangoDB Inc, Cologne, Germany
//
// Author Gergely Brautigam
//

package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	// GenerateCmd is root for various `generate ...` commands
	GenerateCmd = &cobra.Command{
		Use:                "generate-docs",
		Short:              "Generate output",
		Run:                generateMarkdownRun,
		DisableAutoGenTag:  true,
		DisableSuggestions: true,
	}
	generateArgs struct {
		outputDir string
	}
)

const fmTemplate = `---
layout: default
description: %s
title: %s
---
`

func init() {
	RootCmd.AddCommand(GenerateCmd)
	f := GenerateCmd.Flags()
	f.StringVarP(&generateArgs.outputDir, "output-dir", "o", "./docs", "Output directory")
}

func generateMarkdownRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := CLILog

	filePrepender := func(filename string) string {
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		command := strings.Replace(base, "_", " ", -1)
		title := strings.Title(command)
		description := "Description of the " + command + " command"
		return fmt.Sprintf(fmTemplate, description, title)
	}
	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		return strings.Replace(base, "_", "-", -1) + ".html"
	}

	if _, err := os.Stat(generateArgs.outputDir); os.IsNotExist(err) {
		log.Fatal().Err(err).Str("output", generateArgs.outputDir).Msg("Output directory does not exist.")
	}

	err := doc.GenMarkdownTreeCustom(RootCmd, generateArgs.outputDir, filePrepender, linkHandler)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to generate document")
	}
}
