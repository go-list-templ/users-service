import { sleep } from 'k6';

export function createUserTest(client) {
    const payload = {
        name: `User_${__ITER}`,
        email: `mail${__VU}_${__ITER}@example.com`,
    };

    const response =  client.invoke('api.user.v1.UserService/CreateUser', payload);

    sleep(0.5)

    return response
}