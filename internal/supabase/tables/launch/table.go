// nolint:revive // this is just a table columns package
package launch

// Launch is a list of columns of the launch table.
type Launch int

const (
	ID Launch = iota
	CreatedAt
	IsTemplate
	Origin
	UserID
	Duration
	Name
	VersionID
	UpdatedAt
)

var launches = [...]string{
	"id",
	"created_at",
	"is_template",
	"origin",
	"user_id",
	"duration",
	"name",
	"version_id",
	"updated_at",
}

func (s Launch) String() string {
	if ID <= s && s <= UpdatedAt {
		return launches[s]
	}
	return ""
}
