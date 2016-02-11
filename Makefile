all: compile build_image

build_image:
	@docker build -t oberd/aws-rollout .

compile:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o aws-rollout aws-rollout.go

zip:
	tar -zcvf aws-rollout-linux-x86_64.tar.gz aws-rollout
