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

## Pilain
* every user will have 2 buckets 
    * encrypted and locked with password(password will be the aes key)
    * unencrypted and without any lock
* able to show photos and videos like a gallery on local storage, as well as server storage
* able to upload normal files for backup
* file names in minio random/uuid. we will store the actual name in the DB
* we will need to have something similar to range requests and every request should have total count so we can have placeholders in application so user can just keep on scrolling without loading. we will have lazy load for images/videos.