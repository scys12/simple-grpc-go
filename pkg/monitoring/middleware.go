package monitoring

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pth := r.URL.Host
		method := r.Method

		timer := prometheus.NewTimer(latency.WithLabelValues(pth, method))

		w.WriteHeader(http.StatusOK)
		next.ServeHTTP(w, r)

		responseStatus.WithLabelValues(strconv.Itoa(http.StatusOK), pth, method).Inc()
		totalRequests.WithLabelValues(pth, method).Inc()
		timer.ObserveDuration()
	})
}
