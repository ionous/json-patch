package jsonpatch

import (
	"encoding/json"

	"github.com/ionous/errutil"
)

// Migration interface for patching json.
type Migration interface {
	Migrate(doc interface{}) (int, error)
}

// Copy replicates pieces of a document.
type Copy struct {
	From Target `json:"from"`
	To   Target `json:"path"`
}

// Move relocates pieces of a document.
type Move struct {
	From Target `json:"from"`
	To   Target `json:"path"`
}

// Patches runs a series of other migrations.
type Patches []PatchCommand

type Remove struct {
	Path Target `json:"path"`
}

// Replace substitutes new values for pieces of a document.
type Replace struct {
	Path  Target          `json:"path"`
	Value json.RawMessage `json:"value"`
}

// Test validates pieces of a document, then possibly runs sub matches.
// Patches are run against the current document, and not specifically things matched by the test.
// SubPatches are run against all elements matched by the test.
type Test struct {
	Path       Target          `json:"path"`
	Value      json.RawMessage `json:"value,omitempty"`
	Patches    Patches         `json:"patches,omitempty"`
	SubPatches Patches         `json:"subpatches,omitempty"`
}

// Migrate runs the copy command.
func (op *Copy) Migrate(doc interface{}) (ret int, err error) {
	from, to := op.From.Parent.Select(doc), op.To.Parent.Select(doc)
	// get the values without deleting them
	if els, e := ExtractValues(from, op.From.Field, false); e != nil {
		err = e
	} else if len(els) > 0 {
		ret, err = InsertValues(to, op.To.Field, els)
	}
	return
}

// Migrate copies to the new location then removes from the old.
func (op *Move) Migrate(doc interface{}) (ret int, err error) {
	from, to := op.From.Parent.Select(doc), op.To.Parent.Select(doc)
	// delete the old bits first; they'll be in memory still to move them.
	if els, e := ExtractValues(from, op.From.Field, true); e != nil {
		err = e
	} else if len(els) > 0 {
		ret, err = InsertValues(to, op.To.Field, els)
	}
	return
}

// Migrate runs a list of migrations.
func (ps Patches) Migrate(doc interface{}) (ret int, err error) {
	for i, op := range ps {
		if cnt, e := op.Migrate(doc); e != nil {
			err = errutil.Fmt("error encountered during migration command %T at %d, because %s", op, i, e)
			break
		} else {
			ret += cnt
		}
	}
	return
}

// ApplyOver runs this series of patches over all objects matched by the passed cursor.
func (ps Patches) ApplyOverMatches(cs Cursor) (ret int, err error) {
	if cnt, e := cs.Resolve(); e != nil {
		err = e
	} else {
		for i := 0; i < cnt; i++ {
			doc := cs.Element(i)
			if cnt, e := ps.Migrate(doc); e != nil {
				err = errutil.Append(err, e)
			} else {
				ret += cnt
			}
		}
	}
	return
}

// Migrate runs the replace command.
func (op *Remove) Migrate(doc interface{}) (ret int, err error) {
	from := op.Path.Parent.Select(doc)
	if els, e := ExtractValues(from, op.Path.Field, true); e != nil {
		err = e
	} else {
		ret = len(els)
	}
	return
}

// Migrate runs the replace command.
func (op *Replace) Migrate(doc interface{}) (ret int, err error) {
	from := op.Path.Parent.Select(doc)
	return ReplaceValues(from, op.Path.Field, op.Value)
}

// Migrate runs the test command; potentially recursive.
func (op *Test) Migrate(doc interface{}) (ret int, err error) {
	from := op.Path.Parent.Select(doc)
	if len(op.Value) > 0 {
		if matches, misses, e := CompareValues(from, op.Path.Field, op.Value); e != nil {
			err = e
		} else if misses > 0 {
			err = errutil.Fmt("test detected %d mismatches", misses)
		} else {
			ret += matches
		}
	}

	// fix: should handle field not being null
	if ps := op.Patches; err == nil && ps != nil {
		cnt, e := from.Resolve()
		if e != nil {
			err = e
		} else if cnt > 0 {
			ret, err = ps.Migrate(doc)
		}
	}
	if ps := op.SubPatches; err == nil && ps != nil {
		ret, err = ps.ApplyOverMatches(from)
	}
	return
}
