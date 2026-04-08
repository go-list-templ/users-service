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
  --docker-server='ghcr.io/go-list-templ/users-service' \
  --docker-username='go-list-templ' \
  --docker-password='GH_TOKEN' \
  -n users-service
```

Login

```bash
werf cr login -u go-list-templ -p GH_TOKEN ghcr.io/go-list-templ/users-service
```

---

## Install dependency helm

```bash
cd .helm ; helm dependency update
```

---

## Run and build App

Run and deploy to from Helm to Kuber

```bash
werf converge --namespace users-service --repo=ghcr.io/go-list-templ/users-service --platform=linux/amd64 --dev
```

Build docker container

```bash
werf build --platform=linux/amd64 --dev
```

Stop and remove Kuber Pods with Service

```bash
kubectl delete all --all --namespace users-service
```

Forward port on localhost from app

```bash
kubectl port-forward svc/users-service 8080:8080 -n users-service
kubectl port-forward svc/users-service 8081:8081 -n users-service
```