# json-patch

This is a very modest attempt to bridge the gap between JSON Path and JSON Patch.

## Rationale

The JSON Patch method ( [RFC-6902](https://tools.ietf.org/html/rfc6902) ) for patching json documents uses JSON Pointer ( [RFC-6901](https://tools.ietf.org/html/rfc6901) ) to refer to specific predetermined spots within a document. While that's great for some kinds of documents, transforming recursively defined data seems outside of its scope. And while [JSON Path](https://goessner.net/articles/JsonPath/) concept can be used with recursive data, there don't seem to be any off the shelf tools for patching data using its paths.

There's no perfect standard for JSON Paths. This uses [PaesslerAG's](https://github.com/PaesslerAG/). There are [others](https://cburgmer.github.io/json-path-comparison/).

## Status

Currently only "replace" and "copy" are supported. Array support hasn't been verified. A hierarchical select is still needed to apply commands to sub-blocks. Proper patch error handling hasn't been evaluted.

As to other JSON Patch verbs: 
* "add" is basically a replace ( plus or minus some array, error handling jazz );
* "remove" is supported with a `null` replace value;
* "replace" is replace;
* "move" can be implemented with a copy and replace null;
* "copy" is copy;
* "test" will be supported via select.

It's not necessary to have this gracefully decay into the rfc behavior ( ie. so that it can directly support json patch documents ) but that would be cool. ( ex. transform json pointers starting with "/" into their "$" path equivalents; directly support all verbs; correctly handle errors. )

## Sample Patch

Here's a sample(*) migration file. Each command can alter multiple nodes.

```javascript
[{
  "patch": "replace",
  "reason": "restructure/rename 'list' to 'into'. this adds the 'into'(s), we'll then copy bits of 'list' and remove it.",
  "migration": {
    "from": {
      "parent": "$..[?(@.type=='list_push')].value",
      "field": "$INTO"
    },
    "with": {
      "type": "list_target",
      "value": {
        "type": "into_rec_list",
        "value": {
          "$VAR": {
            "type": "variable_name"
}}}}}},
{
  "patch": "copy",
  "migration": {
    "from": {
      "parent": "$..[?(@.type=='list_push')].value['$LIST']",
      "field": "value"
    },
    "to": {
      "parent": "$..[?(@.type=='list_push')].value['$INTO']..['$VAR']",
      "field": "value"
}}},
{
  "patch": "replace",
  "reason": "finish by removing the list nodes",
  "migration": {
    "from": {
      "parent": "$..[?(@.type=='list_push')].value",
      "field": "$LIST"
    },
    "with": null
}},
{
  "patch": "replace",
  "reason": "renames all list_push commands to put_edge",
  "migration": {
    "from": {
      "parent": "$..[?(@.type=='list_push')]",
      "field": "type"
    },
    "with": "put_edge"
}}]
```

(*) _Doesn't really count as a full example without before and after data i guess._
