import http from 'k6/http';

export function healthz(url) {
    return http.get(url + '/healthz')
}