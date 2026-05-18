docker_build(
    'users-service-image',
    '.',
    live_update=[
        sync('.', '/app')
    ]
)


watch_file('.k8s')
k8s_manifests = local(
    'kustomize build .k8s/overlays/dev --load-restrictor LoadRestrictionsNone --enable-alpha-plugins --enable-exec --enable-helm'
)

k8s_yaml(k8s_manifests)

k8s_resource(
    'users-service',
    port_forwards=['8080:8080', '8081:8081']
)