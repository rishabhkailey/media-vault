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
* check requests with multiple ranages https://www.ietf.org/archive/id/draft-ietf-httpbis-p5-range-09.html#internet.media.type.multipart.byteranges
