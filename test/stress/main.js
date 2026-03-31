import grpc from 'k6/net/grpc'
import {check} from 'k6'
import {create, getByEmail, list, verifyCred} from "./grpc/user.js"

const tokens = {}

const client = new grpc.Client()
client.load(['/src/proto/api/user/v1'], 'user.proto')

export const options = {
    scenarios: {
        create: {
            executor: 'constant-arrival-rate',
            rate: 10,
            timeUnit: '1s',
            duration: '1m',
            preAllocatedVUs: 10,
            maxVUs: 10,
            exec: 'runCreate',
        },
        verify_cred: {
            executor: 'constant-arrival-rate',
            rate: 10,
            timeUnit: '1s',
            duration: '1m',
            preAllocatedVUs: 10,
            maxVUs: 10,
            exec: 'runVerifyCred',
        },
        get_by_email: {
            executor: 'constant-arrival-rate',
            rate: 100,
            timeUnit: '1s',
            duration: '1m',
            preAllocatedVUs: 10,
            maxVUs: 50,
            exec: 'runGetByEmail',
        },
        list: {
            executor: 'constant-arrival-rate',
            rate: 100,
            timeUnit: '1s',
            duration: '1m',
            preAllocatedVUs: 10,
            maxVUs: 50,
            exec: 'runList',
        },
    },
    thresholds: {
        'grpc_req_duration{scenario:create}': ['p(95) < 100'],
        'grpc_req_duration{scenario:list}': ['p(95) < 100'],
        'grpc_req_duration{scenario:verify_cred}': ['p(95) < 100'],
        'grpc_req_duration{scenario:get_by_email}': ['p(95) < 100'],
        'checks': ['rate >= 0.9']
    },
    summaryTrendStats: ['min', 'max', 'p(95)', 'p(99)', 'count'],
}

export function setup() {
    return {
        name: `user`,
        email: `example${__VU}@gmail.com`,
        password: "password"
    };
}

export function runCreate(data) {
    const payload = {
        name: data.user,
        email: data.email,
        password: data.password
    };

    const response = create(client, payload)

    check(response, {
        'create status is OK': (r) => r && r.status === grpc.StatusOK,
    })
}

export function runList() {
    const vu = __VU;

    const payload = {
        page_token: tokens[vu]
    };

    const response = list(client, payload)

    check(response, {
        'list status is OK': (r) => r && r.status === grpc.StatusOK,
    })

    if (response && response.status === grpc.StatusOK) {
        tokens[vu] = JSON.parse(JSON.stringify(response.message)).nextPageToken
    }
}

export function runGetByEmail(data) {
    const payload = {
        email: data.email,
    };

    const response = getByEmail(client, payload)

    check(response, {
        'get_by_email status is OK': (r) => r && r.status === grpc.StatusOK,
    })
}

export function runVerifyCred(data) {
    const payload = {
        email: data.email,
        password: data.password
    };

    const response = verifyCred(client, payload)

    check(response, {
        'verify_cred status is OK': (r) => r && r.status === grpc.StatusOK,
    })
}