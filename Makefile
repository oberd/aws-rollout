all: build_image deps install release

clean:
	@rm -rf dist && mkdir dist

compile: clean
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dist/aws-rollout aws-rollout.go

deps:
	go get github.com/c4milo/github-release
	go get github.com/mitchellh/gox

install:
	go install -ldflags "-X main.Version=v1.0.8"

build_image: compile
	@docker build -t oberd/aws-rollout .

release: compile
	@./release.sh
	@docker push oberd/aws-rollout
	@docker tag -f oberd/aws-rollout oberd/aws-rollout:$$(git describe --tags `git rev-list --tags --max-count=1`)
	@docker push oberd/aws-rollout:$$(git describe --tags `git rev-list --tags --max-count=1`)
