package proxyhttpapi

import (
	"net/http"

	"github.com/StanleySun233/python-proxy/apps/panel/api/internal/httpctx"
)

func (r *Router) handleTopology(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeMethodNotAllowed(w, "GET")
		return
	}
	tenantCtx, ok := httpctx.TenantAuth(req.Context())
	if !ok {
		writeError(w, http.StatusBadRequest, "tenant_required")
		return
	}
	writeSuccess(w, http.StatusOK, r.service.Topology(tenantCtx))
}
