// restful api add delete update select

package apis

/*
	RecordAdd	添加记录 POST /api/v1/domain/{id}/records
	RecordDel	删除记录 DELETE /api/v1/domain/{id}/records/{rid}
	RecordList	列出记录 GET /api/v1/domain/{id}/records
 */

import (
	"github.com/emicklei/go-restful"
	"github.com/mefuwei/dns/storage"
	"net/http"
)

var (
	// TODO used config
	storageType = "redis"
	redisAddr = "localhost:6379"
	redisPasswd = ""
	redisDb = 1
)

// todo list rewrite
func DnsList(r *restful.Request, w *restful.Response)  {
	getStorage()
}

func DnsGet(r *restful.Request, w *restful.Response)  {
	//bs := getStorage()
	//records, err := bs.Get(name, qtype)
}

func DnsAdd(r *restful.Request, w *restful.Response)  {

	records := make([]storage.Record, 10)

	err := r.ReadEntity(&records)
	if err != nil {
		FailedResp(r, w, http.StatusBadRequest, err.Error())
		return
	}

	s := getStorage()
	if err := s.Set(records); err != nil {
		if err == storage.RecordParamsNullFailed {
			FailedResp(r, w, http.StatusBadRequest, err.Error())
		}
		FailedResp(r, w, http.StatusInternalServerError, err.Error())
	}

	SuccessResp(r, w, nil)
	return
}

func DnsUpdate(r *restful.Request, w *restful.Response)  {

}

func DnsDelete(r *restful.Request, w *restful.Response)  {

}

// get storage of the config object.
func getStorage() storage.Storage {
	bs := storage.GetStorage(storageType, redisAddr, redisPasswd, redisDb)
	return bs
}



