package models

import (
	"github.com/google/uuid"
)

// Project is the single project that is being inspected
//
// @property {int64} ID - The ID of the project.
// @property {string} Name - The name of the project.
// @property Description - A description of the project.
type Project struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// Version is the version of the project (for ex. supabase-js for project supabase client libs).
//
// @property ID - The ID of the version.
// @property Repo - The name of the repository that the version is in.
// @property {int32} ProjectID - The ID of the project that this version belongs to.
// @property {string} VersionName - The name of the version.
type Version struct {
	ID          *int64  `json:"id,omitempty"`
	Repo        *string `json:"repo,omitempty"`
	ProjectID   int32   `json:"project_id"`
	VersionName string  `json:"version_name"`
}

// Launch represents a single test launch.
//
// @property ID - The ID of the launch.
// @property {bool} IsTemplate - This is a boolean value that indicates whether the launch is a
// template or not.
// @property Origin - The origin of the launch. Can be one of the following:
// @property UserID - The ID of the user who launched the test.
// @property Duration - The duration of the launch in milliseconds.
// @property {string} Name - The name of the launch.
// @property {int64} VersionID - The ID of the version of the test that you want to launch.
type Launch struct {
	ID         *int64  `json:"id,omitempty"`
	IsTemplate bool    `json:"is_template"`
	Origin     *string `json:"origin,omitempty"`
	UserID     *string `json:"user_id,omitempty"`
	Duration   *int32  `json:"duration,omitempty"`
	Name       string  `json:"name"`
	VersionID  int64   `json:"version_id"`
}

// SupaResult is a result of a test DTO for supabase project.
// @property ID - The unique identifier for the test result.
// @property {string} Name - The name of the test
// @property {string} FullName - The full name of the test.
// @property {string} Suite - The name of the suite that the test belongs to.
// @property {string} ParentSuite - The name of the parent suite.
// @property {string} SubSuite - The name of the sub-suite that the test belongs to.
// @property {string} Feature - The name of the feature file
// @property {string} Status - The status of the test. This can be one of the following:
// @property Description - The description of the test case.
// @property {int64} LaunchID - The ID of the launch that this test belongs to.
// @property {int32} Duration - The duration of the test in milliseconds
// @property {string} Steps - This is a JSON string that contains the steps of the test.
// @property {[]*StepContainer} Stps - This is a slice of StepContainer structs.
type SupaResult struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	FullName    string    `json:"fullname"`
	Suite       string    `json:"suite"`
	ParentSuite string    `json:"parent_suite"`
	SubSuite    string    `json:"sub_suite"`
	Feature     string    `json:"feature"`
	Status      string    `json:"status"`
	Description *string   `json:"description,omitempty"`
	LaunchID    int64     `json:"launch_id"`
	Duration    int32     `json:"duration"`
	Steps       string    `json:"steps"`

	Stps []*StepContainer `json:"-"`
}

// StepContainer is a struct that is used to unmarshal the steps field of a SupaResult
// @property {[]*StepContainer} StepContainer - This is a slice of StepContainer objects.
// @property {string} Name - The name of the step
// @property {string} Status - The status of the step.
// @property {int16} Position - The position of the step in the workflow.
type StepContainer struct {
	StepContainer []*StepContainer `json:"steps,omitempty"`
	Name          string           `json:"name"`
	Status        string           `json:"status"`
	Position      int16            `json:"position"`
}
