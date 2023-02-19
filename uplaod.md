* /initUpload
  ```
  request {
    fileName:
    size:
    fileType:
  }
  ```
  ```
  response {
    status:
    requestID: 
  }
  ```
  backend
  ```go
  map[requestId]io.WriteSeeker
  // on failing reset seek to its original position
  
  ```

* /uploadChunk 
  ```
  request {
    requestID:
    index:
    chunkSize:
    chunkData:
  }
  ```
  ```
  response {
    status:
  }
  ```

* /finishUpload 
  ``` 
  request {
    requestID:
    checksum:
  }
  ```
  ```
  response {
    status:
  }
  ```

test
```bash
cat <<EOF | tee part1.txt
1234
EOF
cat <<EOF | tee part2.txt
5678
EOF
cat <<EOF | tee part3.txt
9012
EOF

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"fileName":"test.txt","size":12,"fileType":"txt"}' \
  http://localhost:8090/v1/initChunkUpload

# set the request ID
requestID="a85e19c7-3122-4c84-be0c-87eb2220a7ba"


curl -v \
  -F requestID=${requestID} \
  -F index=0 \
  -F chunkSize=4 \
  -F chunkData=@part1.txt \
  http://localhost:8090/v1/uploadChunk


curl -v \
  -F requestID=${requestID} \
  -F index=4 \
  -F chunkSize=4 \
  -F chunkData=@part2.txt \
  http://localhost:8090/v1/uploadChunk


curl -v \
  -F requestID=${requestID} \
  -F index=8 \
  -F chunkSize=4 \
  -F chunkData=@part3.txt \
  http://localhost:8090/v1/uploadChunk

curl --header "Content-Type: application/json" \
  --request POST \
  --data "{\"requestID\":\"${requestID}\",\"checksum\":\"txt\"}" \
  http://localhost:8090/v1/finishChunkUpload


rm part1.txt part2.txt part3.txt 

```