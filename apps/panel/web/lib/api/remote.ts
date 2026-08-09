import { request } from './client';
import type {
  RemoteCredential,
  RemoteCredentialPayload,
  RemoteCredentialUpdatePayload,
  RemoteAccessDefault,
  RemoteProtocol,
  RemoteSession,
  RemoteSessionPayload
} from '@/lib/types/remote';

export function getRemoteCredentials(accessToken: string, tenantId: string | null, protocol: RemoteProtocol) {
  return request<RemoteCredential[]>(`/remote/credentials?protocol=${encodeURIComponent(protocol)}`, {accessToken, tenantId});
}

export function getRemoteAccessDefaults(accessToken: string, tenantId: string | null) {
  return request<RemoteAccessDefault[]>('/remote/defaults', {accessToken, tenantId});
}

export function setRemoteAccessDefault(accessToken: string, tenantId: string | null, protocol: RemoteProtocol, accessPathId: string) {
  return request<RemoteAccessDefault>(`/remote/defaults/${protocol}`, {method: 'PUT', accessToken, tenantId, body: {accessPathId}});
}

export function createRemoteCredential(accessToken: string, tenantId: string | null, payload: RemoteCredentialPayload) {
  return request<RemoteCredential>('/remote/credentials', {
    method: 'POST',
    accessToken,
    tenantId,
    body: payload
  });
}

export function updateRemoteCredential(accessToken: string, tenantId: string | null, credentialId: string, payload: RemoteCredentialUpdatePayload) {
  return request<RemoteCredential>(`/remote/credentials/${credentialId}`, {
    method: 'PATCH',
    accessToken,
    tenantId,
    body: payload
  });
}

export function deleteRemoteCredential(accessToken: string, tenantId: string | null, credentialId: string) {
  return request<{id: string}>(`/remote/credentials/${credentialId}`, {
    method: 'DELETE',
    accessToken,
    tenantId
  });
}

export function createRemoteSession(accessToken: string, tenantId: string | null, payload: RemoteSessionPayload) {
  return request<RemoteSession>('/remote/sessions', {
    method: 'POST',
    accessToken,
    tenantId,
    body: payload
  });
}
