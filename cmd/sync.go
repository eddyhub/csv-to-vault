// Copyright © 2016 Eduard Angold eddyhub@users.noreply.github.com
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
	"github.com/eddyhub/csv_to_vault/sync"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var csvFile, branch string

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "syncs csv to our vaultserver",
	Long:  `This command syncs a specfied csv to our vaultserver`,
	Run: func(cmd *cobra.Command, args []string) {

		_, data := sync.ReadCsv(csvFile)
		write := sync.DataWriter(
			viper.GetString("address"),
			viper.GetString("caCert"),
			viper.GetBool("insecure"),
			viper.GetString("tlsServerName"),
			viper.GetString("token"))
		for _, i := range data {
			fmt.Printf("Syncing: /secret/%s/%s/%s password=%s\n", branch, i[0], i[1], strings.Repeat("*", 8))
			write(branch, i[0], i[1], i[2])
		}
	},
}

func init() {
	RootCmd.AddCommand(syncCmd)
	syncCmd.Flags().StringVar(&csvFile, "csv-file", "", "-csv-file <path to csv file>")
	syncCmd.Flags().StringVar(&branch, "branch", "db", "-branch db -> /secret/db/...")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
