package jsonpatch

import (
	"strings"
)

// Migration interface for patching json.
type Migration interface {
	Migrate(doc interface{}) (int, error)
}

// Copy replicates pieces of the document.
type Copy struct {
	From Target `json:"from"`
	To   Target `json:"to"`
}

// Replace
type Replace struct {
	From Target      `json:"from"`
	With interface{} `json:"with"`
}

//
type Target struct {
	Parent string `json:"parent"`
	Field  string `json:"field"`
}

// FIX: this is surely going to break something at some point
// Paessler's paths dont handle single quotes by default... but maybe there's a way to make it?
func (t Target) dequote() string {
	return strings.Replace(t.Parent, "'", `"`, -1)
}

// Migrate runs the copy command.
func (op *Copy) Migrate(doc interface{}) (ret int, err error) {
	// FIX: im sure you'll want a Move before long.
	// could be done with Copy and Replace(nil) --
	// to optimize, it might be easier if you resolved paths into objects first.
	// you could pass around a "Cursor" which holds
	// a parent node, original sub/path string, a dynamically cached collection of documents, and an error.
	return ReplicateValues(doc,
		op.From.dequote(), op.From.Field,
		op.To.dequote(), op.To.Field)
}

// Migrate runs the replace command.
func (op *Replace) Migrate(doc interface{}) (ret int, err error) {
	if e, ok := op.With.(error); ok {
		err = e
	} else {
		ret, err = ReplaceValues(doc, op.From.dequote(), op.From.Field, op.With)
	}
	return
}
