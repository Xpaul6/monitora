package models

import "time"

type AuthRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	ID    uint   `json:"id"`
	Token string `json:"token"`
}

type GetAllServersResponse struct {
	Count   int64    `json:"count"`
	Servers []Server `json:"servers"`
}

type AddServerRequest struct {
	Name string `json:"name" binding:"required"`
	IP   string `json:"ip" binding:"required"`
}

type DeleteServerRequest struct {
	ID uint `json:"id" binding:"required"`
}

type GetStatsByPeriodRequest struct {
	ServerID    uint      `json:"server_id" binding:"required"`
	PeriodBegin time.Time `json:"period_begin" binding:"required"`
	PeriodEnd   time.Time `json:"period_end" binding:"required"`
}

type GetStatsByPeriodResponse struct {
	Component  Component  `json:"component"`
	MetricType MetricType `json:"metric_type"`
	Value      float64    `json:"value"`
	TimeStamp  time.Time  `json:"timestamp"`
}
