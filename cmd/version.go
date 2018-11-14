////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

//go:generate go run gen.go
// The above generates: GITVERSION, GLIDEDEPS, and SEMVER

func init() {
	RootCmd.AddCommand(versionCmd)
}

func printVersion() {
	fmt.Printf("Elixxir User Discovery Bot v%s -- %s\n\n", SEMVER,
		GITVERSION)
	fmt.Printf("Dependencies:\n\n%s\n", GLIDEDEPS)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Elixxir UDB",
	Long: `Print the version number of Elixxir User Discovery bot. This
also prints the glide cache versions of all of its dependencies.`,
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}
