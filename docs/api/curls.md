## Search Album
```bash
curl \
    -H 'Authorization: Bearer <token>' \
    -H 'Content-Type: application/json' \
    -X POST \
    -d '{"order": "created_at", "sort": "desc", "page": 1, "per_page": 10, "query": "test"}' \
    http://localhost:8090/v1/search/album

```