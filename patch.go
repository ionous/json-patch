package jsonpatch

import (
	"encoding/json"

	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

// Patch holds a list of commands for de/serializing migrations.
type Patch []PatchCommand

// PatchCommand holds a single command for de/serializing a migration.
type PatchCommand struct {
	Name      string             `json:"patch"` // lowercase name of the migration struct. ex. "replace"
	Migration `json:"migration"` // pointer to the command
}

// Patch runs a list of migrations.
func (p Patch) Migrate(doc interface{}) (ret int, err error) {
	for i, op := range p {
		if cnt, e := op.Migrate(doc); e != nil {
			err = errutil.Fmt("error %v @%d=%v", e, i, pretty.Sprint(op))
			break
		} else {
			ret += cnt
		}
	}
	return
}

// UnmarshalJSON creates concrete implementations of migrations.
func (c *PatchCommand) UnmarshalJSON(data []byte) (err error) {
	// we have to read the name before we can know how to read the particular command.
	var rep struct {
		Name      string          `json:"patch"`
		Migration json.RawMessage `json:"migration"`
	}
	if e := json.Unmarshal(data, &rep); e != nil {
		err = e
	} else {
		c.Name = rep.Name // keep for debugging, serialization, etc.
		switch n, m := rep.Name, rep.Migration; n {
		case "replace":
			err = c.unmarshal(m, &Replace{})
		case "copy":
			err = c.unmarshal(m, &Copy{})
		default:
			err = errutil.New("unknown migration", n)
		}
	}
	return
}

// shared code for parsing migrations
func (c *PatchCommand) unmarshal(msg json.RawMessage, op Migration) (err error) {
	if e := json.Unmarshal(msg, op); e != nil {
		err = e
	} else {
		c.Migration = op
	}
	return
}
