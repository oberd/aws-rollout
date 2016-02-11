#### AWS Rollout Tool

##### Usage:

```bash
aws-rollout [[options]] [service-name] [image]
```

##### Examples

```bash
aws-rollout --cluster=my-app my-app-backend-prod company/my-app:master-1234
```

##### Options

* `cluster` name of cluster