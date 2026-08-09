package store

import (
	"context"

	"github.com/StanleySun233/python-proxy/apps/panel/api/internal/domain"
	"github.com/StanleySun233/python-proxy/apps/panel/api/internal/store/deleteplan"
)

func (s *MySQLStore) ListNodeAccessPaths() []domain.NodeAccessPath {
	items, err := s.proxyRepository().listNodeAccessPaths(context.Background())
	if err != nil {
		return nil
	}
	return deriveNodeAccessPaths(items, s.ListChains(), s.ListNodes())
}

func (s *MySQLStore) ListNodeAccessPathsForTenant(tenantCtx domain.TenantAuthContext) []domain.NodeAccessPath {
	if tenantCtx.SuperAdmin && tenantCtx.ActiveTenant.TenantID == "" {
		return s.ListNodeAccessPaths()
	}
	items, err := s.proxyRepository().listNodeAccessPathsForTenant(context.Background(), tenantCtx)
	if err != nil {
		return nil
	}
	return deriveNodeAccessPaths(items, s.ListChainsForTenant(tenantCtx), s.ListNodesForTenant(tenantCtx))
}

func (s *MySQLStore) CreateNodeAccessPath(input domain.CreateNodeAccessPathInput) (domain.NodeAccessPath, error) {
	ownerID, err := s.defaultOwnerAccountID()
	if err != nil {
		return domain.NodeAccessPath{}, err
	}
	pathID, err := s.nextID("node_access_path")
	if err != nil {
		return domain.NodeAccessPath{}, err
	}
	item := domain.NodeAccessPath{
		ID:             pathID,
		CreateID:       ownerID,
		OwnerID:        ownerID,
		ChainID:        input.ChainID,
		Name:           input.Name,
		Mode:           input.Mode,
		Protocol:       input.Protocol,
		ServiceType:    input.ServiceType,
		RemoteProtocol: input.RemoteProtocol,
		ListenHost:     input.ListenHost,
		ListenPort:     input.ListenPort,
		TargetProtocol: input.TargetProtocol,
		TargetHost:     input.TargetHost,
		TargetPort:     input.TargetPort,
		TargetSNI:      input.TargetSNI,
		TLSMode:        input.TLSMode,
		AuthMode:       input.AuthMode,
		Options:        input.Options,
		Enabled:        true,
	}
	if err := s.proxyRepository().createNodeAccessPath(context.Background(), item, ""); err != nil {
		return domain.NodeAccessPath{}, err
	}
	return s.deriveNodeAccessPath(item), nil
}

func (s *MySQLStore) CreateNodeAccessPathForTenant(tenantCtx domain.TenantAuthContext, input domain.CreateNodeAccessPathInput) (domain.NodeAccessPath, error) {
	pathID, err := s.nextID("node_access_path")
	if err != nil {
		return domain.NodeAccessPath{}, err
	}
	item := domain.NodeAccessPath{
		ID:             pathID,
		CreateID:       tenantCtx.Account.ID,
		OwnerID:        tenantCtx.Account.ID,
		ChainID:        input.ChainID,
		Name:           input.Name,
		Mode:           input.Mode,
		Protocol:       input.Protocol,
		ServiceType:    input.ServiceType,
		RemoteProtocol: input.RemoteProtocol,
		ListenHost:     input.ListenHost,
		ListenPort:     input.ListenPort,
		TargetProtocol: input.TargetProtocol,
		TargetHost:     input.TargetHost,
		TargetPort:     input.TargetPort,
		TargetSNI:      input.TargetSNI,
		TLSMode:        input.TLSMode,
		AuthMode:       input.AuthMode,
		Options:        input.Options,
		Enabled:        true,
	}
	if err := s.proxyRepository().createNodeAccessPath(context.Background(), item, tenantCtx.ActiveTenant.TenantID); err != nil {
		return domain.NodeAccessPath{}, err
	}
	items := deriveNodeAccessPaths([]domain.NodeAccessPath{item}, s.ListChainsForTenant(tenantCtx), s.ListNodesForTenant(tenantCtx))
	return items[0], nil
}

func (s *MySQLStore) UpdateNodeAccessPath(pathID string, input domain.UpdateNodeAccessPathInput) (domain.NodeAccessPath, error) {
	item, err := s.proxyRepository().updateNodeAccessPath(context.Background(), pathID, input)
	if err != nil {
		return domain.NodeAccessPath{}, err
	}
	return s.deriveNodeAccessPath(item), nil
}

func (s *MySQLStore) deriveNodeAccessPath(path domain.NodeAccessPath) domain.NodeAccessPath {
	return deriveNodeAccessPaths([]domain.NodeAccessPath{path}, s.ListChains(), s.ListNodes())[0]
}

func (s *MySQLStore) DeleteNodeAccessPath(pathID string) error {
	plan, err := s.proxyRepository().buildNodeAccessPathDeletePlan(context.Background(), pathID, false)
	if err != nil {
		return err
	}
	_, err = deleteplan.NewMySQLExecutor(s.db).Execute(context.Background(), plan)
	return err
}

func (s *MySQLStore) NodeAccessPathBindingPermission(tenantCtx domain.TenantAuthContext, pathID string) (domain.BindingPermission, bool) {
	return s.tenantResourcePermission(tenantCtx, "tenant_access_paths", "access_path_id", pathID)
}

func (s *MySQLStore) CountNodeAccessPathBindings(pathID string) int {
	return s.countTenantResourceBindings("tenant_access_paths", "access_path_id", pathID)
}
