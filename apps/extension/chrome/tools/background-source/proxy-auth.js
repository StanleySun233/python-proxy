import { appendLog } from './diagnostics.js';
import { accessPathById, selectedEntrypoint } from './state.js';
import { isUsableAccessPath } from './routing.js';

let proxyAuthCache = {
  targets: new Set(),
  token: ''
};

function proxyTargetKey(host, port) {
  return `${String(host || '').toLowerCase()}:${Number(port || 0)}`;
}

function proxyAuthTargetsFrom(state) {
  return new Set((state.remote.accessPaths || [])
    .map((path) => accessPathById(state, path.id))
    .filter(isUsableAccessPath)
    .map((path) => selectedEntrypoint(path))
    .filter(Boolean)
    .map((entrypoint) => proxyTargetKey(entrypoint.host, entrypoint.port)));
}

export function updateProxyAuthCache(state) {
  proxyAuthCache = {
    targets: proxyAuthTargetsFrom(state),
    token: state.session && state.session.proxyToken ? String(state.session.proxyToken) : ''
  };
}

function matchesProxyChallenge(details) {
  if (!details || !details.isProxy || !proxyAuthCache.token) {
    return false;
  }
  const challenger = details.challenger || {};
  return proxyAuthCache.targets.has(proxyTargetKey(challenger.host, challenger.port));
}

export function registerProxyAuthHandler() {
  if (!chrome.webRequest || !chrome.webRequest.onAuthRequired) {
    return;
  }
  chrome.webRequest.onAuthRequired.addListener(
    (details, callback) => {
      if (!matchesProxyChallenge(details)) {
        callback({});
        return;
      }
      const challenger = details.challenger || {};
      appendLog('info', 'proxy_auth_supplied', {
        host: String(challenger.host || ''),
        port: Number(challenger.port || 0)
      }).catch(() => {});
      callback({
        authCredentials: {
          username: 'token',
          password: proxyAuthCache.token
        }
      });
    },
    { urls: ['<all_urls>'] },
    ['asyncBlocking']
  );
}
