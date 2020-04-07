# sodago
`sodago` is an implementation of a subset of the [soda4LCA
API](https://bitbucket.org/okusche/soda4lca). It is not intended to be used in
production but only for testing purposes. However, it comes with a very fast
server and is trivial to set up which makes it a fun tool when
developing/testing the soda4LCA interface in
[openLCA](https://github.com/GreenDelta/olca-app) or the [EPD
Editor](https://github.com/GreenDelta/epd-editor). 

## Usage
`sodago` is written in [Go](https://golang.org) and compiles to single binary:

```bash
cd sodago
go build  # compile it
./sodago  # run it
```

After this you just need to start the `sodago` executable.
