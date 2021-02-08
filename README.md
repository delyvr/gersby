# Gersby - go-billy Extensions

Extensions to the go filesystem abstraction [go-billy](https://github.com/go-git/go-billy)

## Extensions

### gersby.Walk

```go
func Walk(fs billy.Filesystem, root string, walkFn filepath.WalkFunc) error
```

Walks a `billy.Filesystem`, starting at `root`. Modeled after `filepath.Walk`


## Why gersby?

go-billy is named after the popular Ikea bookcase. Gersby is also an Ikea bookcase, slightly bigger - i.e. extended.