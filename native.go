// gowinsvc project native.go
package gowinsvc

import (
	"syscall"
	"unsafe"
)

// about advapi32.dll
var (
	advapi32dll                      = syscall.NewLazyDLL("advapi32.dll")
	procStartServiceCtrlDispatcher   = advapi32dll.NewProc("StartServiceCtrlDispatcherW")
	procRegisterServiceCtrlHandlerEx = advapi32dll.NewProc("RegisterServiceCtrlHandlerExW")
	procSetServiceStatus             = advapi32dll.NewProc("SetServiceStatus")
)

// about kernel32.dll
var (
	kernel32dll           = syscall.NewLazyDLL("kernel32.dll")
	procOutputDebugString = kernel32dll.NewProc("OutputDebugStringW")
)

// StartServiceCtrlDispatcher Wrapper
func startServiceCtrlDispatcher(serviceName string, fnServiceMain interface{}) bool {
	svcTableEntry := new([2]SERVICE_TABLE_ENTRY)
	svcTableEntry[0].serviceName = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(serviceName)))
	svcTableEntry[0].serviceProc = syscall.NewCallback(fnServiceMain)
	svcTableEntry[1].serviceName = 0
	svcTableEntry[1].serviceProc = 0

	ret, _, _ := procStartServiceCtrlDispatcher.Call(uintptr(unsafe.Pointer(&svcTableEntry[0])))
	return ret != 0
}

func registerServiceCtrlHandlerEx(serviceName string, handler interface{}) SERVICE_STATUS_HANDLE {
	ret, _, _ := procRegisterServiceCtrlHandlerEx.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(serviceName))),
		syscall.NewCallback(handler),
		0,
	)

	return SERVICE_STATUS_HANDLE(ret)
}

// SetServiceStatus Wrapper
func setServiceStatusFunction(serviceStatusHandle SERVICE_STATUS_HANDLE, serviceStatus *SERVICE_STATUS) {
	procSetServiceStatus.Call(
		uintptr(serviceStatusHandle),
		uintptr(unsafe.Pointer(serviceStatus)),
	)
}

// OutputDebugString Wrapper
func outputDebugString(logString string) {
	utf16LogString := syscall.StringToUTF16Ptr(logString)
	procOutputDebugString.Call(uintptr(unsafe.Pointer(utf16LogString)))
}
