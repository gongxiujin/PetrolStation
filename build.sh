CGO_ENABLED=1 GOARCH=amd64 CC=x86_64-linux-musl-gcc GOOS=linux go build -o station main.go
scp -P 63000 ~/code/PetrolStation/station root@116.62.194.251:/root/