# json-patch
[![Go Reference](https://pkg.go.dev/badge/github.com/ionous/json-patch.svg)](https://pkg.go.dev/github.com/ionous/json-patch)
[![Go Report Card](https://goreportcard.com/badge/github.com/ionous/json-patch)](https://goreportcard.com/report/github.com/ionous/json-patch)

This is a very modest attempt to bridge the gap between JSON Path and JSON Patch.

## Rationale

The JSON Patch ( [RFC-6902](https://tools.ietf.org/html/rfc6902) ) way of patching uses JSON Pointer ( [RFC-6901](https://tools.ietf.org/html/rfc6901) ) to refer to specific predetermined spots within a document. While that's great for some kinds of documents, transforming recursively defined data seems outside of its scope. And while the [JSON Path](https://goessner.net/articles/JsonPath/) concept can be used with recursive data, there don't seem to be any off-the-shelf tools for patching data using its paths.

There's no perfect standard for JSON Paths. This uses [PaesslerAG's](https://github.com/PaesslerAG/). There are [others](https://cburgmer.github.io/json-path-comparison/).

## Differences b/t the RFC and this lib.

* JSON Paths affect multiple values; so each operation affects multiple nodes.
* RFC compliant error handling hasn't been evaluated. 
* Array handling hasn't been explored deeply.
* The operation "test" in the RFC requires a value, here it does not. Also, "test" here supports arrays of recursive "patches" and "subpatches", both of which are processed should the test succeed. The "patches" applies its operations to the current document. The "subpatches" applies its operations to each object matched by the test.
* Defines a "reason" key as a way to add comments to patch files.
* Paths can (optionally) be specified as a pair of values -- a `parent path` targeting one or more json objects, and a `child field` string within each matching object. This is required for some case where paths use filters ( #1 )


While it's not necessary to have this gracefully decay into the RFC behavior ( ie. so that it can directly support json patch documents ) that would be cool. ( ex. replace the path pair with a single string value; transform json pointers starting with "/" into their "$" path equivalents; directly support all verbs; correctly handle errors; ... )  

It also might be nice to merge the 'test' op's "patches" and "subpatches". Each path could itself could indicate the desired matching behavior. For instance, paths starting with the standard `$` could match against the current document, while `@` could match against the object matched by the test.

## Sample Patch

Here's a sample(*) migration file. Each command can alter multiple nodes.

```javascript
[{
    "op": "test",
    "reason": "select the nodes containing type #choose",
    "path": "$..[?(@.type=='choose')]",
    "subpatches": [{
      "op": "test",
      "reason": "change #false to #else *if* false isnt empty",
      "path": "$.value['$FALSE'].value['$EXE'].*",
      "patches": [{
        "op": "move",
        "reason": "notes 'patches' applies to the current doc, not the elements selected by the test",
        "from": "$.value['$FALSE']",
        "path": "$.value['$ELSE']"
      }]
    }, {
      "op": "remove",
      "reason": "remove #false block in case it wasnt moved in the test",
      "path": "$.value['$FALSE']"
    }, {
      "op": "replace",
      "reason": "#do is required; #true was optional, so first create a blank.",
      "path": "$.value['$DO']",
      "value":{
        "type": "activity",
        "value": {
          "$EXE": []
        }
      }
    }, {
      "op": "copy",
      "reason": "now... copy the #true actions (if they exist).",
      "from": "$.value['$TRUE'].value['$EXE']",
      "path": "$.value['$DO'].value['$EXE']"
    }, {
      "op": "remove",
      "reason": "remove #true ( if it existed. )",
      "path": "$.value['$TRUE']"
    },{
      "op": "replace",
      "reason": "finally, rename #choose to #choose_action",
      "path": "$.type",
      "value": "choose_action"
    }]
}]
```

(*) _Doesn't really count as a full example without before and after data i guess._
