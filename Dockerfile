FROM busybox
ADD ca-bundle.crt /etc/ssl/certs/ca-certificates.crt
ADD dist/aws-rollout /aws-rollout
ENTRYPOINT ["/aws-rollout"]