/*
Package cmd contains all the commands that are available in the test-inspector CLI.
Copyright © 2022 Egor Romanov egor@supabase.io

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"test-inspector/internal/supabase"
	"test-inspector/pkg/allure"
	"test-inspector/pkg/color"
	"test-inspector/pkg/junit"
	"test-inspector/pkg/models"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "inspect test results comparing to the reference run for your project",
	Run: func(cmd *cobra.Command, args []string) {
		if err := validateVersionID(); err != nil {
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

		templates, err := supa.GetTemplate(int64(versionID))
		if err != nil {
			fmt.Printf("error trying to retrieve template test results: %v", err)
			return
		}

		errors := 0
		warns := 0
		fmt.Print(color.Blue + "Test Results comparison report:\n" + color.Reset)

		fmt.Printf("\n%s%d%s test results found in local run (%s)\n",
			color.Blue, len(results), color.Reset, resultsPath)
		fmt.Printf("%s%d%s test results found in template run\n",
			color.Green, len(templates), color.Reset)
		if len(results) < len(templates) {
			fmt.Printf("\n%sWARNING%s: number of test results in local run (%d) is less "+
				"then number of test results in template run (%d)\n",
				color.Yellow, color.Reset, len(results), len(templates))
		} else if len(results) > len(templates) {
			fmt.Printf("\n%sWARNING%s: number of test results in local run (%d) differs "+
				"from number of test results in template run (%d)\n",
				color.Yellow, color.Reset, len(results), len(templates))
		}
		fmt.Print("\n")

		var wg sync.WaitGroup
		var mu sync.Mutex
		for _, t := range templates {
			t := t
			wg.Add(1)
			go func() {
				defer wg.Done()
				report := ""
				r := findSameResult(t, results)
				if r == nil {
					report += fmt.Sprintf(
						"%s[ERROR]%s: no test result found for template: %s - %s\n",
						color.Red, color.Reset, t.Name, t.ParentSuite)
					mu.Lock()
					errors++
					fmt.Print(report)
					mu.Unlock()
					return
				}
				if t.Status == "passed" && r.Status == "passed" &&
					t.Steps != "" && r.Steps != "" && t.Steps != r.Steps {
					var templateSteps []*models.StepContainer
					err = json.Unmarshal([]byte(t.Steps), &templateSteps)
					if err != nil {
						mu.Lock()
						fmt.Printf("error when unmarshal template steps: %+v", err)
						fmt.Print(report)
						mu.Unlock()
						return
					}
					var resultSteps []*models.StepContainer
					err = json.Unmarshal([]byte(r.Steps), &resultSteps)
					if err != nil {
						mu.Lock()
						fmt.Printf("error when unmarshal result steps: %+v", err)
						fmt.Print(report)
						mu.Unlock()
						return
					}
					stepsComp := compareSteps(templateSteps, resultSteps, t.Name, t.Name)
					if stepsComp != "" {
						report += fmt.Sprint(stepsComp)
						mu.Lock()
						fmt.Print(report)
						warns++
						mu.Unlock()
						return
					}
					mu.Lock()
					fmt.Print(report)
					mu.Unlock()
					return
				}
			}()
		}
		wg.Wait()

		fmt.Printf("\n%s%d errors%s and %s%d warnings%s found\n",
			color.Red, errors, color.Reset,
			color.Yellow, warns, color.Reset)
		if errors == 0 && warns == 0 {
			fmt.Print(color.Green + "All checks passed!\n" + color.Reset)
		}
		os.Exit(errors)
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inspectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inspectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func findSameResult(
	template models.SupaResult,
	results map[uuid.UUID]models.SupaResult) *models.SupaResult {
	for _, r := range results {
		if normalizeName(r.Name) == normalizeName(template.Name) && checkSuiteNames(template, r) {
			return &r
		}
	}
	return nil
}

// nolint:gocyclo // this is just trying to find a match for suite/subsuite/parentsuite
func checkSuiteNames(template, result models.SupaResult) bool {
	if (normalizeName(result.Suite) == normalizeName(template.Suite)) && (template.Suite != "") ||
		(normalizeName(result.ParentSuite) == normalizeName(template.Suite)) && (template.Suite != "") ||
		(normalizeName(result.SubSuite) == normalizeName(template.Suite)) && (template.Suite != "") ||
		(normalizeName(result.Suite) == normalizeName(template.ParentSuite)) && (template.ParentSuite != "") ||
		(normalizeName(result.ParentSuite) == normalizeName(template.ParentSuite)) && (template.ParentSuite != "") ||
		(normalizeName(result.SubSuite) == normalizeName(template.ParentSuite)) && (template.ParentSuite != "") ||
		(normalizeName(result.Suite) == normalizeName(template.SubSuite)) && (template.SubSuite != "") ||
		(normalizeName(result.ParentSuite) == normalizeName(template.SubSuite)) && (template.SubSuite != "") ||
		(normalizeName(result.SubSuite) == normalizeName(template.SubSuite)) && (template.SubSuite != "") {
		return true
	}
	return false
}

func compareSteps(templateSteps, resultSteps []*models.StepContainer, parent, test string) string {
	if len(templateSteps) != len(resultSteps) {
		return fmt.Sprintf("%s[WARN]%s: number of steps in template - %s - (%d) is not equal to "+
			"number of steps in result (%d) for parent: %s\n",
			color.Yellow, color.Reset, test, len(templateSteps), len(resultSteps), parent)
	}
	for i, t := range templateSteps {
		if replaceAllSubstringsInBrackets(t.Name) != replaceAllSubstringsInBrackets(resultSteps[i].Name) {
			return fmt.Sprintf("%s[WARN]%s: step name in template - %s - (%s) is not equal to "+
				"step name in result (%s) for parent: %s, pos: %d\n",
				color.Yellow, color.Reset, test, t.Name, resultSteps[i].Name, parent, t.Position)
		}
		if t.StepContainer != nil && len(t.StepContainer) != 0 {
			if resultSteps[i].StepContainer != nil && len(resultSteps[i].StepContainer) != 0 {
				if compRes := compareSteps(t.StepContainer, resultSteps[i].StepContainer, t.Name, test); compRes != "" {
					return compRes
				}
			} else {
				return fmt.Sprintf("%s[WARN]%s: inner steps exist in template - %s - (%s) but not "+
					"in result (%s) for parent: %s, pos: %d\n",
					color.Yellow, color.Reset, test, t.Name, resultSteps[i].Name, parent, t.Position)
			}
		}
	}
	return ""
}

func replaceAllSubstringsInBrackets(str string) string {
	re := regexp.MustCompile(`\{.*\}`)
	return re.ReplaceAllString(str, "")
}

func normalizeName(str string) string {
	re := regexp.MustCompile(`\s|_`)
	s := re.ReplaceAllString(str, "")
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "test", "")
	return s
}
