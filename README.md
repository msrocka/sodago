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

http://localhost/resource/flows/a2b32f97-3fc7-4af2-b209-525bc6426f33
http://localhost/resource/datastocks/982aa819-c7bc-4f89-97c2-783606f7fa66/flows/a2b32f97-3fc7-4af2-b209-525bc6426f33?version=30.00.000

### Profiles
Profiles are stored in the sub-folder `profiles` of the data directory.
Each profile is stored in a separate JSON file where the name has the
following pattern: `<profile ID>.json`

TODO:
* implement: GET [/sources/{uuid}/digitalfile](https://bitbucket.org/okusche/soda4lca/src/c78970a1d3ddaf855745b938082cee9cac1363e7/Doc/src/Service_API/Service_API_Dataset_Source_GET_DigitalFile.md)
* put types to lower case (e.g. Version)
* normalize versions (1 == 1.00 == 1.00.000)
* simple search
* LCIA methods
* write lock via simple mutex