import http from 'k6/http';
import { sleep } from 'k6';

export function healthz(url) {
    const response =  http.get(url + '/healthz');

    sleep(3)

    return response
}