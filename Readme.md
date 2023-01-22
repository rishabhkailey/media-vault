## Range Requests
[RFC](https://www.ietf.org/archive/id/draft-ietf-httpbis-p5-range-09.html)
* request with if-range header
* 206 response for the if-range header or normal partial response
* 416 status code for invalid range request
* cache validator?
* server send Accept-Ranges: bytes header for indicating it supports range requests. Accept-Ranges: none for the opposite
* Content-Range: bytes 21010-47021/47022 (0 indexed bytes range/total bytes) if total length is unkown use * for total lenght
* Content-Length: 26012 (length of the response range, not the total lenght)
* If-Range request header - if the condition is fulfilled, the range request is issued, and the server sends back a 206 Partial Content answer with the appropriate body. If the condition is not fulfilled, the full resource is sent back with a 200 OK status.
* Range request header- bytes=<range-start>-<range-end>, ...
atleast 1 of range-start or range-end should be present. both will also work. if the range is invalid we can return the whole file with 200 status code.


## Todo
* check requests with multiple ranages https://www.ietf.org/archive/id/draft-ietf-httpbis-p5-range-09.html#internet.media.type.multipart.byteranges, https://developer.mozilla.org/en-US/docs/Web/HTTP/Range_requests#multipart_ranges

## Test video files

https://samples.mplayerhq.hu/
https://file-examples.com/index.php/sample-video-files/
https://sample-videos.com/



## Setup


Register client in auth server
```bash
docker exec media-service_devcontainer-auth-service-1 /etc/auth-server/auth-server --help
docker exec media-service_devcontainer-auth-service-1 /etc/auth-server/auth-server client create --scopes user --scopes admin --defaultScopes user --domain localhost --id media-service --secret aslkjfdalksdfjlksadfjlkasjlkfjasdfjlkasjd --returnUri localhost:8090
# to verify if client is registered try <auth-service-url>/v1/<client-id>/.well-known/openid-configuration
# e.g. http://localhost:8080/v1/media-service/.well-known/openid-configuration
```

scopes
* user
* admin
* anonymous