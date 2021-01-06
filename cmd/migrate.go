package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ionous/errutil"
	jp "github.com/ionous/json-patch"
)

func main() {
	// fix: remove
	paths := "/Users/ionous/Dev/go/src/github.com/ionous/iffy/stories"
	patchPath := "/Users/ionous/Dev/go/src/github.com/ionous/iffy/cmd/migrate/push.patch.js"

	// fix: default to using stdin, stdout.
	flag.StringVar(&paths, "in", paths, "comma separated input files or directory names")
	flag.StringVar(&patchPath, "patch", patchPath, "patch file")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.Parse()
	//
	var patch jp.Patch
	if e := readJson(patchPath, &patch); e != nil {
		panic(e)
	} else if e := migratePaths(paths, patch); e != nil {
		panic(e)
	}
}

func migratePaths(paths string, patch jp.Migration) (err error) {
	return readPaths(paths, func(path string, doc interface{}) (err error) {
		log.Printf("migrating %q...", path)
		if cnt, e := patch.Migrate(doc); e != nil {
			err = e
		} else if cnt == 0 {
			log.Println("unchanged.")
		} else {
			if f, e := os.Create(path); e != nil {
				err = e
			} else {
				defer f.Close()
				js := json.NewEncoder(f)
				js.SetIndent("", "  ")
				err = js.Encode(doc)
				log.Println("migrated.")
			}
		}
		return //
	})
}

// read a comma-separated list of files and directories
func readPaths(filePaths string, cb func(path string, data interface{}) error) (err error) {
	split := strings.Split(filePaths, ",")
	for _, path := range split {
		if info, e := os.Stat(path); e != nil {
			err = e
		} else {
			if !info.IsDir() {
				var one interface{}
				if e := readJson(path, &one); e != nil {
					err = e
					break
				} else if e := cb(path, one); e != nil {
					err = e
					break
				}
			} else {
				if !strings.HasSuffix(path, "/") {
					path += "/" // trailing slash needed for opening symbolic directories on macos
				}
				// walk files...
				if e := filepath.Walk(path, func(path string, info os.FileInfo, e error) (err error) {
					if e != nil {
						err = e
					} else if !info.IsDir() && filepath.Ext(path) == ".if" {
						var one interface{}
						if e := readJson(path, &one); e != nil {
							err = e
						} else if e := cb(path, one); e != nil {
							err = e
						}
					}
					return // walk
				}); e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}

func readJson(path string, out interface{}) (err error) {
	if f, e := os.Open(path); e != nil {
		err = errutil.Fmt("can't open %q because %v", path, e)
	} else {
		defer f.Close()
		if e := json.NewDecoder(f).Decode(out); e != nil && e != io.EOF {
			err = errutil.Fmt("can't decode %q because %v", path, e)
		}
	}
	return
}
