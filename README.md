# Users service

Template users service on go

---

## Create secret and login in docker registry

Create namespace

```bash
kubectl create namespace users-service
```

Create secret

```bash
kubectl create secret docker-registry registrysecret \
  --docker-server=ghcr.io/go-list-templ \
  --docker-username=go-list-templ \
  --docker-password=GH_TOKEN
```

Login

```bash
werf cr login ghcr.io/go-list-templ -u go-list-templ -p GH_TOKEN 
```

---

## Install dependency helm

```bash
werf helm dependency update .helm
```

---

## Run and build App

Run and deploy to from Helm to Kuber

```bash
werf converge --repo=ghcr.io/go-list-templ/users-service --platform=linux/amd64
```

Stop and remove release in kuber

```bash
werf dismiss
```

Forward port on localhost from app

```bash
werf kubectl port-forward svc/users-service 8080:8080 -n users-service
werf kubectl port-forward svc/users-service 8081:8081 -n users-service
```

Delete all images from container registry (token with rules on write+delete packages)

```bash
werf purge --repo ghcr.io/go-list-templ/users-service --repo-github-token GH_TOKEN
```