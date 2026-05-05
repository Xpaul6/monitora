// General models
export type Server = {
  ID: number,
  user_id: number,
  name: string,
  ip: string,
  status: string
}

export type Component = {
  ID: number,
  server_id: number,
  type: string,
  address: string
}

export type MetricType = {
  ID: number,
  name: string,
  unit: string,
  description: string
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

export type GetServerComponentsRequest = {
  id: number
}

export type DeleteServerRequest = {
  id: number
}

export type GetStatsByPeriodRequest = {
  server_id: number,
  period_begin: string,
  period_end: string
}

export type GetStatsByPeriodResponse = {
  component: Component,
  metric_type: MetricType,
  value: number,
  timestamp: Date
}
