import grpc from 'k6/net/grpc';
import {check, sleep} from 'k6';

export function createUserTest(client) {
    const payload = {
        name: `User_${__ITER}`,
        email: `parallel_${__VU}_${__ITER}@example.com`,
    };

    const response = client.invoke('v1.UserService/CreateUser', payload);

    check(response, {
        'create_user: status is OK': (r) => r && r.status === grpc.StatusOK,
    });

    sleep(1);
}