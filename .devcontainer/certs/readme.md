## Minio sever certs 
> Note: these steps are for dev setup only
```bash
cd .devcontainer/certs
wget https://github.com/minio/certgen/releases/download/v1.2.1/certgen-linux-amd64 -O certgen
chmod +x certgen
./certgen -host "localhost,minio,172.17.0.1"
```