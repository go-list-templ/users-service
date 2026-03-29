import grpc from 'k6/net/grpc'
import {check} from 'k6'
import {healthz} from "./http/healthz.js"
import {create, getByEmail, list, verifyCred} from "./grpc/user.js"

const grpcUrl = 'app:8080'
const diagnosticUrl = 'http://app:8081'

const tokens = {};
const usersCred = {name: "", password: "", email: ""};

const client = new grpc.Client()
client.load(['/src/proto/api/user/v1'], 'user.proto')

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
        verify_cred_user_grpc: {
            executor: 'ramping-arrival-rate',
            startRate: 5,
            timeUnit: '1s',
            preAllocatedVUs: 5,
            maxVUs: 5,
            stages: [
                {target: 50, duration: '1m'},
            ],
            exec: 'runVerifyCred',
        },
        get_by_email_user_grpc: {
            executor: 'ramping-arrival-rate',
            startRate: 5,
            timeUnit: '1s',
            preAllocatedVUs: 5,
            maxVUs: 5,
            stages: [
                {target: 50, duration: '1m'},
            ],
            exec: 'runGetByEmail',
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
            rate: 2,
            timeUnit: '5s',
            duration: '1m',
            preAllocatedVUs: 1,
            maxVUs: 1,
            exec: 'runHealthz',
        },
    },
    thresholds: {
        'grpc_req_duration{scenario:create_user_grpc}': ['p(95) < 100'],
        'grpc_req_duration{scenario:list_users_grpc}': ['p(95) < 100'],
        'grpc_req_duration{scenario:verify_cred_user_grpc}': ['p(95) < 100'],
        'grpc_req_duration{scenario:get_by_email_user_grpc}': ['p(95) < 100'],
        'http_req_duration{scenario:healthz_http}': ['p(95) < 500'],
        'checks': ['rate >= 0.9']
    },
    summaryTrendStats: ['min', 'max', 'p(95)', 'p(99)', 'count'],
}

export function runCreate() {
    const vu = __VU;
    const iter = __ITER;

    client.connect(grpcUrl, {plaintext: true});

    const payload = {
        name: `user_${iter}`,
        email: `mail${vu}_${iter}@example.com`,
        password: "password"
    };

    const response = create(client, payload)

    check(response, {
        'create_user status is OK': (r) => r && r.status === grpc.StatusOK,
    })

    if (usersCred.name === "") {
        usersCred.name = payload.name
        usersCred.email = payload.email
        usersCred.password = payload.password
    }
}

export function runList() {
    const vu = __VU;

    client.connect(grpcUrl, {plaintext: true});

    const payload = {page_token: tokens[vu]};

    const response = list(client, payload)

    check(response, {
        'list_users status is OK': (r) => r && r.status === grpc.StatusOK,
    })

    if (response && response.status === grpc.StatusOK) {
        tokens[vu] = JSON.parse(JSON.stringify(response.message)).nextPageToken
    }
}

export function runGetByEmail() {
    client.connect(grpcUrl, {plaintext: true});

    const payload = {
        email: usersCred.email,
    };

    const response = getByEmail(client, payload)

    check(response, {
        'create_user status is OK': (r) => r && r.status === grpc.StatusOK,
    })
}

export function runVerifyCred() {
    client.connect(grpcUrl, {plaintext: true});

    const payload = {
        email: usersCred.email,
        password: usersCred.password
    };

    const response = verifyCred(client, payload)

    check(response, {
        'create_user status is OK': (r) => r && r.status === grpc.StatusOK,
    })
}

export function runHealthz() {
    const response = healthz(diagnosticUrl)

    check(response, {
        'healthz is 200': (r) => r.status === 200,
    })
}