export function allUsersTest(client) {
    return client.invoke('api.user.v1.UserService/AllUsers', {})
}