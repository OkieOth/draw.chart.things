package sub

import (
	"fmt"

	"github.com/spf13/cobra"
)

var BoxesCmd = &cobra.Command{
	Use:   "boxes",
	Short: "Draws boxes based on given layout",
	Long:  `Draws boxes and their connections and layouts them.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: I will draw boxes :D")
	},
}

func init() {
	initDefaultFlags(BoxesCmd)
}
