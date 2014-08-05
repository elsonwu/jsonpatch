jsonpatch
=========

To handle json PATCH request 

###Json request:

```
[
    { "op": "remove", "path": "/a/b/c" },
    { "op": "add", "path": "/a/b/c", "value": [ "foo", "bar" ] },
    { "op": "replace", "path": "/a/b/c", "value": 42 }
]
```

Now only support op=`remove`, `replace` and `add`

### For Map specially

If you need to update a map, you can't pass the key in "path", you need to pass a sub-map, such as:

A document:

```
{
    "a": {
        "b": "hello"
    }
}
```

Willing to insert a `c: "world"` inside `a`, then you should pass the json like this:

```
[
	{ "op": "add", "path": "/a/b", "value": {"c": "world" } }
]
```

Below json will not work:

```
[
	{ "op": "add", "path": "/a/b/c", "value": "world" }
]
```

### Example

Please read `test/main.go`

### More

Please read this standard

[https://tools.ietf.org/html/rfc6902](https://tools.ietf.org/html/rfc6902)