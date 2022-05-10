// nolint:revive // this is just a table columns package
package result

// Result is a list of columns of the result table.
type Result int

const (
	ID Result = iota
	Name
	FullName
	Suite
	ParentSuite
	SubSuite
	Feature
	Status
	Description
	Steps
	LaunchID
	Duration
	CreatedAt
)

var results = [...]string{
	"id",
	"name",
	"fullname",
	"suite",
	"parent_suite",
	"sub_suite",
	"feature",
	"status",
	"description",
	"steps",
	"launch_id",
	"duration",
	"created_at",
}

func (s Result) String() string {
	if ID <= s && s <= CreatedAt {
		return results[s]
	}
	return ""
}
