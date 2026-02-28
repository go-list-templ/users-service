import grpc from 'k6/net/grpc'
import {check} from 'k6'
import {create, list} from './grpc/user.js'
import {healthz} from "./http/healthz.js"

const grpcUrl = 'app:8080'
const diagnosticUrl = 'http://app:8081'

const client = new grpc.Client()
client.load(['/src/proto/api/user/v1'], 'user.proto')

const tokens = {};

export const options = {
    scenarios: {
        create_user_grpc: {
            executor: 'ramping-arrival-rate',
            startRate: 50,
            timeUnit: '1s',
            preAllocatedVUs: 10,
            maxVUs: 10,
            stages: [
                {target: 100, duration: '1m'},
            ],
            exec: 'runCreate',
        },
        list_users_grpc: {
            executor: 'ramping-arrival-rate',
            startRate: 50,
            timeUnit: '1s',
            preAllocatedVUs: 10,
            maxVUs: 50,
            stages: [
                {target: 100, duration: '30s'},
                {target: 500, duration: '30s'},
            ],
            exec: 'runList',
        },
        healthz_http: {
            executor: 'constant-arrival-rate',
            exec: 'runHealthz',
            rate: 2,
            timeUnit: '5s',
            duration: '1m',
            preAllocatedVUs: 1,
            maxVUs: 1,
        },
    },
    thresholds: {
        'grpc_req_duration{scenario:create_user_grpc}': ['p(95) < 100'],
        'grpc_req_duration{scenario:list_users_grpc}': ['p(95) < 100'],
        'http_req_duration{scenario:healthz_http}': ['p(95) < 500'],
        'checks': ['rate >= 0.9']
    },
    summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(95)', 'p(99)', 'count'],
}

export function runCreate() {
    client.connect(grpcUrl, {plaintext: true})

    const payload = {
        name: `user_${__ITER}`,
        email: `mail${__VU}_${__ITER}@example.com`,
    };

    const response = create(client, payload)

    check(response, {
        'create_user status is OK': (r) => r && r.status === grpc.StatusOK,
    })

    client.close()
}

export function runList() {
    client.connect(grpcUrl, {plaintext: true})

    const vu = __VU;
    const payload = {page_token: tokens[vu]};

    const response = list(client, payload)

    check(response, {
        'list_users status is OK': (r) => r && r.status === grpc.StatusOK,
    })

    if (response && response.status === grpc.StatusOK) {
        tokens[vu] = JSON.parse(JSON.stringify(response.message)).nextPageToken
    }

    client.close()
}

export function runHealthz() {
    const response = healthz(diagnosticUrl)

    check(response, {
        'healthz is 200': (r) => r.status === 200,
    })
}