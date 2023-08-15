# Services
we will be starting with minimum number of services and we will keep divinding big services to small services if we see benifits

* Search Service
* Media
* Media Storage
* User Service
* Album Service

## Media Service
### Tables
* media
* media_metadata
* user_media_bindings


### Methods
* Create -> store.mediaMetadata.create() -> store.media.create() -> store.userMediaBinding.create()
<!-- or just 1 store -->
* Create -> store.media.create(Media{Metadata: {}}, userID, uploadRequestID)
* DeleteOne
* DeleteMany
* GetMediaByUploadRequestID(uploadRequestID)
* GetMediaWithMetadataByUploadRequestID()
<!-- or combine above 2 methods with preloadMetadata parameter -->
<!-- seems better we can do something like db = s.db, if preloadMetadata; s = s.preload() -->
* GetMediaByFileName
* GetMediaListByUserID
* GetUserMediaByID
* GetMediaByIDs
* GetMediaByID
* GetTypeByFileName
* CheckFileBelongsToUser
* CheckMediaBelongsToUser
* UpdateThumbnail
<!-- lets add description for all the methods -->

```go
		uploadRequest, err = server.UploadRequests.Create(c.Request.Context(), uploadrequests.CreateUploadRequestCommand{
			UserID: userID,
		})
		mediaMetadata, err := server.MediaMetadata.Create(c.Request.Context(), mediametadata.CreateCommand{
			Metadata: mediametadata.Metadata{
				Name: requestBody.FileName,
				Date: time.UnixMilli(requestBody.Date),
				Size: uint64(requestBody.Size),
				Type: requestBody.MediaType,
			},
		})
		uploadingMedia, err = server.Media.Create(c.Request.Context(), media.CreateMediaCommand{
			UploadRequestID: uploadRequest.ID,
			MetadataID:      mediaMetadata.ID,
		})
		_, err = server.UserMediaBindings.Create(c.Request.Context(), usermediabindings.CreateCommand{
			UserID:  userID,
			MediaID: uploadingMedia.ID,
		})
	err = server.MediaStorage.InitChunkUpload(c.Request.Context(), mediastorage.InitChunkUploadCmd{
		UserID:    userID,
		RequestID: uploadRequest.ID,
		FileName:  uploadingMedia.FileName,
		FileSize:  requestBody.Size,
	})

	
	uploadRequestID, fileName := InitUpload() // it will create uploadRequest and media in media storage
	createMedia()
```

### store

