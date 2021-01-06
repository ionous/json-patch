# json-patch

This is a very modest attempt to bridge the gap between JSON Path and JSON Patch.

## Rationale:

The JSON Patch method ( [RFC-6902](https://tools.ietf.org/html/rfc6902) ) for patching json documents uses JSON Pointer ( [RFC-6901](https://tools.ietf.org/html/rfc6901) ) to refer to specific predetermined spots within a document. While that's great for some kinds of documents, transforming recursively defined data seems outside of its scope. And while [JSON Path](https://goessner.net/articles/JsonPath/) concept can be used with recursive data, there don't seem to be any off the shelf tools for patching data using its paths.

There's no perfect standard for JSON Paths. This uses [PaesslerAG's](https://github.com/PaesslerAG/)] uses. There are [others](https://cburgmer.github.io/json-path-comparison/).
