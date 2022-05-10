/*
Package cmd contains all the commands that are available in the test-inspector CLI.
Copyright © 2022 Egor Romanov egor@supabase.io

*/
package cmd

import (
	"fmt"
	"sync"
	"test-inspector/internal/supabase"
	"test-inspector/pkg/allure"
	"test-inspector/pkg/junit"
	"test-inspector/pkg/models"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	launch     string
	isTemplate bool
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload latest results to test-inspector",

	Run: func(cmd *cobra.Command, args []string) {
		if err := validateFlags(); err != nil {
			fmt.Printf("%v", err)
			return
		}

		supa, err := supabase.CreateClient(host, SupabaseKey, supabase.UserCredentials{
			Email:    user,
			Password: password,
		})
		if err != nil {
			fmt.Printf("error trying to connect to supabase: %v", err)
			return
		}
		if _, err = supa.GetVersion(versionID); err != nil {
			fmt.Printf("error trying to get version: %v", err)
			return
		}

		launchID, err := supa.CreateLaunch(models.Launch{
			IsTemplate: isTemplate,
			Name:       launch,
			VersionID:  int64(versionID),
		})
		if err != nil || launchID == 0 {
			fmt.Printf("error trying to create launch: %v", err)
			return
		}

		var results map[uuid.UUID]models.SupaResult
		switch reportType {
		case "allure":
			results, err = allure.ReadResults(resultsPath)
		case "junit":
			results, err = junit.ReadResults(resultsPath)
		default:
			err = fmt.Errorf("only 'junit' and 'allure' types supported")
		}
		if err != nil {
			fmt.Printf("error trying to parse results folder: %v", err)
			return
		}

		var wg sync.WaitGroup
		for _, r := range results {
			r := r
			wg.Add(1)
			go func() {
				defer wg.Done()
				r.LaunchID = launchID
				err := supa.CreateResult(r)
				if err != nil {
					fmt.Printf("problems with inserting Result: %s name - %s. %v", r.ID, r.Name, err)
					return
				}
			}()
		}
		wg.Wait()
		fmt.Println("upload succeed")
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringVarP(
		&launch, "launch", "l", "", "launch id (optional)")
	uploadCmd.Flags().BoolVarP(
		&isTemplate, "isReference", "r", false,
		"should this run be used as reference (requires admin permission)")

	viper.BindPFlag("launch", uploadCmd.Flags().Lookup("launch"))
	viper.BindPFlag("isReference", uploadCmd.Flags().Lookup("isReference"))
}
