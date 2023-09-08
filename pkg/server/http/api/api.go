package api

import (
	echo "github.com/labstack/echo/v4"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *string     `json:"error,omitempty"`
}

type Factory interface {
	Register(model RegisterModel, schemaMap RegisterSchemaMap) echo.HandlerFunc
	Login(model LoginModel, schemaMap LoginSchemaMap) echo.HandlerFunc
}

type FactoryImplementation struct{}

func New() *FactoryImplementation {
	return &FactoryImplementation{}
}
