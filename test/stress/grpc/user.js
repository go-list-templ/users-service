export function create(client, payload) {
    return client.invoke('api.user.v1.UserService/Create', payload)
}

export function list(client, payload) {
    return client.invoke('api.user.v1.UserService/List', payload)
}

export function getByEmail(client, payload) {
    return client.invoke('api.user.v1.UserService/GetByEmail', payload)
}

export function verifyCred(client, payload) {
    return client.invoke('api.user.v1.UserService/VerifyCred', payload)
}