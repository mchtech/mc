// +build ignore

/*
 * Minio Client (C) 2014, 2015 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
)

var versionTemplate = `// --------  DO NOT EDIT --------
// This file is autogenerated by genversion.go during the release process.

package main

const mcVersion = {{if .Version}}"{{.Version}}"{{else}}""{{end}}
const mcReleaseTag = {{if .ReleaseTag}}"{{.ReleaseTag}}"{{else}}""{{end}}
`

// genVersion generates ‘version.go’.
func genVersion(version string) {
	t := template.Must(template.New("version").Parse(versionTemplate))
	versionFile, err := os.OpenFile("version.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("genversion: Unable to generate ‘version.go’. Error: %s.\n", err)
		os.Exit(1)
	}
	defer versionFile.Close()

	err = t.Execute(versionFile, struct {
		Version    string
		ReleaseTag string
	}{version, genReleaseTag(version)})

	if err != nil {
		fmt.Printf("genversion: Unable to generate ‘version.go’. Error: %s.\n", err)
		os.Exit(1)
	}
}

// genReleaseTag prints release tag to the console for easy git tagging.
func genReleaseTag(version string) string {
	relTag := strings.Replace(version, " ", "-", -1)
	relTag = strings.Replace(relTag, ":", "-", -1)
	relTag = strings.Replace(relTag, ",", "", -1)
	return "RELEASE." + relTag
}

func main() {
	// Version is in HTTP TimeFormat.
	version := time.Now().UTC().Format(http.TimeFormat)

	// generate ‘version.go’ file.
	genVersion(version)

	fmt.Println("Version: \"" + version + "\"")
	fmt.Println("Release-Tag: " + genReleaseTag(version))
}