import { type AuthRequest, type ErrorMessage, type LoginResponse } from "./models";

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
