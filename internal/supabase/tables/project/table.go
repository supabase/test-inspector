// nolint:revive // this is just a table columns package
package project

// Project is a list of columns of the project table.
type Project int

const (
	ID Project = iota
	CreatedAt
	Name
	Description
)

var projects = [...]string{
	"id",
	"created_at",
	"name",
	"description",
}

func (s Project) String() string {
	if ID <= s && s <= Description {
		return projects[s]
	}
	return ""
}
