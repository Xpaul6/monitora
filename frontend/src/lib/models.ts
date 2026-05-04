// General models
export type Server = {
  user_id: number,
  name: string,
  ip: string,
  status: string
}

// API models
export type ErrorMessage = {
  error: string
}

export type AuthRequest = {
  email: string,
  password: string
}

export type LoginResponse = {
  id: number,
  token: string
}

export type GetAllServersResponse = {
  count: number,
  servers: Server[]
}

export type AddServerRequest = {
  name: string,
  ip: string
}
