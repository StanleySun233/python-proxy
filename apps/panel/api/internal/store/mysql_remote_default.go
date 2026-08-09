package store

import (
	"database/sql"

	"github.com/StanleySun233/python-proxy/apps/panel/api/internal/domain"
)

func (s *MySQLStore) ListRemoteAccessDefaults(tenantID string) []domain.RemoteAccessDefault {
	rows, err := s.db.Query(remoteAccessDefaultSelect()+" WHERE tenant_id = ? ORDER BY remote_protocol", tenantID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	items := make([]domain.RemoteAccessDefault, 0, 2)
	for rows.Next() {
		var item domain.RemoteAccessDefault
		if err := rows.Scan(&item.TenantID, &item.Protocol, &item.AccessPathID, &item.UpdatedBy, &item.UpdatedAt); err == nil {
			items = append(items, item)
		}
	}
	return items
}

func (s *MySQLStore) RemoteAccessDefault(tenantID string, protocol string) (domain.RemoteAccessDefault, bool) {
	var item domain.RemoteAccessDefault
	err := s.db.QueryRow(remoteAccessDefaultSelect()+" WHERE tenant_id = ? AND remote_protocol = ?", tenantID, protocol).
		Scan(&item.TenantID, &item.Protocol, &item.AccessPathID, &item.UpdatedBy, &item.UpdatedAt)
	return item, err == nil
}

func (s *MySQLStore) SetRemoteAccessDefault(tenantID string, protocol string, accessPathID string, updatedBy string) (domain.RemoteAccessDefault, error) {
	now := nowRFC3339()
	_, err := s.db.Exec(
		`INSERT INTO tenant_remote_access_defaults (tenant_id, remote_protocol, access_path_id, updated_by, updated_at)
		 VALUES (?, ?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE access_path_id = VALUES(access_path_id), updated_by = VALUES(updated_by), updated_at = VALUES(updated_at)`,
		tenantID, protocol, accessPathID, updatedBy, now,
	)
	if err != nil {
		return domain.RemoteAccessDefault{}, err
	}
	item, ok := s.RemoteAccessDefault(tenantID, protocol)
	if !ok {
		return domain.RemoteAccessDefault{}, sql.ErrNoRows
	}
	return item, nil
}

func remoteAccessDefaultSelect() string {
	return "SELECT tenant_id, remote_protocol, access_path_id, updated_by, updated_at FROM tenant_remote_access_defaults"
}
