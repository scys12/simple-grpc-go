package gateway

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/scys12/simple-grpc-go/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

func OpenAPIServer(dir string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, ".swagger.json") {
			logger.Log.Error("Not Found:", zap.String("url", r.URL.Path))
			http.NotFound(rw, r)
			return
		}

		logger.Log.Info("Serving -> ", zap.String("url", r.URL.Path))
		p := strings.TrimPrefix(r.URL.Path, "/openapiv2/")
		p = path.Join(dir, p)
		http.ServeFile(rw, r, p)
	}
}

func CheckHealth(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if s := conn.GetState(); s != connectivity.Ready {
			http.Error(w, fmt.Sprintf("grpc server is %s", s), http.StatusBadGateway)
			return
		}
		fmt.Fprintln(w, "ok")
	}
}
