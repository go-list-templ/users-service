import grpc from 'k6/net/grpc';
import {check} from 'k6';
import {createUserTest} from './grpc/user/create.js';
import {allUsersTest} from './grpc/user/all.js';
import {healthz} from "./http/diagnostic/healthz.js";

const grpcUrl = 'app:8080'
const diagnosticUrl = 'http://app:8081'

const client = new grpc.Client();

export const options = {
    scenarios: {
        create_user_grpc: {
            executor: 'constant-vus',
            exec: 'runCreateUser',
            vus: 100,
            duration: '30s',
        },
        all_users_grpc: {
            executor: 'constant-vus',
            exec: 'runAllUsers',
            vus: 100,
            duration: '30s',
        },
        healthz_http: {
            executor: 'constant-vus',
            exec: 'runHealthz',
            vus: 1,
            duration: '30s',
        },
    },
    thresholds: {
        'grpc_req_duration{scenario:create_user_grpc}': ['p(95) < 100'],
        'grpc_req_duration{scenario:all_users_grpc}': ['p(95) < 100'],
        'http_req_duration{scenario:healthz_http}': ['p(95) < 500'],
        'checks': ['rate > 0.9'],
    },
    summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(95)', 'p(99)', 'count'],
};

export function runCreateUser() {
    client.connect(grpcUrl, {plaintext: true, reflect: true});
    const response = createUserTest(client);

    check(response, {
        'create_user status is OK': (r) => r && r.status === grpc.StatusOK,
    }, {scenario: 'create_user'});
}

export function runAllUsers() {
    client.connect(grpcUrl, {plaintext: true, reflect: true});
    const response = allUsersTest(client);

    check(response, {
        'all_users status is OK': (r) => r && r.status === grpc.StatusOK,
    }, {scenario: 'all_users'});
}

export function runHealthz() {
    const response = healthz(diagnosticUrl)

    check(response, {
        'healthz is 200': (r) => r.status === 200,
    }, {scenario: 'healthz'});
}