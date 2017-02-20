// Copyright © 2017 Sam Tang <samyftang@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yftang/seqapp/batch/rna"
)

// batchCmd represents the batch command
var batchCmd = &cobra.Command{
	Use:   "batch",
	Short: "Generates series of pbs files for similar bioinformatics analysis",
	Long: `If you have multiple samples/files which need analysis for similar pipeline,
you should need this command. For example:

seqapp batch diff -a RNA
`,
}

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Generates series of pbs file for diff analysis",
	Run: func(cmd *cobra.Command, args []string) {
		app, appErr := cmd.Flags().GetString("app")
		if appErr != nil {
			fmt.Println("Flag app not found")
		}
		conf, confErr := cmd.Flags().GetString("conf")
		if conf == "" || confErr != nil {
			fmt.Println("Error: Param conf not found")
			os.Exit(-1)
		}
		user, userErr := cmd.Flags().GetString("user")
		if user == "" || userErr != nil {
			fmt.Println("Error: Param user not found")
			os.Exit(-1)
		}

		switch app {
		case "RNA":
			dc := rna.DiffConfig{ConfigFile: conf}

			if err := dc.Parse(); err != nil {
				fmt.Println(err)
			}

			diff := rna.Diff{Config: dc}
			if err := diff.GenPbsFiles(user); err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Println("Error: Please specify the application name, such as RNA")
		}
	},
}

func init() {
	RootCmd.AddCommand(batchCmd)

	batchCmd.PersistentFlags().StringP("app", "a", "", "Application to use, such as 'RNA'")
	batchCmd.PersistentFlags().StringP("conf", "c", "", "Config file path")
	batchCmd.PersistentFlags().StringP("user", "u", "", "Username of cluster")

	batchCmd.AddCommand(diffCmd)
}
