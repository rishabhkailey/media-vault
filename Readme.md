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
https://upload.wikimedia.org/wikipedia/commons/a/a4/BBH_gravitational_lensing_of_gw150914.webm
https://en.wikipedia.org/wiki/Category:Video_samples_of_films


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


moov atom
```bash
# check position of moov
ffmpeg -v trace -i test-files/verticle.mp4 2>&1 | grep -e type:\'mdat\' -e type:\'moov\'

# change moov atom position
# https://stackoverflow.com/questions/8061798/post-processing-in-ffmpeg-to-move-moov-atom-in-mp4-files-qt-faststart
```

https://trac.ffmpeg.org/wiki

# Let's start with non e2e encrypted, unencrypted storage and without transcoding
* will add feature for encrypted storage and transcoding later
* will add feature for e2e encryption later


Go Profiling
```bash
# cpu
wget -O cpu.pprof http://localhost:8090/debug/pprof/profile?seconds=120
sudo apt install -y graphviz
go tool pprof -http 0.0.0.0:8989 cpu.pprof
```

system CPU debugging (helpful if other service like minio, postgres or redis is using high cpu)
[detailed steps](https://github.com/brendangregg/FlameGraph)
```bash
# on host machine not inside the container
# if perf command is not installed
# ubuntu
sudo apt install linux-tools-common linux-tools-generic linux-tools-`uname -r`
# debian
sudo apt install linux-perf

# will generate perf.data file
perf record -F 99 -a -g -- sleep 120
# read perf.data file
perf script > out.perf
# if flamegraph is not installed
git clone https://github.com/brendangregg/FlameGraph.git
export PATH=$PATH:$PWD/FlameGraph
cd ..

stackcollapse-perf.pl out.perf > out.folded
flamegraph.pl out.folded > kernel.svg
```



## DB

<!-- todo sas token for URLs or use cookie/session -->
<!-- todo do we need url in table? we can also generate url from id/name -->
<!-- todo use id or name for the file name? -->
Media
| id | name  | upload status | Thumbnail | Media type | media date
| -- | -- |  -- | -- | -- | -- | 
| `int64` | `text` | `failed`/`success`/`inProgress` | `boolean` | `video/mp4, video/webm` ...  | `timestamp` |


<!-- https://developer.mozilla.org/en-US/docs/Web/Media/Formats/Image_types -->
<!-- for now lets just go ahead with video mp4/webm and image png/jpeg -->
Media Types
| extension | type |
| -- | -- |
| `.apng` | `image/apng` |

<!-- different metadata table? -->

UserMediaBinding
| User ID | Media ID | 
| -- | -- |
| `int64` | `int64` |



Media - id, random-uuid file name

media metadata - media id (pkey), actual file name, file date, type

thumbnail - media id (pkey), random uuid file name

upload requests - id uuid?, media id, status, user id