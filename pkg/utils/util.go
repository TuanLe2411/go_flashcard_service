package utils

import (
	"flashcard_service/pkg/constant"
	"net/http"
	"strconv"
	"fmt"
)

func ChainMiddlewares(handler http.Handler, middlewares ...constant.Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i].Do(handler)
	}
	return handler
}

// Int64ToString chuyển đổi int64 sang string sử dụng strconv
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

// StringToInt64 chuyển đổi string sang int64, trả về giá trị và lỗi (nếu có)
func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// Int64ToStringFmt chuyển đổi int64 sang string sử dụng fmt.Sprintf
// Phương pháp này có thể đơn giản hơn nhưng chậm hơn strconv
func Int64ToStringFmt(num int64) string {
	return fmt.Sprintf("%d", num)
}
