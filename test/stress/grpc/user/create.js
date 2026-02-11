export function createUserTest(client) {
    const payload = {
        name: `User_${__ITER}`,
        email: `mail${__VU}_${__ITER}@example.com`,
    };

    return client.invoke('api.user.v1.UserService/CreateUser', payload)
}