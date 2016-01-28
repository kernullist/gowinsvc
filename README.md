# gowinsvc
Go 로 Windows Service 를 쉽게 만들게 도와주는 패키지

# Demo
[![Demo Video](http://img.youtube.com/vi/raFxSoF6aU4/0.jpg)](http://www.youtube.com/watch?v=raFxSoF6aU4)

# Setup
1. mingw 설치

[http://sourceforge.net/projects/mingw-w64/](http://sourceforge.net/projects/mingw-w64/)

2. gowinsvc 패키지 다운로드

> go get github.com/kernullist/gowinsvc

# Example
1초마다 현재 시간을 출력하는 서비스 예제

```go
package main

import (
	"time"

	"github.com/kernullist/gowinsvc"
)

type MySvc struct {
	service *gowinsvc.ServiceObject
}

func (mysvc MySvc) Serve(serviceExit <-chan bool) {
	for {
		select {
		case <-serviceExit:
			mysvc.service.OutputDebugString("[MYSERVICE] My Service Exit~~~\n")
			return
		case <-time.After(1 * time.Second):
			mysvc.service.OutputDebugString(
				"[MYSERVICE] Now : %d:%d:%d\n",
				time.Now().Hour(),
				time.Now().Minute(),
				time.Now().Second())
		}
	}
}

func main() {
	mysvc := new(MySvc)
	mysvc.service = gowinsvc.NewService("myservice")
	mysvc.service.StartServe(mysvc)
}
```
