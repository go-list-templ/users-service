const grpcUrl = 'app:8080'

export function create(client, payload) {
    client.connect(grpcUrl, {plaintext: true});

    const response = client.invoke('api.user.v1.UserService/Create', payload)

    client.close()

    return response
}

export function list(client, payload) {
    client.connect(grpcUrl, {plaintext: true});

    const response = client.invoke('api.user.v1.UserService/List', payload)

    client.close()

    return response
}

export function getByEmail(client, payload) {
    client.connect(grpcUrl, {plaintext: true});

    const response = client.invoke('api.user.v1.UserService/GetByEmail', payload)

    client.close()

    return response
}

export function verifyCred(client, payload) {
    client.connect(grpcUrl, {plaintext: true});

    const response = client.invoke('api.user.v1.UserService/VerifyCred', payload)

    client.close()

    return response
}