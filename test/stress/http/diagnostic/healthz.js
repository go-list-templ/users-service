import http from 'k6/http';
import {check, sleep} from 'k6';

export function healthCheckTest(url) {
    const endpoint = url + '/healthz'
    const res = http.get(endpoint);

    check(res, {
        'http status is 200': (r) => r.status === 200,
    });

    sleep(5);
}