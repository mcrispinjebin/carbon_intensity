IMAGE_NAME = carbon-intensity

build:
	docker build -t ${IMAGE_NAME} .

run:
	docker run -p 3000:3000 ${IMAGE_NAME}

test:
	go test -v ./...