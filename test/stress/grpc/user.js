export function create(client, payload) {
    return client.invoke('api.user.v1.UserService/Create', payload)
}

export function list(client, payload) {
    return client.invoke('api.user.v1.UserService/List', payload)
}