# Users service

Template users service on go

---

Run and deploy to from Helm to Kuber

```bash
werf converge --repo=ghcr.io/go-list-templ/users-service --dev
```

Run docker app container

```bash
werf run --dev
```

Build docker container

```bash
werf build --dev
```