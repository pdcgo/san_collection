package san_mcp

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type McpRegisterHandler func(server *mcp.Server) error
type McpServer struct {
	handler McpRegisterHandler
}

func (m *McpServer) Run() error {
	var err error

	// creating server
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// creating mcp server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "Developer MCP Tool Server PDC",
		Version: "0.1.0",
	}, nil)

	// registering mcp
	err = m.handler(server)
	if err != nil {
		return err
	}

	mcpHandler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{
		DisableLocalhostProtection: true,
	})

	mux.Handle("/", stripOrigin(mcpHandler))

	// serve webserver
	httpSrv := &http.Server{
		Addr:              "localhost:8090",
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
	slog.Info("running on localhost:8090")
	return httpSrv.ListenAndServe()
}

func NewMcpServer(handler McpRegisterHandler) *McpServer {
	return &McpServer{handler}
}

func stripOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Del("Origin")
		next.ServeHTTP(w, r)
	})
}
