/*
   Copyright The containerd Authors.

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

package main

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var volumeRmCommand = &cli.Command{
	Name:         "rm",
	Aliases:      []string{"remove"},
	Usage:        "Remove one or more volumes",
	ArgsUsage:    "[flags] VOLUME [VOLUME, ...]",
	Description:  "NOTE: volume in use is deleted without caution",
	Action:       volumeRmAction,
	BashComplete: volumeRmBashComplete,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "(unimplemented yet)",
		},
	},
}

func volumeRmAction(clicontext *cli.Context) error {
	if clicontext.NArg() == 0 {
		return errors.Errorf("requires at least 1 argument")
	}
	volStore, err := getVolumeStore(clicontext)
	if err != nil {
		return err
	}
	names := clicontext.Args().Slice()
	removedNames, err := volStore.Remove(names)
	for _, removed := range removedNames {
		fmt.Fprintln(clicontext.App.Writer, removed)
	}
	return err
}

func volumeRmBashComplete(clicontext *cli.Context) {
	coco := parseCompletionContext(clicontext)
	if coco.boring || coco.flagTakesValue {
		defaultBashComplete(clicontext)
		return
	}
	// show voume names
	bashCompleteVolumeNames(clicontext)
}
