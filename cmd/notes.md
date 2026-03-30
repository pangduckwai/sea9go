# Publishing this module

### Commit source code to git
```bash
$ git status
$ git add .
$ git commit -m "[commit message]"
```

### Tag version
```bash
$ git tag vX.Y.Z
```

### Publish to github
```bash
$ git push --atomic origin master vX.Y.Z
```

### Update go modules index
```bash
$ GOPROXY=proxy.golang.org go list -m github.com/pangduckwai/sea9go@vX.Y.Z
```
