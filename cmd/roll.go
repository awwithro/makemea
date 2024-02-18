package cmd

import (
	"errors"
	"fmt"

	"github.com/justinian/dice"
	"github.com/spf13/cobra"
)

var rollCmd = &cobra.Command{
	Use:   "roll [dice string xdy[+z]]",
	Short: "Roll some dice ",
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args{
			result, _, err := dice.Roll(arg)
			if err != nil{
				fmt.Printf("%v\n", err)
			}else{
				fmt.Println(result.Int())
			}
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("No dice string given")
		}
		return nil
	},
}
