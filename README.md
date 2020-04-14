# sodago
`sodago` is an implementation of a subset of the [soda4LCA
API](https://bitbucket.org/okusche/soda4lca). It is not intended to be used in
production but only for testing purposes. However, it comes with a very fast
server and is trivial to set up which makes it a fun tool when
developing/testing against the [soda4LCA service API](https://bitbucket.org/okusche/soda4lca/src/5.0.3/Doc/src/Service_API/Service_API.md) in
[openLCA](https://github.com/GreenDelta/olca-app) or the [EPD
Editor](https://github.com/GreenDelta/epd-editor). 

## Usage
`sodago` is written in [Go](https://golang.org) and compiles to a single binary:

```bash
cd sodago
go build  # compile it
./sodago  # run it
```

## Data storage
All data are stored as plain files in the `data` folder. Within this folder
the content of each data stock is stored in a sub-folder. The name of the
data stock is the name of this sub-folder and the UUID of the data stock is
stored in a file `.stock` in that folder. Additionally, the content of a
data stock is stored in the `index.json` file. The data sets are stored under
their respective paths in the data folder (flows in `flows`, processes in
`processes` etc.). The file name of a data set is simply `<uuid>_<version>.xml`.  

External documents of sources are stored under `external_docs/<source uuid>/<file>`.
Profiles are stored in the sub-folder `profiles` of the `data` directory. Thus,
you should not name a data stock `profiles`. Each profile is stored in a
separate JSON file where the name has the following pattern: `<profile ID>.json`

## Implemented service routes
The prefix `/resource` is always added to all service routes (as in soda4LCA):

* `GET /datastocks`
* `GET /profiles`
* `GET /profiles/{id}`
* `/authenticate/login` (returns always OK!)
* `GET [/datastocks/{datastock}]/{path}`
* `GET [/datastocks/{datastock}]/{path}/{id}[?version={version}]`
* `GET [/datastocks/{datastock}]/sources/{id}/{file}`
* `POST /{path}`
* `POST /sources/withBinaries`

TODO:
* implement: GET [/sources/{uuid}/digitalfile](https://bitbucket.org/okusche/soda4lca/src/c78970a1d3ddaf855745b938082cee9cac1363e7/Doc/src/Service_API/Service_API_Dataset_Source_GET_DigitalFile.md)
* put types to lower case (e.g. Version)
* normalize versions (1 == 1.00 == 1.00.000)
* simple search
* LCIA methods