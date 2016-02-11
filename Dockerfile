FROM scratch
ADD aws-rollout /aws-rollout
ENTRYPOINT ["/aws-rollout"]