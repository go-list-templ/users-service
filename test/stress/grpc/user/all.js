import { sleep } from 'k6';

export function allUsersTest(client) {
    const response =  client.invoke('api.user.v1.UserService/AllUsers', {});

    sleep(0.5)

    return response
}