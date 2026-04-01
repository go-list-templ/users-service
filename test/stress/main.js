import grpc from 'k6/net/grpc'
import {check} from 'k6'
import {SharedArray} from 'k6/data';
import {create, getByEmail, list, verifyCred} from "./grpc/user.js"
import {generateUsersData, getRandomItem} from "./helpers/helpers.js";

const tokens = {}

const usersData = new SharedArray('users data', () => {
    return generateUsersData(100)
});

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
            exec: 'runCreate',
        },
        verify_cred: {
            executor: 'constant-arrival-rate',
            rate: 10,
            timeUnit: '1s',
            duration: '1m',
            preAllocatedVUs: 10,
            exec: 'runVerifyCred',
        },
        get_by_email: {
            executor: 'constant-arrival-rate',
            rate: 100,
            timeUnit: '1s',
            duration: '1m',
            preAllocatedVUs: 10,
            exec: 'runGetByEmail',
        },
        list: {
            executor: 'constant-arrival-rate',
            rate: 100,
            timeUnit: '1s',
            duration: '1m',
            preAllocatedVUs: 10,
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

export function runCreate() {
    const user = getRandomItem(usersData)

    const payload = {
        name: user.name,
        email: user.email,
        password: user.password
    };

    const response = create(client, payload)

    check(response, {
        'create status is OK': (r) => r && r.status !== grpc.StatusInternal,
    })
}

export function runList() {
    const vu = __VU;

    const payload = {
        page_token: tokens[vu]
    };

    const response = list(client, payload)

    check(response, {
        'list status is OK': (r) => r && r.status !== grpc.StatusInternal,
    })

    if (response && response.status === grpc.StatusOK) {
        tokens[vu] = JSON.parse(JSON.stringify(response.message)).nextPageToken
    }
}

export function runGetByEmail() {
    const user = getRandomItem(usersData)

    const payload = {
        email: user.email,
    };

    const response = getByEmail(client, payload)

    check(response, {
        'get_by_email status is OK': (r) => r && r.status !== grpc.StatusInternal,
    })
}

export function runVerifyCred() {
    const user = getRandomItem(usersData)

    const payload = {
        email: user.email,
        password: user.password
    };

    const response = verifyCred(client, payload)

    check(response, {
        'verify_cred status is OK': (r) => r && r.status !== grpc.StatusInternal,
    })
}