build:
	go build

docker:
	docker build -t sosedoff/slacklet .

docker-release: docker
	docker push sosedoff/slacklet