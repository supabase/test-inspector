// nolint:revive // this is just a table columns package
package version

// Version is columns of the version table.
type Version int

const (
	ID Version = iota
	CreatedAt
	Repo
	ProjectID
	VersionName
)

var versionCols = [...]string{
	"id",
	"created_at",
	"repo",
	"project_id",
	"version_name",
}

func (s Version) String() string {
	if ID <= s && s <= VersionName {
		return versionCols[s]
	}
	return ""
}
