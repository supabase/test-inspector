package supatms

import "test-inspector/pkg/models"

// ToResult takes an Allure result and a launch ID and converts to a SupaResult
func ToResult(launchID int64, r models.AllureResult, steps string) models.SupaResult {
	if r.FullName == nil {
		r.FullName = stringRef("")
	}
	suite, _ := r.FindLabel("suite")
	parentSuite, _ := r.FindLabel("parentSuite")
	subSuite, _ := r.FindLabel("subSuite")
	feature, _ := r.FindLabel("feature")
	res := models.SupaResult{
		ID:          r.UUID,
		Name:        r.Name,
		FullName:    *r.FullName,
		Suite:       suite,
		ParentSuite: parentSuite,
		SubSuite:    subSuite,
		Feature:     feature,
		Status:      r.Status,
		Description: r.Description,
		LaunchID:    launchID,
		Duration:    int32(r.Stop - r.Start),
		Steps:       steps,
	}

	return res
}

func stringRef(s string) *string {
	return &s
}
