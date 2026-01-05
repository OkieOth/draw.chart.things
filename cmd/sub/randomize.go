package sub

import (
	"fmt"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/brianvoe/gofakeit/v7/source"
	"github.com/okieoth/draw.chart.things/pkg/boxesimpl"
	"github.com/okieoth/draw.chart.things/pkg/types/boxes"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var RandomizeCmd = &cobra.Command{
	Use:   "randomize",
	Short: "Randomize the texts in input files",
	Long:  `Randomize the texts in input files`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("call this command with one of the available sub commands")
	},
}

var RandomizeBoxesCmd = &cobra.Command{
	Use:   "boxes",
	Short: "Randomize the texts in boxes input files",
	Long:  `Randomize the texts in boxes input files`,
	Run: func(cmd *cobra.Command, args []string) {
		boxes, err := boxesimpl.LoadBoxesFromFile(From)
		if err != nil {
			fmt.Println("Error while loading input:", err)
			os.Exit(1)
		}
		faker := gofakeit.NewFaker(source.NewJSF(11), true)
		randomizeBoxesLayout(&boxes.Boxes, faker)
		bytes, err := yaml.Marshal(boxes)
		if err != nil {
			fmt.Println("Error while serialize to yaml:", err)
			os.Exit(1)
		}
		if Output == "" {
			fmt.Println(string(bytes))
		} else {
			output, err := os.Create(Output)
			if err != nil {
				fmt.Println("Error while creating output file:", err)
				os.Exit(1)
			}
			defer output.Close()
			_, err = output.Write(bytes)
			if err != nil {
				fmt.Println("Error while writing output file:", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	initDefaultFlags(RandomizeCmd)
	RandomizeCmd.AddCommand(RandomizeBoxesCmd)
}

func randomizeText(txt string, faker *gofakeit.Faker) string {
	if txt == "" {
		return txt
	}
	return faker.Sentence()
}

func randomizeBoxesLayoutContainer(cont []boxes.Layout, faker *gofakeit.Faker) []boxes.Layout {
	if len(cont) == 0 {
		return []boxes.Layout{}
	}
	ret := make([]boxes.Layout, 0)
	for _, c := range cont {
		c.Caption = randomizeText(c.Caption, faker)
		c.Text1 = randomizeText(c.Text1, faker)
		c.Text2 = randomizeText(c.Text2, faker)
		ret = append(ret, c)
	}
	return ret
}

func randomizeBoxesLayout(l *boxes.Layout, faker *gofakeit.Faker) {
	l.Caption = randomizeText(l.Caption, faker)
	l.Text1 = randomizeText(l.Text1, faker)
	l.Text2 = randomizeText(l.Text2, faker)
	l.Horizontal = randomizeBoxesLayoutContainer(l.Horizontal, faker)
	l.Vertical = randomizeBoxesLayoutContainer(l.Vertical, faker)
}
