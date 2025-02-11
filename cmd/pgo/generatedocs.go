// +build ignore

package main

/*
 Copyright 2019 - 2021 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra/doc"

	"github.com/qingcloud/postgres-operator/cmd/pgo/cmd"
)

const fmTemplate = `---
title: "%s"
---
`

func main() {
	fmt.Println("generate CLI markdown")

	filePrepender := func(filename string) string {
		// now := time.Now().Format(time.RFC3339)
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		fmt.Println(base)
		// url := "/commands/" + strings.ToLower(base) + "/"
		return fmt.Sprintf(fmTemplate, strings.ReplaceAll(base, "_", " "))
	}

	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		return "/pgo-client/reference/" + strings.ToLower(base) + "/"
	}

	err := doc.GenMarkdownTreeCustom(cmd.RootCmd, "./", filePrepender, linkHandler)
	if err != nil {
		log.Fatal(err)
	}
}
