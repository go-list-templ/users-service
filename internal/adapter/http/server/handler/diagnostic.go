package handler

import "net/http"

type Diagnostic struct{}

func RegisterDiagnostic() {
	d := &Diagnostic{}

	http.HandleFunc("/healthz", d.Health())
}

func (d *Diagnostic) Health() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
