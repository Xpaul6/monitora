import type { AddServerRequest, AuthRequest, Component, DeleteLimitRequest, DeleteServerRequest, ErrorMessage, GetAllServersResponse, GetLimitsRequest, GetNotificationsRequest, GetServerComponentsRequest, GetStatsByPeriodRequest, GetStatsByPeriodResponse, Limit, LimitNotification, LoginResponse, MetricType, SetLimitRequest } from "./models";

// Service functions
function getToken() {
  return localStorage.getItem("token") || "";
}

export function getHeaders(): HeadersInit {
  return {
    "Content-Type": "application/json",
    Authorization: `Bearer ${getToken()}`,
  };
}

// API functions
export async function checkAuth(): Promise<boolean> {
  const res = await fetch('/api/auth/check', { headers: getHeaders() });
  if (!res.ok) return false;
  return true;
}

export async function register(reqBody: AuthRequest): Promise<null | string> {
  const res = await fetch('/api/auth/register', {
    method: "POST",
    headers: getHeaders(),
    body: JSON.stringify(reqBody)
  });
  if (res.ok) return null;

  const data: ErrorMessage = await res.json();
  return data.error;
}

export async function login(reqBody: AuthRequest): Promise<LoginResponse | string> {
  const res = await fetch('/api/auth/login', {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(reqBody)
  });
  if (res.ok) return res.json();

  const data: ErrorMessage = await res.json();
  return data.error;
}

export async function getAllServers(): Promise<GetAllServersResponse | string> {
  const res = await fetch('/api/user/servers', {
    method: 'GET',
    headers: getHeaders()
  });
  if (res.ok) return res.json();

  const data: ErrorMessage = await res.json();
  return data.error;
}

export async function addServer(reqBody: AddServerRequest): Promise<null | string> {
  const res = await fetch('/api/user/add-server', {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(reqBody)
  });
  if (res.ok) return null;

  const data: ErrorMessage = await res.json();
  return data.error;
}

export async function getServerComponents(reqBody: GetServerComponentsRequest): Promise<Component[] | string> {
  const res = await fetch('/api/user/components', {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(reqBody)
  });
  if (res.ok) return res.json();

  const data: ErrorMessage = await res.json();
  return data.error;
}

export async function getMetricTypes(): Promise<MetricType[] | string> {
  const res = await fetch('/api/service/metric-types', {
    method: 'GET'
  });
  if (res.ok) return res.json();

  const data: ErrorMessage = await res.json();
  return data.error;
}

export async function getStatsByPeriod(reqBody: GetStatsByPeriodRequest): Promise<GetStatsByPeriodResponse[] | string> {
  const res = await fetch('/api/stats/by-period', {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(reqBody)
  });
  if (res.ok) return res.json();

  const data: ErrorMessage = await res.json();
  return data.error;
}

export async function deleteServer(reqBody: DeleteServerRequest): Promise<null | string> {
  const res = await fetch('api/user/delete-server', {
    method: 'DELETE',
    headers: getHeaders(),
    body: JSON.stringify(reqBody)
  });
  if (res.ok) return null;

  const data = await res.json();
  return data.error;
}

export async function getLimits(reqBody: GetLimitsRequest): Promise<Limit[] | string> {
  const res = await fetch('api/user/limits', {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(reqBody)
  });
  if (res.ok) return res.json();

  const data = await res.json();
  return data.error;
}

export async function setLimit(reqBody: SetLimitRequest): Promise<null | string> {
  const res = await fetch('api/user/set-limit', {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(reqBody)
  });
  if (res.ok) return res.json();

  const data = await res.json();
  return data.error;
}

export async function deleteLimit(reqBody: DeleteLimitRequest): Promise<null | string> {
  const res = await fetch('api/user/delete-limit', {
    method: 'DELETE',
    headers: getHeaders(),
    body: JSON.stringify(reqBody)
  });
  if (res.ok) return null;

  const data = await res.json();
  return data.error;
}

export async function getNotifications(reqBody: GetNotificationsRequest): Promise<LimitNotification[] | string> {
  const res = await fetch('api/user/notifications', {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(reqBody)
  });
  if (res.ok) return res.json();

  const data = await res.json();
  return data.error;
}
