// gowinsvc project gowinsvc.go
package gowinsvc

import (
	"C"
	"fmt"
	"os"
	"sync"
)

// Create New ServiceObject Function
func NewService(serviceName string) *ServiceObject {
	serviceObject := ServiceObject{}
	serviceObject.serviceName = serviceName
	serviceObject.serviceStatusHandle = 0
	serviceObject.serviceExit = make(chan bool)
	return &serviceObject
}

// Serve Function
func (service *ServiceObject) StartServe(serviceInterface Service) bool {
	if service.serviceStatusHandle != 0 {
		return false
	}

	service.serviceInterface = serviceInterface

	return startServiceCtrlDispatcher(service.serviceName, service.serviceMain)
}

// Internal Service Main Function
func (service *ServiceObject) serviceMain(args uint32, ppServiceArgsVectors uintptr) uintptr {
	if service.initService() == false {
		return 0
	}

	defer os.Exit(0)
	defer service.setServiceStatus(SERVICE_STOPPED)

	var wg sync.WaitGroup

	wg.Add(1)

	go func(service *ServiceObject) {
		defer wg.Done()
		service.serviceInterface.Serve(service.serviceExit)
	}(service)

	wg.Wait()

	return 0
}

// Service Handler Function
func (service *ServiceObject) serviceHandler(control, eventType uint32, eventData, context uintptr) uintptr {
	if control == SERVICE_CONTROL_STOP {
		service.setServiceStatus(SERVICE_STOP_PENDING)
		service.serviceExit <- true
	}

	return 0
}

// Initialize Service Function
func (service *ServiceObject) initService() bool {
	service.serviceStatusHandle = registerServiceCtrlHandlerEx(service.serviceName, service.serviceHandler)
	if service.serviceStatusHandle == 0 {
		return false
	}

	service.currentServiceStatus.serviceType = SERVICE_WIN32_OWN_PROCESS
	service.currentServiceStatus.controlsAccepted = SERVICE_ACCEPT_STOP
	service.currentServiceStatus.checkPoint = 0
	service.currentServiceStatus.serviceSpecificExitCode = 0
	service.currentServiceStatus.win32ExitCode = 0
	service.setServiceStatus(SERVICE_RUNNING)

	return true
}

// Change Service Status Function
func (service *ServiceObject) setServiceStatus(status uint32) {
	if service.serviceStatusHandle == 0 {
		return
	}

	service.currentServiceStatus.currentState = status
	setServiceStatusFunction(service.serviceStatusHandle, &service.currentServiceStatus)
}

// OutputDebugString Wrapper Function
func (service *ServiceObject) OutputDebugString(format string, a ...interface{}) {
	logString := fmt.Sprintf(format, a...)
	outputDebugString(logString)
}
