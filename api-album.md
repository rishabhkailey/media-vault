# Tables
## albums
| ID | Name | thumbnail_url | 
| -- | -- | -- | 
| `uint` | `text` | `text` | 


## album_media_bindings
| album_id | media_id |
| -- | -- |
| `uint` | `uint` |

## user_album_bindings
| album_id|  user_id| 
| -- | -- |
| `uint` | `text` |

* Get user albums request
  ```json
  GET /v1/albums?page=1&perPage=10
  ```
  Response
  ```json
  [
    {
      "id": "uint",
      "name": "string",
      "thumbnail_url": "string"
    }
  ]
  ```
* Get album media list
  ```
  GET /v1/album/:album_id/media?page=1&perPage=10&order=date&sort=desc
  ```
  Response - normal media list response

* add media in album
  ```json
  POST /v1/album/:album_id/media
  {
    "media_ids": []
  }
  ```
* update thumbnail
  ```json
  POST /v1/album/:album_id/thumbnail
  {
    "media_id": 123
  }
  ```
