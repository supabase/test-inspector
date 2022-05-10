// nolint:revive // this is just a table columns package
package label

// Label is a list of columns of the label table.
type Label int

const (
	ID Label = iota
	ProjectID
	Name
	Value
)

var labels = [...]string{
	"id",
	"project_id",
	"name",
	"value",
}

func (s Label) String() string {
	if ID <= s && s <= Value {
		return labels[s]
	}
	return ""
}
