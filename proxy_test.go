package main

import (
	"github.com/stretchr/testify/mock"
)

// Mock

type ProxyMock struct{
	mock.Mock
}

func (m *ProxyMock) Provision(host, reconfPort, certPath, scAddress string) error {
	args := m.Called(host, certPath, scAddress)
	return args.Error(0)
}

func (m *ProxyMock) Reconfigure(host, reconfPort, serviceName string, servicePath []string) error {
	args := m.Called(host, reconfPort, serviceName, servicePath)
	return args.Error(0)
}


func getProxyMock(skipMethod string) *ProxyMock {
	mockObj := new(ProxyMock)
	if skipMethod != "Provision" {
		mockObj.On("Provision", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	}
	if skipMethod != "Reconfigure" {
		mockObj.On("Reconfigure", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	}
	return mockObj
}