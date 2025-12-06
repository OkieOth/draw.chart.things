package sub

import (
	"github.com/spf13/cobra"
)

var From string
var Output string

func initDefaultFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&From, "from", "f", "", "Path to the input file (required)")
	cmd.PersistentFlags().StringVarP(&Output, "output", "o", "", "Output to create, default is stdout")
	cmd.MarkPersistentFlagRequired("from")
	cmd.MarkPersistentFlagFilename("from")
	cmd.MarkPersistentFlagFilename("output")
}
