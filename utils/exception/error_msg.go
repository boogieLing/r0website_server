package exception

import "fmt"

type SysError struct {
	msg string
}

func (e SysError) Error() string {
	return fmt.Sprintf("发生异常：%s", e.msg)
}

func NewSysError(msg string) SysError {
	return SysError{
		msg: msg,
	}
}
