import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

export function allUsersTest(client) {
    const response = client.invoke('api.user.v1.UserService/AllUsers', {});

    check(response, {
        'all_users: status is OK': (r) => r && r.status === grpc.StatusOK,
    });

    sleep(1);
}