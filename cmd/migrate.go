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
	jsonpatch "github.com/ionous/json-patch"
)

func main() {
	// fix: remove
	var paths, patchPath, outdir string
	// fix: default to using stdin, stdout.
	flag.StringVar(&paths, "in", paths, "comma separated input files or directory names")
	flag.StringVar(&outdir, "out", outdir, "optional output directory. without this output overwrites input.")
	flag.StringVar(&patchPath, "patch", patchPath, "patch file")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.Parse()

	//
	var patch jp.Patches
	if e := readJson(patchPath, &patch); e != nil {
		panic(e)
	} else if e := migratePaths(paths, outdir, patch); e != nil {
		panic(e)
	}
}

func migratePaths(paths, outdir string, patch jp.Migration) (err error) {
	return readPaths(paths, func(dir, file string, doc interface{}) (err error) {
		var outpath string
		if len(outdir) == 0 {
			outpath = dir + file
		} else {
			if filepath.IsAbs(outdir) {
				outpath = filepath.Join(outdir, file)
			} else {
				outpath = filepath.Join(paths, outdir, file)
			}
		}
		path := dir + file
		log.Printf("reading %q...", path)
		if cnt, e := patch.Migrate(doc); e != nil {
			err = errutil.Fmt("migration of %q failed because %v", path, e)
		} else if cnt == 0 && len(outdir) == 0 {
			log.Println("unchanged.")
		} else if f, e := os.Create(outpath); e != nil {
			err = errutil.New("couldn't write to", path)
		} else {
			defer f.Close()
			js := json.NewEncoder(f)
			js.SetEscapeHTML(jsonpatch.EscapeHTML)
			js.SetIndent("", "  ")
			if e := js.Encode(doc); e != nil {
				err = errutil.New("couldnt encode output", e)
			} else if len(outdir) == 0 {
				log.Println("CHANGED.")
			} else {
				log.Printf("wrote: %q ( w/ %d changes ).", outpath, cnt)
			}
		}
		return
	})
}

// read a comma-separated list of files and directories
func readPaths(filePaths string, cb func(dir, file string, data interface{}) error) (err error) {
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
				} else {
					dir, file := filepath.Split(path)
					if e := cb(dir, file, one); e != nil {
						err = e
						break
					}
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
						} else {
							dir, file := filepath.Split(path)
							if e := cb(dir, file, one); e != nil {
								err = e
							}
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
