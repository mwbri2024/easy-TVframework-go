# easy-TVframework-go

## Golang multi-platform build support
go install mvdan.cc/garble@latest

## Docker multi-platform build support Using BuildKit local container
docker buildx create --name mybuilder --use --driver docker-container  
docker buildx inspect --bootstrap

## Build multi-platform image and push to the HUB repository _latest
sudo docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t g12777/easy-TV:latest --push .

## Deploy and run _latest
sudo docker run -d --restart always -p 8123:8123 --name easy-TV g12777/easy-TV:latest -token=mytoken

## Auto-update Sync_Update.sh Modify line 39 for server SSH information
### The website directory is defined in the config. If the local version differs from the online version, it triggers an automatic update.
