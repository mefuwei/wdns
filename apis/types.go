package apis

import (
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"net/http"
)

func SuccessResp(r *restful.Request, w *restful.Response, data interface{})  {
	resp := Resp{
		r:       r,
		w:       w,
		Data:    data,
	}
	resp.write(w, http.StatusOK)
}

func FailedResp(r *restful.Request, w *restful.Response, code int, message string)  {
	resp := Resp{
		r:       r,
		w:       w,
		Data:    nil,
	}
	resp.Failed(w, code, message)
}

type Resp struct {
	r *restful.Request
	w *restful.Response

	Title   string      `json:"title"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (r *Resp) Failed(w *restful.Response, code int, message string) {
	switch code {
	case http.StatusUnauthorized:
		r.Title = "请求未认证"
		if message == "" {
			r.Message = "您的请求未认证，请重新认证"
		}
	case http.StatusForbidden:
		r.Title = "请求无权限"
		if message == "" {
			r.Message = "您的请求无权限，请联系管理员加权限或取消访问"
		}
	case http.StatusInternalServerError:
		r.Title = "服务器内部错误"
		if message == "" {
			r.Message = "服务器内部错误，请联系管理员处理"
		}
	case http.StatusNotFound:
		r.Title = "未找到资源"
		if message == "" {
			r.Message = "未找到你要的资源"
		}
	case http.StatusBadRequest:
		r.Title = "错误的请求"
		if message == "" {
			r.Message = "你的请求参数错误"
		}
	}
	r.write(w, code)
}

func (r *Resp) write(w *restful.Response, code int) {
	if err := w.WriteHeaderAndEntity(code, r); err != nil {
		glog.Errorf("write msg to restful.Response failed, %s", err.Error())
	}
}
