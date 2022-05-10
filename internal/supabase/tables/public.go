// nolint:revive // this is just a list of tables
package tables

// Table is a list of tables.
type Table int

const (
	Labels Table = iota
	Launches
	Projects
	Results
	Versions
)

var tables = [...]string{
	"labels",
	"launches",
	"projects",
	"results",
	"versions",
}

func (s Table) String() string {
	if Labels <= s && s <= Versions {
		return tables[s]
	}
	return ""
}
