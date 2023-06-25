* /medialist
  ```go
  struct Media type {
    id uint
    metadata struct {
      name string
      date string
      mediaType string
      size uint
    }
  }
  ```
  ```go
  request {
    body {
      page int
      perPage int // some max limit set
      mediaType []string
      orderBy string // e.g. upload_date, file_creation_date
      sort string // asc, desc
    }
    headers {
      Authorization: 
    }
  }
  ```
  ```go
  response {
    body []Media,
    headers {
      x-total uint // number of media
      x-total-page uint // number of pages
    }
  }
  ```
* /media/:id 
  ```
  request {
    headers {
      Authorization: 
    }
  }
  ```
  ```
  response {
    header {
      content-type
      content-size
    }
  }
  ```

* /media/:mediaID/delete
multiple
response multiple failed and success