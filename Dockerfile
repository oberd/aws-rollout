FROM scratch
ADD ca-bundle.crt /etc/ssl/certs/
ADD ca-bundle.trust.crt /etc/ssl/certs/
ADD aws-rollout /aws-rollout
ENTRYPOINT ["/aws-rollout"]