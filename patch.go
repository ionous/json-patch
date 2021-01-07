package jsonpatch

import (
	"encoding/json"

	"github.com/ionous/errutil"
)

// PatchCommand holds a single command for de/serializing a migration.
type PatchCommand struct {
	Migration // pointer to the command
}

// UnmarshalJSON creates concrete implementations of migrations.
func (c *PatchCommand) UnmarshalJSON(data []byte) (err error) {
	// we have to read the name before we can know how to read the particular command.
	var op struct {
		Name string `json:"op"`
	}
	if e := json.Unmarshal(data, &op); e != nil {
		err = errutil.New("couldnt unmarshal patch command", e)
	} else if newCmd, ok := opFactory[op.Name]; !ok {
		err = errutil.Fmt("unknown migration %q", op.Name)
	} else {
		err = c.unmarshal(data, newCmd())
	}
	return
}

var opFactory = map[string]func() Migration{
	"copy":    func() Migration { return new(Copy) },
	"move":    func() Migration { return new(Move) },
	"remove":  func() Migration { return new(Remove) },
	"replace": func() Migration { return new(Replace) },
	"test":    func() Migration { return new(Test) },
}

// shared code for parsing migrations
func (c *PatchCommand) unmarshal(msg []byte, op Migration) (err error) {
	if e := json.Unmarshal(msg, op); e != nil {
		err = errutil.New("couldnt read patch command because", e)
	} else {
		c.Migration = op
	}
	return
}
