package app

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// SuccessOrAbort 會判斷是否讓程式繼續向後執行，若 err 不存在則繼續執行
// 若 err 存在，則根據錯誤類型將錯誤訊息帶入 gin 的 ctx.Errors 中，並中止在此
// 最後到 middleware 的 error_handle 財會回傳錯誤訊息
func SuccessOrAbort(
	ctx *gin.Context,
	statusCode int,
	err error,
) bool {
	// 表示在 ctx 中已經就有錯誤存在
	if len(ctx.Errors) > 0 {
		return false
	}

	// 當錯誤是不合法的 JSON syntax 時
	var jsonSyntaxErr *json.SyntaxError
	if errors.As(err, &jsonSyntaxErr) {
		_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid JSON syntax"))
		return false
	}

	// 當錯誤是 request 中不正確的 params，把 ErrorType 設成 ErrorTypeBind
	var validationErrs validator.ValidationErrors
	if ok := errors.As(err, &validationErrs); ok {
		_ = ctx.AbortWithError(http.StatusBadRequest, validationErrs).SetType(gin.ErrorTypeBind)
		return false
	}

	// // 當錯誤是 GORM 找不到該資料
	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	_ = ctx.AbortWithError(http.StatusNotFound, errors.New("record not found"))
	// 	return false
	// }

	// // if error is PointClickCare Error
	// var pccErrors *pcc.ErrorResp
	// if ok := errors.As(err, &pccErrors); ok {
	// 	firstError := pccErrors.Errors[0]
	// 	statusCode, err := strconv.Atoi(firstError.Status)
	// 	if err != nil {
	// 		log.Warn("[pkg/app] SuccessOrAbort - strconv.Atoi failed", err)
	// 	}

	// 	_ = ctx.AbortWithError(statusCode, fmt.Errorf("%s:%s", firstError.Title, firstError.Detail))
	// 	return false
	// }

	// if pccErrors, ok := err.(*pcc.ErrorResp); ok {
	//	firstError := pccErrors.Errors[0]
	//	statusCode, err := strconv.Atoi(firstError.Status)
	//	if err != nil {
	//		log.Warn("[pkg/app] SuccessOrAbort - strconv.Atoi failed", err)
	//	}
	//
	//	_ = ctx.AbortWithError(statusCode, fmt.Errorf("%s:%s", firstError.Title, firstError.Detail))
	//	return false
	//}

	// 如果是其他類型的錯誤
	if err != nil {
		_ = ctx.AbortWithError(statusCode, err)
		return false
	}

	return true
}
