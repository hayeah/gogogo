
Build commands from a JSON stream and run them, possibly in parallel.

In data.json:

```
{"foo": 1, "bar": "a string"}
{"foo": 2, "bar": "another string"}
```

The command is specified as a golang [text template](http://golang.org/pkg/text/template/):

```
cat data.json > gogo 'echo {{.foo}} {{.bar}}'
```