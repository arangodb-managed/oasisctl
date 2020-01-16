//
// DISCLAIMER
//
// Copyright 2020 ArangoDB Inc, Cologne, Germany
//
// Author Gergely Brautigam
//

package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path"
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
		filePrepend string
		outputDir   string
	}
)

func init() {
	RootCmd.AddCommand(GenerateCmd)
	f := GenerateCmd.Flags()
	f.StringVarP(&generateArgs.filePrepend, "prepend", "p", "", "Content to preppend to the generated content")
	f.StringVarP(&generateArgs.outputDir, "output-dir", "o", "./docs", "Output directory")
}

func generateMarkdownRun(c *cobra.Command, args []string) {
	var prepend string

	if generateArgs.filePrepend != "" {
		content, err := ioutil.ReadFile(generateArgs.filePrepend)
		if err != nil {
			log.Fatal(err)
		}
		prepend = string(content)
	}
	filePrepender := func(filename string) string {
		return prepend
	}
	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		return "/commands/" + strings.ToLower(base) + "/"
	}

	if _, err := os.Stat(generateArgs.outputDir); os.IsNotExist(err) {
		log.Fatalf("Directory %s does not exist.", generateArgs.outputDir)
	}

	err := doc.GenMarkdownTreeCustom(RootCmd, generateArgs.outputDir, filePrepender, linkHandler)
	if err != nil {
		log.Fatal("Failed to generate markdown.", err)
	}
}
