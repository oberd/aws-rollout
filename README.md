#### AWS Rollout Tool

##### Usage:

```bash
aws-rollout [[options]] [service-name] [image]
```

##### Examples

```bash
aws-rollout --cluster=my-app my-app-backend-prod company/my-app:master-1234
```

```bash
docker run -it --rm \
  -e "AWS_REGION=****" \
  -e "AWS_ACCESS_KEY_ID=****" \
  -e "AWS_SECRET_ACCESS_KEY=****" \
  oberd/aws-rollout --cluster=my-app my-app-backend-prod company/my-app:master-1234
```

##### Options

* `cluster` name of cluster


##### Development

