package models

import (
	"strconv"

	"github.com/google/uuid"
)

// AllureStepContainer is an interface that has a method `GetSteps()` that returns a slice of `Step`
// pointers. Helper to treat Step and Test as the same type.
// @property {[]*Step} GetSteps - Returns a slice of steps.
type AllureStepContainer interface {
	GetSteps() []*Step
}

// AllureResult is a struct that holds the result information of a test.
// @property {string} Name - The name of the test case.
// @property {string} Status - The status of the test case. Possible values are: PASSED, FAILED,
// BROKEN, SKIPPED, PENDING, CANCELED, UNKNOWN.
// @property Stage - The stage of the test.
// @property StatusDetails - This is a struct that contains the message and trace properties.
// @property {[]*Step} Steps - An array of steps that are part of the test case.
// @property {[]*Attachment} Attachments - A list of attachments.
// @property {[]*Parameter} Parameters - A list of parameters that were used in the test.
// @property {int64} Start - The start time of the test in milliseconds since the epoch.
// @property {int64} Stop - The time when the test finished.
// @property UUID - A unique identifier for the test case.
// @property HistoryID - The ID of the test run.
// @property FullName - The full name of the test case.
// @property Description - The description of the test case.
// @property DescriptionHTML - The description of the test case in HTML format.
// @property {[]*Label} Labels - Labels are used to group tests. For example, you can group tests by
// severity, by feature, by component, etc.
// @property {[]*Link} Links - Links to other results.
type AllureResult struct {
	Name            string         `json:"name"`
	Status          string         `json:"status"`
	Stage           *string        `json:"stage,omitempty"`
	StatusDetails   *StatusDetails `json:"statusDetails,omitempty"`
	Steps           []*Step        `json:"steps,omitempty"`
	Attachments     []*Attachment  `json:"attachments,omitempty"`
	Parameters      []*Parameter   `json:"parameters,omitempty"`
	Start           int64          `json:"start"`
	Stop            int64          `json:"stop"`
	UUID            uuid.UUID      `json:"uuid"`
	HistoryID       *string        `json:"historyId,omitempty"`
	FullName        *string        `json:"fullName,omitempty"`
	Description     *string        `json:"description,omitempty"`
	DescriptionHTML *string        `json:"descriptionHtml,omitempty"`
	Labels          []*Label       `json:"labels,omitempty"`
	Links           []*Link        `json:"links,omitempty"`
}

// GetSteps returns the steps of a test.
func (r *AllureResult) GetSteps() []*Step {
	return r.Steps
}

// FindLabel trying to find the label by name.
func (r *AllureResult) FindLabel(name string) (string, bool) {
	for _, l := range r.Labels {
		if l.Name == name {
			return l.Value, true
		}
	}
	return "", false
}

// FindLinkByType trying to find the link in the labels.
func (r *AllureResult) FindLinkByType(t string) (string, bool) {
	for _, l := range r.Links {
		if *l.Type == t {
			return l.URL, true
		}
	}
	return "", false
}

// FindID trying to find the test id in the labels.
func (r *AllureResult) FindID() (int32, bool) {
	var id string
	var ok bool
	if id, ok = r.FindLabel("testId"); !ok {
		if id, ok = r.FindLinkByType("supatms"); !ok {
			id, ok = r.FindLinkByType("tms")
		}
	}
	if !ok {
		return 0, ok
	}
	// nolint:gosec // no way this will overflow
	intID, err := strconv.Atoi(id)
	if err != nil {
		return 0, false
	}
	return int32(intID), true
}

// Container is a named collection of tests or suites.
// @property UUID - The UUID of the container.
// @property {string} Name - The name of the container.
// @property {[]uuid.UUID} Children - An array of UUIDs of the children of this container.
// @property Description - A description of the container.
// @property DescriptionHTML - The description of the container, in HTML format.
// @property {[]*Step} Befores - A list of steps that should be run before the container is run.
// @property {[]*Step} Afters - A list of steps that will be executed after the container is finished.
// @property {[]*Label} Labels - A list of labels that can be used to filter the container.
// @property {[]*Link} Links - A list of links to external resources.
// @property {int64} Start - The start time of the container in milliseconds since the epoch.
// @property {int64} Stop - The time the container finished running.
type Container struct {
	UUID            uuid.UUID   `json:"uuid"`
	Name            string      `json:"name"`
	Children        []uuid.UUID `json:"children"`
	Description     *string     `json:"description,omitempty"`
	DescriptionHTML *string     `json:"descriptionHtml,omitempty"`
	Befores         []*Step     `json:"befores,omitempty"`
	Afters          []*Step     `json:"afters,omitempty"`
	Labels          []*Label    `json:"labels,omitempty"`
	Links           []*Link     `json:"links,omitempty"`
	Start           int64       `json:"start"`
	Stop            int64       `json:"stop"`
}

// FindLabel function that returns a label value by label key.
func (c *Container) FindLabel(name string) (string, bool) {
	for _, l := range c.Labels {
		if l.Name == name {
			return l.Value, true
		}
	}
	return "", false
}

// Suite is a container for other suites or tests.
// @property UUID - A unique identifier for the suite.
// @property Container - The container that the suite is running in.
// @property {[]*Suite} Suites - A list of suites that are children of this suite.
// @property {[]*AllureResult} Tests - The tests that are part of this suite.
// @property Parent - The parent suite of the current suite.
type Suite struct {
	UUID      uuid.UUID       `json:"uuid"`
	Container *Container      `json:"container"`
	Suites    []*Suite        `json:"suites,omitempty"`
	Tests     []*AllureResult `json:"tests,omitempty"`
	Parent    *Suite          `json:"-"`
}

// Step is a single step in a test.
// @property {string} Name - The name of the step.
// @property Description - A description of the step.
// @property Status - The status of the step. Possible values are:
// @property StatusDetails - This is a nested object that contains the following properties:
// @property Stage - The stage of the step.
// @property {[]*Step} Steps - A list of steps that are part of this step.
// @property {[]*Attachment} Attachments - An array of attachments that are associated with the step.
// @property {[]*Parameter} Parameters - A list of parameters that are used by the step.
// @property {int64} Start - The time the step started.
// @property {int64} Stop - The time the step stopped.
type Step struct {
	Name          string         `json:"name"`
	Description   *string        `json:"description,omitempty"`
	Status        *string        `json:"status,omitempty"`
	StatusDetails *StatusDetails `json:"statusDetails,omitempty"`
	Stage         *string        `json:"stage,omitempty"`
	Steps         []*Step        `json:"steps,omitempty"`
	Attachments   []*Attachment  `json:"attachments,omitempty"`
	Parameters    []*Parameter   `json:"parameters,omitempty"`
	Start         int64          `json:"start"`
	Stop          int64          `json:"stop"`
}

// GetSteps returns the steps of a step or test.
func (s *Step) GetSteps() []*Step {
	return s.Steps
}

// Attachment is a struct that holds info about an attachment to the step
// @property {string} Name - The name of the attachment.
// @property {string} Source - The URL of the attachment.
// @property {string} Type - The type of attachment.
type Attachment struct {
	Name   string `json:"name"`
	Source string `json:"source"`
	Type   string `json:"type"`
}

// A Label is a name and value pair for labels.
// @property {string} Name - The name of the label.
// @property {string} Value - The value of the label.
type Label struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// A Link is link object that can be attached to a test or test suite to connect with TMS.
// @property {string} Name - The name of the link.
// @property {string} URL - The URL of the link.
// @property Type - The type of the link. This is an optional property.
type Link struct {
	Name string  `json:"name"`
	URL  string  `json:"url"`
	Type *string `json:"type,omitempty"`
}

// A Parameter is a name/value pair of test parameters.
// @property {string} Name - The name of the parameter.
// @property {string} Value - The value of the parameter.
type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// StatusDetails is a struct with fields `Known`, `Muted`, `Flaky`, `Message`, and `Trace`.
//
// @property {bool} Known - Whether the test is known to be failing.
// @property {bool} Muted - If the test is muted, this will be true.
// @property {bool} Flaky - If the test is flaky, this will be true.
// @property Message - A message of the error or skip.
// @property Trace - The stack trace of the failure.
type StatusDetails struct {
	Known   bool    `json:"known"`
	Muted   bool    `json:"muted"`
	Flaky   bool    `json:"flaky"`
	Message *string `json:"message,omitempty"`
	Trace   *string `json:"trace,omitempty"`
}

// "labels": [
// 		{
// 				"name": "package",
// 				"value": "com.example.project.rest.package.TestClass"
// 		},
// 		{
// 				"name": "testClass",
// 				"value": "com.example.project.rest.package.TestClass"
// 		},
// 		{
// 				"name": "testMethod",
// 				"value": "testMethodName"
// 		},
// 		{
// 				"name": "parentSuite",
// 				"value": "Root suite name"
// 		},
// 		{
// 				"name": "suite",
// 				"value": "Parent suite name"
// 		},
// 		{
// 				"name": "subSuite",
// 				"value": "com.example.project.rest.package.TestClass"
// 		},
// 		{
// 				"name": "host",
// 				"value": "NB-EROMANOV.local"
// 		},
// 		{
// 				"name": "thread",
// 				"value": "49476@NB-EROMANOV.local.TestNG-test=Notification tests-2(30)"
// 		},
// 		{
// 				"name": "framework",
// 				"value": "testng"
// 		},
// 		{
// 				"name": "language",
// 				"value": "java"
// 		},
// 		{
// 				"name": "feature",
// 				"value": "Test feature"
// 		},
// 		{
// 			"name": "suite",
// 			"value": "SuiteOrTestName"
// 		},
// 		{
// 			"name": "severity",
// 			"value": "blocker"
// 		},
// 		{
// 			"name": "tag",
// 			"value": "smoke"
// 		}
// ],
