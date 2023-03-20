package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yagikota/git-multi/pkg/multiclone"
)

// multicloneCmd represents the multiclone command
var multicloneCmd = &cobra.Command{
	Use:   "multiclone",
	Short: "multiclone clones multiple git repositories in parallel",
	RunE: func(cmd *cobra.Command, args []string) error {
		multiclone.MultiClone(args)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(multicloneCmd)
}
