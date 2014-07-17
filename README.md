og
==

### OG a demo of a reverse proxy in GO

This small script was born with the idea to proxy [elasticseach](http://elasticsearch.org) calls and add a very basic authentication in front of it. The idea was to have elasticsearch bound to a private NIC and the proxy bound on a public NIC. The proxy would then forward all http calls to the local NIC as long as they were authenticated. Is pretty simple to add a more complex authentication scheme than a basic one.

This is an exercise towards my learning of [go lang](http://golang.org), showing how easy it was to create a very small application that consumes minimal resources (~300kb of memory under 10 concurrent connections)

For real applications you really should use nginx as a reverse proxy.

#### Usage

After building it, just call: ./og -config=*config_file*

Sample config:

```
{
  "url": "http://<local_nic>:<port>",
  "port": <listening port>,
  "authentication": {
    "credentials": [
      {
        "username": "es",
        "password": "bonsai"
      }
     
    ]
  }
}

```

The proxy will search for any Basic authentication (using base64 decode) header sent, and match against the credentials list.



