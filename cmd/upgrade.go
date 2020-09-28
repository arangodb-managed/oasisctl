//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
//
// Author Gergely Brautigam
//

package cmd

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	common "github.com/arangodb-managed/apis/common/v1"
	tools "github.com/arangodb-managed/apis/tools/v1"

	"github.com/arangodb-managed/oasisctl/pkg/format"
)

func init() {
	InitCommand(
		RootCmd,
		&cobra.Command{
			Use:   "upgrade",
			Short: "Upgrade Oasisctl tool",
			Long:  "Check the latest, compatible version and upgrade this tool to that.",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				dryRun bool
				force  bool
			}{}
			f.BoolVarP(&cargs.dryRun, "dry-run", "d", false, "Do an upgrade without applying the version.")
			f.BoolVarP(&cargs.force, "force", "f", false, "Force an upgrade even if the versions match.")

			c.Run = func(c *cobra.Command, args []string) {
				log := CLILog
				conn := MustDialAPI(WithoutVersionCheck())
				toolsc := tools.NewToolsServiceClient(conn)
				versions := collectCurrentAPIVersions()
				resp, err := toolsc.GetLatestVersion(context.Background(), &tools.GetLatestVersionRequest{
					Name:                tools.ToolNameOasisctl,
					ExpectedApiVersions: versions,
				})
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to get latest version for tool.")
				}
				if cargs.dryRun {
					fmt.Println("dry-run enabled, skipping applying version")
					fmt.Println(format.ToolsVersion(resp, RootArgs.Format))
					return
				}

				latestVersion := semver.New(strings.TrimPrefix(resp.GetLatestVersion(), "v"))
				if latestVersion.Equal(*currentVersion) && !cargs.force {
					log.Info().
						Str("latest_version", latestVersion.String()).
						Str("current_version", currentVersion.String()).
						Msg("Already using latest version. Nothing to do.")
					return
				}
				log.Info().Str("latest_version", resp.GetLatestVersion()).Msg("Applying latest version...")
				if err := upgradeBinary(log, resp.GetDownloadUrl()); err != nil {
					log.Fatal().Err(err).Msg("Error while upgrading to latest compatible version.")
				}
				log.Info().Msg("Done.")
			}
		},
	)
}

// upgradeBinary takes a url to download the latest release from, extracts it
// and moves the os and architecture based file to the current executables location.
func upgradeBinary(log zerolog.Logger, url string) error {
	log.Info().Msg("Downloading...")
	response, err := http.Get(url)
	if err != nil {
		log.Debug().Err(err).Str("url", url).Msg("Failed to get the download url.")
		return err
	}
	defer response.Body.Close()

	log.Info().Msg("done. Extracting...")
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to read body.")
		return err
	}
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		log.Debug().Err(err).Msg("Failed to create zip reader.")
		return err
	}

	dest, err := ioutil.TempDir("", "unzip")
	if err != nil {
		log.Debug().Err(err).Msg("Failed to create temp folder for decrompressing.")
		return err
	}

	ops := runtime.GOOS
	arch := runtime.GOARCH

	defer func() {
		if err := os.RemoveAll(dest); err != nil {
			log.Debug().Err(err).Str("folder", dest).Msg("Failed to perform cleanup. Please remove manually.")
		}
	}()

	currentExecutable, err := os.Executable()
	if err != nil {
		log.Debug().Err(err).Msg("Failed to find executable.")
		return err
	}
	originalPath := path.Dir(currentExecutable)
	execName := filepath.Base(os.Args[0])

	binaryFound := false
	for _, f := range zipReader.File {
		dir := path.Dir(f.Name)
		// only decompress the binary this system needs and don't bother with creating the folders.
		if !strings.Contains(filepath.ToSlash(dir), filepath.ToSlash(filepath.Join(ops, arch))) || f.Mode().IsDir() {
			continue
		}
		binaryFound = true
		filename := filepath.Base(f.Name)
		outFile, err := os.OpenFile(filepath.Join(dest, filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			log.Debug().Err(err).Str("dest", dest).Str("file", filename).Msg("Failed to create file.")
			return err
		}

		rc, err := f.Open()
		if err != nil {
			log.Debug().Err(err).Str("dest", dest).Str("file", filename).Msg("Failed to open file.")
			return err
		}
		if _, err := io.Copy(outFile, rc); err != nil {
			log.Debug().Err(err).Str("dest", dest).Str("file", filename).Msg("Failed to copy file content.")
			return err
		}

		outFile.Close()
		rc.Close()
	}
	if !binaryFound {
		return common.Unknown("No binary found for your os/architecture %s/%s", ops, arch)
	}
	log.Info().Msg("done. Updating binary...")

	// Rename the original running binary so updating on Windows can also work.
	old := filepath.Join(originalPath, execName+".old")
	if err := os.Rename(filepath.Join(originalPath, execName), old); err != nil {
		log.Debug().Err(err).Msg("Failed to rename current running binary.")
		return err
	}

	defer func() {
		if ops == "windows" {
			log.Info().Msg("Please remove binary with extension .old.")
			return
		}
		if err := os.RemoveAll(old); err != nil {
			log.Debug().Err(err).Str("old", old).Msg("Failed to perform cleanup. Please remove manually.")
		}
	}()

	if err := copyFile(log, filepath.Join(dest, execName), filepath.Join(originalPath, execName)); err != nil {
		log.Debug().Err(err).Msg("Failed to copy new version.")
		return err
	}
	return nil
}

// copyFile does a safe copy which works over partitions as well as locally.
func copyFile(log zerolog.Logger, src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		log.Debug().Err(err).Str("dest", dst).Msg("Failed to open src file.")
		return err
	}
	defer srcFile.Close()

	// get src file mode
	srcMode, err := os.Stat(srcFile.Name())
	if err != nil {
		log.Debug().Err(err).Str("dest", dst).Msg("Failed to stat src file.")
		return err
	}
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcMode.Mode())
	if err != nil {
		log.Debug().Err(err).Str("dest", dst).Msg("Failed to create destination file.")
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		log.Debug().Err(err).Str("dest", dst).Str("src", src).Msg("Failed to copy file content.")
		return err
	}
	return nil
}
