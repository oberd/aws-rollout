FROM busybox
ADD ca-bundle.crt /etc/ssl/certs/ca-certificates.crt
ADD aws-rollout /aws-rollout
ENTRYPOINT ["/aws-rollout"]