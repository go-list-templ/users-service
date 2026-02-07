import grpc from 'k6/net/grpc';
import {createUserTest} from './grpc/user/create.js';
import {allUsersTest} from './grpc/user/all.js';
import {healthCheckTest} from "./http/diagnostic/healthz.js";

const grpcUrl = 'app:8080'
const diagnosticUrl = 'http://app:8081'
const client = new grpc.Client();

export const options = {
    scenarios: {
        create_user_grpc: {
            executor: 'constant-vus',
            exec: 'runCreateUser',
            vus: 1,
            stages: [
                {duration: '10s', target: 10},
                {duration: '20s', target: 10},
            ],
        },
        all_users_grpc: {
            executor: 'ramping-vus',
            exec: 'runAllUsers',
            vus: 1,
            stages: [
                {duration: '10s', target: 10},
                {duration: '20s', target: 10},
            ],
        },
        healthz_http: {
            executor: 'constant-vus',
            exec: 'runHealthCheck',
            vus: 2,
            stages: [
                {duration: '10s', target: 10},
                {duration: '20s', target: 10},
            ],
        },
    },
    thresholds: {
        'grpc_req_duration{scenario:create_user_grpc}': ['p(95)<100'],
        'http_req_duration{scenario:healthz_http}': ['p(99)<500'],
    },
};

export function runCreateUser() {
    client.connect(grpcUrl, { plaintext: true, reflect: true });

    createUserTest(client);
}

export function runAllUsers() {
    client.connect(grpcUrl, { plaintext: true, reflect: true });

    allUsersTest(client);
}

export function runHealthCheck() {
    healthCheckTest(diagnosticUrl)
}