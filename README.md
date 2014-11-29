`gogo` takes a JSON stream as input. For each object in the stream, it expands a template, and run the result.

By default, the commands are run serially. With the `-c` flag, you can run the commands in parallel.

# Installation

```
go get github.com/hayeah/gogo
```

# Example

Given a stream of json objects in data.json:

```
{"foo": 1, "bar": "a string"}
{"foo": 2, "bar": "another string"}
```

We can expand each of these with a Golang [text template](http://golang.org/pkg/text/template/):

```
cat data.json | gogo 'echo foo = {{.foo}}, bar = "{{.bar}}"'
```

And it produces the output:

```
2014/11/28 22:07:44 run cmd: echo foo = 1, bar = "a string"
foo = 1, bar = a string
2014/11/28 22:07:44 run cmd: echo foo = 2, bar = "another string"
foo = 2, bar = another string
```

# Changing the Command Evaluator

By default the commands are evaluated with the program invokation `sh -c`. You can change it to something else using the `-e` flag. To use Ruby as the evaluator:

```
cat data.json | gogo -e "ruby -e" 'puts "hello {{.foo}} {{.bar}}"'
```

# Run Processes Concurrently

Use the `-c` flag.

```
gogo -c 3 Template
```