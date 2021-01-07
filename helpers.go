package jsonpatch

import (
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/ionous/errutil"
)

// Target a location in a json document
// with a json path pointing to (one or more) objects,
// and a field addressing a member of each objects.
type Target struct {
	Parent Path   `json:"parent"`
	Field  string `json:"field"`
}

// Path provides a JSONPath ready string.
// See: https://goessner.net/articles/JsonPath/, and https://github.com/PaesslerAG/
type Path string

// Select prepares a collection of elements from a json doc pointed to by the path.
func (p Path) Select(doc interface{}) Cursor {
	return Cursor{els: doc, path: p, res: nil}
}

// Cursor reads (caches) a collection of objects from a json docs.
type Cursor struct {
	res  error       // nil, cached, or error
	path Path        // location of the cursor; mainly for debugging.
	els  interface{} // doc or elements depending
}

const cached = errutil.Error("github.com/ionous/json-patch/cached")

// Path that describes the collection of objects targeted by the cursor.
func (c *Cursor) Path() Path {
	return c.path
}

// Resolve reads (caches) the objects targeted by a json path and returns the number of matches.
func (c *Cursor) Resolve() (ret int, err error) {
	if els, e := c.resolve(); e != nil {
		err = e
	} else {
		ret = len(els)
	}
	return
}

// Element returns one of the targeted objects.
func (c *Cursor) Element(i int) (ret interface{}) {
	if els, e := c.resolve(); e != nil {
		panic(e)
	} else {
		ret = els[i]
	}
	return
}

// returns a collection of objects originally pointed to by a path.
func (c *Cursor) resolve() (ret []interface{}, err error) {
	switch e := c.res; e {
	default:
		err = e
	case cached:
		ret = c.els.([]interface{})
	case nil:
		// FIX: this is surely going to break something at some point
		// Paessler's paths dont handle single quotes by default... but maybe there's a way to make it?
		path := strings.Replace(string(c.path), "'", `"`, -1)
		if tgt, e := jsonpath.Get(path, c.els); e != nil {
			err = errutil.New("error selecting", c.path, e)
			c.res = err
		} else {
			els, ok := tgt.([]interface{})
			if !ok {
				els = []interface{}{tgt} // convert a single result to a uniform array
			}
			ret, c.els, c.res = els, els, cached
		}
	}
	return
}
