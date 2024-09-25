package app

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	val "github.com/go-playground/validator/v10"
	"github.com/go-programming-tour-book/blog-service/pkg/errcode"
)

type Response struct {
	Ctx *gin.Context
}

type Pager struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	TotalRows int `json:"total_rows"`
}

type StandardResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{Ctx: ctx}
}

func (r *Response) ToResponse(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	r.Ctx.JSON(http.StatusOK, StandardResponse{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

func (r *Response) ToResponseList(list interface{}, totalRows int) {
	r.Ctx.JSON(http.StatusOK, StandardResponse{
		Code: 0,
		Msg:  "success",
		Data: gin.H{
			"list": list,
			"pager": Pager{
				Page:      GetPage(r.Ctx),
				PageSize:  GetPageSize(r.Ctx),
				TotalRows: totalRows,
			},
		},
	})
}

func (r *Response) ToErrorResponse(err error) {
	var response StandardResponse
	var statusCode int
	switch e := err.(type) {

	case *errcode.Error:
		response = StandardResponse{
			Code: e.Code(),
			Msg:  e.Msg(),
			Data: nil,
		}
		details := e.Details()
		if len(details) > 0 {
			response.Data = gin.H{"details": details}
		}

		statusCode = e.StatusCode()
	default:
		// 处理非 errcode.Error 类型的错误（如内部服务器错误）
		serverError := errcode.ServerError
		response = StandardResponse{
			Code: serverError.Code(),
			Msg:  serverError.Msg(),
			Data: gin.H{"error": e.Error()},
		}
		statusCode = serverError.StatusCode()
	}
	fmt.Printf("Responding with: Status Code: %d, Response: %+v\n", statusCode, response)
	r.Ctx.JSON(statusCode, response)
}

type ValidError struct {
	Key     string
	Message string
}

type ValidErrors []*ValidError

func BindAndValid(c *gin.Context, v interface{}) (bool, ValidErrors) {
	var errs ValidErrors
	err := c.ShouldBind(v)
	if err != nil {
		v := c.Value("trans")
		trans, _ := v.(ut.Translator)
		verrs, ok := err.(val.ValidationErrors)
		if !ok {
			return false, errs
		}

		for key, value := range verrs.Translate(trans) {
			errs = append(errs, &ValidError{
				Key:     key,
				Message: value,
			})
		}

		return false, errs
	}

	return true, nil
}

func (v *ValidError) Error() string {
	return v.Message
}

func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}
