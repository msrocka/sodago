# sodago

This is a small server application to mock the [soda4LCA Service
API](https://bitbucket.org/okusche/soda4lca). We use it to develop and test the
soda4LCA integration in [openLCA](https://github.com/GreenDelta/olca-app) and
the [EPD Editor](https://github.com/GreenDelta/epd-editor).


## Usage

`sodago` is written in [Go](https://golang.org) and compiles to a single binary:

```bash
cd sodago
go build  # compile it
./sodago  # run it
```

This will try to start a server at port `8080` using the `data` folder to store
the configuration and data files. Both can be changed with the following command
line options:

```bash
./sodago -port 8081 -data ./my-data
```

The data folder is created automatically when it does not exist yet. On start,
`sodago` also checks if users are configured and creates a default admin user if
this is not the case.


## Users and roles

The users of the server are configured in the `config.json` file in the data
folder. It has the following format:

```json
{
  "users": [
    {
      "user": "admin",
      "password": "default",
      "roles": [ "admin" ],
      "tokens": [ "a-fixed-admin-token" ]
    }
  ]
}
```

Currently the `roles` do not have a real meaning. They are only reported in the
`GET /resource/authenticate/status` response.

The optional `tokens` list contains pre-configured authentication tokens of a
user. They can be used directly in the `Authorization` header, without
requesting a token via the service API first.


## Authentication

Reading data is always allowed, also for anonymous users. Writing data (`POST`
requests) requires authentication, either with a session cookie or with a token:

* **Session**: `POST /resource/authenticate/login` with the form parameters
  `username` and `password`. The legacy
  `/resource/authenticate/login?userName={user}&password={password}` route is
  also supported. The session cookie is sent with each request and is closed
  with `GET /resource/authenticate/logout`.
* **Token**: `POST /resource/authenticate/getToken` with the form parameters
  `username` and `password` (legacy:
  `/resource/authenticate/getToken?userName={user}&password={password}`). The
  response is a signed token that is sent with each request in the
  `Authorization` header (`Authorization: Bearer {token}`). Tokens are valid
  for 90 days and do not need a logout. A requested token is appended to the
  `tokens` list of the user in the `config.json`, so that issued tokens are
tracked there.

`GET /resource/authenticate/status` reports the authenticated user of a
request. Requests with invalid credentials, for example an expired token, are
rejected with `403 Permission denied.`


## Data storage
All data are stored as plain files in the `data` folder with the following layout:

```
data
├── config.json         # the configured users (see above)
├── cookie_auth.key     # the signing key for the session cookie
├── token_auth.key      # the signing key for the authentication tokens
└── root                # a data stock, the folder name is its name
    ├── .stock          # the UUID of the data stock
    ├── index.json      # the index of the data sets stored in this stock
    ├── processes       # data sets by type, the folder name is the type
    │   └── <uuid>_<version>.xml   # the version is normalized (01.00.000)
    ├── flows           # (or flowproperties, lciamethods, sources,
    └── external_docs   #  unitgroups, contacts)
        └── <source uuid>
            └── <file>
```

Each sub-folder of the `data` folder that contains a `.stock` file is a data
stock. The name of the sub-folder is the name of the data stock and the `.stock`
file contains the UUID of the data stock. The `root` is created automatically
when it does not exist. The `index.json` file maps the data set types to the
UUIDs, versions and names of the stored data sets for faster lookups.

## Implemented service routes
The prefix `/resource` is always added to all service routes (as in soda4LCA):

* `GET /datastocks`
* `POST /authenticate/login` (form: `username`, `password`)
* `/authenticate/login?userName={user}&password={password}`
* `POST /authenticate/getToken` (form: `username`, `password`)
* `/authenticate/getToken?userName={user}&password={password}`
* `GET /authenticate/logout`
* `GET /authenticate/status`
* `Authorization: Bearer {token}` (session or token for all routes)
* `GET [/datastocks/{datastock}]/{path}`
* `GET [/datastocks/{datastock}]/{path}/{id}[?version={version}]`
* `GET [/datastocks/{datastock}]/sources/{id}/{file}`
* `GET [/datastocks/{datastock}]/sources/{id}/digitalfile`
* `POST /{path}`
* `POST /sources/withBinaries`

List requests (`GET [/datastocks/{datastock}]/{path}`) are paged and support
the query parameters `startIndex`, `pageSize`, `countOnly` and `allVersions`.
By default, only the most recent version of a data set is returned. With
`search=true` the result is filtered by the `name` parameter, e.g.
`GET /contacts?search=true&name=electricity`: the search phrase is split into
whitespace separated keywords that all have to occur in the name of a data set
(case-insensitive). With a `uuid:` prefix, the search phrase is interpreted as
a UUID prefix instead, e.g. `GET /contacts?search=true&name=uuid:42217b74`.

Versions are stored in the normalized ILCD format, so `1`, `1.0` and `01.00.000`
refer to the same version of a data set.

The route `GET [/datastocks/{datastock}]/sources/{id}/digitalfile` returns the
first digital file (the first `referenceToDigitalFile` entry) of a source, while
`.../sources/{id}/{file}` returns the file with the given name.
