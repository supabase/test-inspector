package allure

import (
	"encoding/json"
	"fmt"

	// nolint:depguard // we need to import this package to read result files
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"test-inspector/pkg/models"
	"test-inspector/pkg/supatms"

	"github.com/google/uuid"
)

// ReadResults reads all files from the allure results folder,
// parses them and returns a map of SupaResults
func ReadResults(resultsPath string) (map[uuid.UUID]models.SupaResult, error) {
	files, err := ioutil.ReadDir(resultsPath)
	if err != nil {
		return nil, fmt.Errorf("error trying to read files from dir: %v", err)
	}

	results := map[uuid.UUID]models.SupaResult{}
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, f := range files {
		f := f
		wg.Add(1)
		go func() {
			defer wg.Done()
			switch {
			case f.IsDir():
				// allure does not create inner folders
				return
			case strings.Contains(f.Name(), attachment.String()):
				// we don't need to upload attachments
				return
			case strings.Contains(f.Name(), container.String()):
				// suites provided by allure container json are not supported
				return
			case strings.Contains(f.Name(), result.String()):
				{
					res, err := parseResult(resultsPath, f)
					if err != nil {
						fmt.Printf("error trying to parse result %s: %v", f.Name(), err)
						return
					}
					steps := ""
					if res.Steps != nil && len(res.Steps) > 0 {
						stepsStruct := parseSteps(res)
						stepsRaw, err := json.Marshal(stepsStruct)
						if err != nil {
							fmt.Printf("problems with parsing steps: %s name - %s. %v", res.UUID, res.Name, err)
						}
						steps = string(stepsRaw)
					}
					result := supatms.ToResult(0, *res, steps)
					mu.Lock()
					results[result.ID] = result
					mu.Unlock()
				}
			}
		}()
	}
	wg.Wait()
	return results, nil
}

func parseResult(resultsPath string, f os.FileInfo) (*models.AllureResult, error) {
	var res models.AllureResult
	jsonFile, err := os.ReadFile(filepath.Join(resultsPath, f.Name()))
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(jsonFile, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func parseSteps(r models.AllureStepContainer) []*models.StepContainer {
	steps := []*models.StepContainer{}
	var ctr int16
	for _, s := range r.GetSteps() {
		if s.Status == nil {
			s.Status = stringRef("undefined")
		}
		stepInfo := models.StepContainer{
			StepContainer: []*models.StepContainer{},
			Name:          s.Name,
			Status:        *s.Status,
			Position:      ctr,
		}

		stepInfo.StepContainer = parseSteps(s)

		steps = append(steps, &stepInfo)
		ctr++
	}
	return steps
}

func stringRef(s string) *string {
	return &s
}

// Suffix is used to get the type of allure result file
type Suffix int

const (
	result Suffix = iota
	container
	attachment
)

var suffixes = [...]string{
	"result",
	"container",
	"attachment",
}

func (s Suffix) String() string {
	if result <= s && s <= attachment {
		return suffixes[s]
	}
	return ""
}
