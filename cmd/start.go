package cmd

import (
	"fmt"
	"github.com/SwapnilDaddikar07/todo/app"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use: "start",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := app.NewDefaultStore()
		if err != nil {
			return err
		}

		view := app.NewView(store)

		err = view.Build()
		if err != nil {
			fmt.Printf("error building view %v", err)
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
