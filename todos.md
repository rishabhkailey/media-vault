## todos

* image/video upload
* image/video stream 
* encryption
* try js blob for custom encrypted stream https://github.com/videojs/video.js/issues/5926
* try https://developer.mozilla.org/en-US/docs/Web/API/Service_Worker_API for updating incoming stream
* dash + aes https://github.com/Dash-Industry-Forum/dash.js/issues/1993
* check Minio.GetObject they do seems to have logic on caching and range we may don't need to do that
* video thumbnail https://jsfiddle.net/rodrigo_silveira/tq1u07tz/1/
    * try webassembly or ffmpeg wasm
    * https://github.com/gpac/mp4box.js/ (https://jsbin.com/mugoguxiha/edit?html,output)
    * webcodecs 
* video metadata (length).
    * try mp4box.js
* https://sanjeev-pandey.medium.com/understanding-the-mpeg-4-moov-atom-pseudo-streaming-in-mp4-93935e1b9e9a
    * https://stackoverflow.com/questions/23787712/qt-faststart-windows-how-to-run-it
* offscreencanvas https://developer.mozilla.org/en-US/docs/Web/API/OffscreenCanvas
* test memory usage with big files
* post proccessing encoding to v9 and fast mov (it should be optional) (vp9 ffmpeg commands https://developers.google.com/media/vp9/get-started and https://stackoverflow.com/questions/6954845/how-to-create-a-webm-video-file) **lets go ahead with vp9 codec and webm container**
* jpeg for thumbnail
* https://github.com/FFmpeg/FFmpeg/blob/master/doc/examples/transcoding.c and https://github.com/FFmpeg/FFmpeg/tree/master/doc/examples
* try https://developer.mozilla.org/en-US/docs/Web/Guide/Audio_and_video_manipulation and  https://developer.mozilla.org/en-US/docs/Web/API/Canvas_API/Manipulating_video_using_canvas for thumbnail for video's on client side
* https://www.youtube.com/watch?v=SdePc87Ffik&ab_channel=Pusher (e2e encryption)
* https://github.com/mozilla/send
* we cannot have 2 service workers in 1 scope so let's write our own code for stream download
* continue download even after closing the tab https://developer.mozilla.org/en-US/docs/Web/API/Background_Fetch_API
* image location from metadata (use for search) (https://www.npmjs.com/package/exif-js)
* password storage client side (https://developer.mozilla.org/en-US/docs/Learn/JavaScript/Client-side_web_APIs/Client-side_storage)
* status column in media table
* code structure https://github.com/grafana/grafana/blob/main/pkg/services
* upload request returns the media so it can be added to the list and user can see it without page refresh 

## Problems
* avi format not supported so we might need to transcode videos and files
* videos in firefox not working

## Useful
* https://github.com/minio/minio-go/blob/master/examples/s3/



## Check later

* minio [oidc](https://min.io/docs/minio/linux/developers/security-token-service.html?ref=docs)
* 1 bucket per user? then we can just use minio quota feature for limiting users
* as this will be self hosted MINM/proxied verify end to end cert?
* do we need to close reader bodies of minio object? 
* do we need to cache minio objects?
* getting `time="2023-01-24T02:19:24Z" level=error msg="write tcp 127.0.0.1:8090->127.0.0.1:57824: write: broken pipe"`
* https://github.com/seaweedfs/seaweedfs
* https://blog.min.io/minio-optimizes-small-objects/
* https://en.wikipedia.org/wiki/VP9
* https://www.npmjs.com/package/resumablejs
* js performance test https://jsben.ch/
* for file meteada - https://www.npmjs.com/package/mediainfo.js (for now continuing with file extensions)

## Pilain
* every user will have 2 buckets 
    * encrypted and locked with password(password will be the aes key)
    * unencrypted and without any lock
* able to show photos and videos like a gallery on local storage, as well as server storage
* able to upload normal files for backup
* file names in minio random/uuid. we will store the actual name in the DB
* we will need to have something similar to range requests and every request should have total count so we can have placeholders in application so user can just keep on scrolling without loading. we will have lazy load for images/videos.


## Auth
* https://www.ory.sh/docs/ecosystem/projects 
    * roles - used for user access - openid scopes e.g. annnonymous, user, admin
    * scopes - used by applications to access user data - oauth scopes e.g. media.read, media.read_write

## Searching
* labels
    * file name (/\-_ .+) separated
    * file type (image, video, document, other)
    * month name
    * year
    * manual tags
    * location
    * album name
* https://github.com/typesense/typesense
* https://github.com/apache/lucene
* https://docs.meilisearch.com/learn/what_is_meilisearch/overview.html
* https://blog.meilisearch.com/why-should-you-use-meilisearch-over-elasticsearch/

* user time zones

new user_info table?
```js
const tz = Intl.DateTimeFormat().resolvedOptions().timeZone
```

```go
//init the loc
loc, _ := time.LoadLocation("Asia/Shanghai")
//set timezone,  
now := time.Now().In(loc)
```

### MeiliSearch
* facet based on media type
* [filtering](https://docs.meilisearch.com/learn/getting_started/filtering_and_sorting.html#settings) based on userID
* we can [add or update](https://docs.meilisearch.com/reference/api/documents.html#add-or-update-documents) the document using put request 
#### document Schema
> todo geo location
```json
{
    "media_id": "mediaID", // primary key
    "user_id": "userID", // for filter
    "metadata": {
        "name": "", // search name
        "timestamp": "", // unix time stamp for filters, before and after some time
        "date": "", // wed 19 feb 2023
        "type": "", // like video, audio
    }
}
```
synonyms for day and months
## Issues
* todo mutex for map
```log
"concurrent map writes"
Stack:
	2  0x0000000001235719 in github.com/rishabhkailey/media-service/internal/router/api/v1.newUploadRequest
	    at /workspaces/media-service/internal/router/api/v1/chunkedUpload.go:42
	3  0x0000000001235f2d in github.com/rishabhkailey/media-service/internal/router/api/v1.(*Server).startUploadInBackground
	    at /workspaces/media-service/internal/router/api/v1/chunkedUpload.go:87
	4  0x0000000001235df1 in github.com/rishabhkailey/media-service/internal/router/api/v1.(*Server).InitChunkUpload.func1
	    at /workspaces/media-service/internal/router/api/v1/chunkedUpload.go:78
```
* happened when upload and download at same time?
```log
    [GIN] 2023/03/08 - 10:40:47 | 200 |     363.397µs |       127.0.0.1 | POST     "/v1/finishChunkUpload"
    time="2023-03-08T10:40:52Z" level=info msg="request received" range="{0 1169637}"
    time="2023-03-08T10:40:52Z" level=info msg=sent bytes=0
    time="2023-03-08T10:40:52Z" level=error msg="At least one of the pre-conditions you specified did not hold"
```
* make login persistence
* make the file upload dialog persistence? like it restart the upload if state is updated
* test handling of expired tokens
* todo domain for the cookies
* way to clear sessions from redis
* redis key eviction policy and max memory (also different db of redis for different puposes with different limits and eviction policy)
* increase auth server token expire time, it is set to 2 hour right now
* encrypt file name? but then we will not be able to use file name in search



## Test Commands
```bash
curl -X GET 'http://localhost:8090/v1/mediaList?perPage=5&sort=desc&page=1' -H "Authorization: Bearer <token>" | tee test.json
curl -v -X GET 'http://localhost:8090/v1/media?file=8379ada2-e309-4d3a-b4b8-18d49211748e' -H "Authorization: Bearer <token>" | tee test.json
```
## next
* cache integration in all stores
* album and 1 default album (favourite)
* api tests https://gin-gonic.com/docs/testing/
* services and store tests
* health endpoint
* local storage store implementation
* on finish upload add the media to the list in client

* add cache for custom file type
* local store implementation and migration command
* Delete api using transactions - https://gorm.io/docs/transactions.html
* Try https://vuetifyjs.com/en/api/v-img/#props-gradient for hover effect