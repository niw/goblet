Simple Proxy Server
===================

A simple version of proxy server glue code.

Usage
-----

Build it and run as proxy server.

```
go build simple-proxy-server/main.go
./main --cache_root=`pwd`/cache
```

When clone repositories, give `http_proxy` also use `http://` instead of
`https://` for clone URL.

```
env http_proxy=localhost:8080 git clone http://host/path.git
```

Then it will create a cache in `cache` directory cloned with TLS.
