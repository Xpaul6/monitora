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
