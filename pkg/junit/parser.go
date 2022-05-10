package junit

import (
	"strings"
	"test-inspector/pkg/models"
	"test-inspector/pkg/supatms"

	"github.com/google/uuid"
	"github.com/joshdk/go-junit"
)

// ReadResults reads the results from the junit report and returns them as a map of SupaResults
func ReadResults(resultsPath string) (map[uuid.UUID]models.SupaResult, error) {
	results, err := parseResults(resultsPath)
	if err != nil {
		return nil, err
	}

	suparesults := map[uuid.UUID]models.SupaResult{}
	for _, r := range results {
		supares := supatms.ToResult(0, r, r.Steps[0].Name)
		suparesults[supares.ID] = supares
	}

	return suparesults, nil
}

func parseResults(resultsPath string) ([]models.AllureResult, error) {
	var suites []junit.Suite
	var err error
	if strings.Contains(resultsPath, ".xml") {
		suites, err = junit.IngestFile(resultsPath)
	} else {
		suites, err = junit.IngestDir(resultsPath)
	}
	if err != nil {
		return nil, err
	}
	res := []models.AllureResult{}

	for _, suite := range suites {
		res = append(res, convertTests(suite, suite.Package)...)
	}
	return res, nil
}

func convertTests(suite junit.Suite, parentSuite string) []models.AllureResult {
	res := []models.AllureResult{}
	for _, test := range suite.Tests {
		msg := test.Message
		if test.Error != nil {
			msg = msg + "\n" + test.Error.Error()
		}
		ar := models.AllureResult{
			Name:   test.Name,
			Status: string(test.Status),
			StatusDetails: &models.StatusDetails{
				Known:   false,
				Muted:   false,
				Flaky:   false,
				Message: &msg,
				Trace:   &test.SystemErr,
			},
			Steps: []*models.Step{
				{
					Name: test.SystemOut,
				},
			},
			Start:    0,
			Stop:     test.Duration.Milliseconds(),
			UUID:     uuid.New(),
			FullName: &test.Classname,
			Labels: []*models.Label{
				{
					Name:  "suite",
					Value: suite.Name,
				},
				{
					Name:  "parentSuite",
					Value: parentSuite,
				},
			},
		}
		for k, v := range test.Properties {
			ar.Labels = append(ar.Labels, &models.Label{
				Name:  k,
				Value: v,
			})
		}
		res = append(res, ar)
	}
	for _, s := range suite.Suites {
		res = append(res, convertTests(s, suite.Name)...)
	}
	return res
}
