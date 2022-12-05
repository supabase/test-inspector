/*
Package cmd contains all the commands that are available in the test-inspector CLI.
Copyright © 2022 Egor Romanov egor@supabase.io

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"test-inspector/internal/supabase"
	"test-inspector/pkg/color"
	"test-inspector/pkg/models"

	"github.com/spf13/cobra"
)

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "print reference test results for your project",
	Run: func(cmd *cobra.Command, args []string) {
		if err := validateVersionID(); err != nil {
			fmt.Printf("%v", err)
			return
		}

		supa, err := supabase.CreateClient(host, SupabaseKey, supabase.UserCredentials{})
		if err != nil {
			fmt.Printf("error trying to connect to supabase: %v", err)
			return
		}
		if _, err = supa.GetVersion(versionID); err != nil {
			fmt.Printf("error trying to get version: %v", err)
			return
		}

		templates, err := supa.GetTemplate(int64(versionID))
		if err != nil {
			fmt.Printf("error trying to retrieve template test results: %v", err)
			return
		}

		features, err := supa.GetFeatures()
		if err != nil {
			fmt.Printf("error trying to retrieve features for test results: %v", err)
			return
		}

		fmt.Print("\n" + color.Blue + "Reference Test Results:\n\n" + color.Reset)

		fmt.Printf("%s%d%s test results found in reference run\n\n\n",
			color.Green, len(templates), color.Reset)

		for _, feature := range features {
			ctr := 0
			fmt.Printf("%s%s%s\n\n", color.Blue, feature, color.Reset)
			for _, t := range templates {
				if t.Feature == feature ||
					t.Suite == feature ||
					t.ParentSuite == feature ||
					t.SubSuite == feature {
					ctr++
					fmt.Printf("\t%d. %s%s%s\n", ctr, color.Green, t.Name, color.Reset)

					if t.Steps != "" {
						var templateSteps []*models.StepContainer
						err = json.Unmarshal([]byte(t.Steps), &templateSteps)
						if err == nil {
							for _, step := range templateSteps {
								fmt.Printf("\t\t%d. %s%s%s\n", step.Position+1, color.Yellow, step.Name, color.Reset)
								printSteps(step, 3)
							}
						}
					}
					fmt.Print("\n")
				}
			}
			fmt.Print("\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(printCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// printCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// printCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func printSteps(parent *models.StepContainer, depth int) {
	for _, step := range parent.StepContainer {
		tabs := strings.Repeat("\t", depth)
		fmt.Printf("%s%d. %s%s%s\n", tabs, step.Position+1, color.Yellow, step.Name, color.Reset)
		printSteps(step, depth+1)
	}
}
