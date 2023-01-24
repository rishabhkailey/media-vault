## todos

* image/video upload
* image/video stream 
* encryption
* try js blob for custom encrypted stream https://github.com/videojs/video.js/issues/5926
* dash + aes https://github.com/Dash-Industry-Forum/dash.js/issues/1993
* check Minio.GetObject they do seems to have logic on caching and range we may don't need to do that


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

## Pilain
* every user will have 2 buckets 
    * encrypted and locked with password(password will be the aes key)
    * unencrypted and without any lock
* able to show photos and videos like a gallery on local storage, as well as server storage
* able to upload normal files for backup
