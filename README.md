# sodago
`sodago` is an implementation of a subset of the
[soda4LCA API](https://bitbucket.org/okusche/soda4lca). It is not intended to
be used in production but only for testing purposes. However, it comes with
a very fast server and is trivial to set up which makes it a fun tool when
developing/testing the soda4LCA interface in 
[openLCA](https://github.com/GreenDelta/olca-app) or the
[EPD Editor](https://github.com/GreenDelta/epd-editor). 

## Installation
Currently no binary distributions are provided. Thus, you have to compile it
from source with [Go](https://golang.org):

```bash
cd sodago
go build
```

After this you just need to start the `sodago` executable.

## Testing with curl

```
curl -d "@6b1b696a-3ea8-48c6-9e8a-f0cc5ce00dba_01.00.000.xml" -X POST http://localhost/resource/flows
```