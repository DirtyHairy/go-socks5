This is a simple SOCKS5 proxy server forked of and based on
https://github.com/armon/go-socks5.

In particular, it allows to remap connections to particular IPs and ports.

Example:

```
    $ socks5-server -m 1.2.3.4:5.6.7.8 -m 10.9.44.2:80:127.0.0.1:8080
```

This will remap requests all requests to `1.2.3.4` to `5.6.7.8` and requests
to `10.9.44.2:80` to `127.0.0.1:8080`.
