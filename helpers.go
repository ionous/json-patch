package jsonpatch

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/ionous/errutil"
)

// Target a location in a json document
// with a json path pointing to (one or more) objects,
// and a field addressing a member of each objects.
type Target struct {
	Parent   Path   `json:"parent"`
	Field    string `json:"field"`
	FullPath Path   // original path, the join of parent and field.
}

// https://regoio.herokuapp.com/
var parts = regexp.MustCompile(`(\[.+?\]|\.[^\[.]*)`)
var pieces = regexp.MustCompile(`^\.(\w+)$|^\['(.*)'\]$`)

// / UnmarshalJSON creates concrete implementations of migrations.
func (c *Target) UnmarshalJSON(data []byte) (err error) {
	// we have to read the name before we can know how to read the particular command.
	var val interface{}
	if e := json.Unmarshal(data, &val); e != nil {
		err = errutil.New("couldnt unmarshal patch command", e)
	} else {
		switch path := val.(type) {
		default:
			err = errutil.New("unknown path", val)
		case string:
			// split the strings into parts ( either dotted fields, or bracketed ones )
			var field string
			parts := parts.FindAllString(path, -1)
			ogpath := path
			if end := len(parts) - 1; end >= 0 {
				pieces := pieces.FindStringSubmatch(parts[end])
				if len(pieces) == 3 {
					// rejoin all the front parts together
					path = "$" + strings.Join(parts[:end], "")
					// except for the last piece sans any leading dot.
					if unbracketed := pieces[2]; len(unbracketed) > 0 {
						field = unbracketed
					} else if undotted := pieces[1]; len(undotted) > 0 {
						field = undotted
					} else {
						panic("impossible?")
					}
				}
			}
			c.FullPath, c.Parent, c.Field = Path(ogpath), Path(path), field

		case map[string]interface{}:
			if p, ok := path["parent"].(string); !ok {
				err = errutil.New("couldnt find parent in", path)
			} else {
				f, _ := path["field"].(string) // its okay for some ops if the field is missing
				*c = At(p, f)
			}
		}
	}
	return
}

// Path provides a JSONPath ready string.
// See: https://goessner.net/articles/JsonPath/, and https://github.com/PaesslerAG/
type Path string

func (p Path) String() string { return string(p) }

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
			if estr := e.Error(); strings.Contains(estr, "unknown key") {
				err = UnknownKey{estr, c.path.String()}
			} else {
				err = errutil.New("error selecting", c.path, e)
			}
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
