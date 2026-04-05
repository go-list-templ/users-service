# Users service

Template users service on go

---

Run and deploy to from Helm to Kuber

```bash
werf converge --repo=ghcr.io/go-list-templ/users-service --platform=linux/amd64 --dev
```

Build docker container

```bash
werf build --platform=linux/amd64 --dev
```

Stop and remove Kuber Pods with Service

```bash
kubectl delete all --all --namespace users-service
```