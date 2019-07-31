FROM busybox:1.31.0
COPY ca-bundle.crt /etc/ssl/certs/ca-certificates.crt
COPY dist/aws-rollout /usr/local/bin/aws-rollout
RUN ln -s /usr/local/bin/aws-rollout /aws-rollout