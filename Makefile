all: build_image deps install release

compile:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o aws-rollout aws-rollout.go

deps:
	go get github.com/c4milo/github-release
	go get github.com/mitchellh/gox

install:
	go install -ldflags "-X main.Version=v1.0.8"

build_image: compile
	@docker build -t oberd/aws-rollout .

release:
	@./release.sh
	@docker push oberd/aws-rollout:latest
	@docker push oberd/aws-rollout:$$(git describe --tags `git rev-list --tags --max-count=1`)
