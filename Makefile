build:
	go build

build-linux:
	GOOS=linux GOARCH=amd64 go build -o makemea-linux

deploy: build-linux
	cat makemea-linux | ssh makemea "cat > /home/ubuntu/tables/makemea-new && sudo service makemea stop && cd /home/ubuntu/tables && mv makemea-new makemea && chmod +x makeamea && sudo service makemea start"