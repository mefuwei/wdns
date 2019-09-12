// restful api add delete update select

package apis

import (
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
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
func List(r *restful.Request, w *restful.Response)  {
	getStorage()
}

func Get(r *restful.Request, w *restful.Response)  {
	//bs := getStorage()
	//records, err := bs.Get(name, qtype)
}

func Add(r *restful.Request, w *restful.Response)  {
	records := []storage.Record
	err := r.ReadEntity(&records)
	if err != nil {
		FailedResp(r, w, http.StatusBadRequest, err.Error())
		return
	}

	s := getStorage()
	if err := s.Set(records); err != nil {
		FailedResp(r, w, http.StatusInternalServerError, err.Error())
	}

	SuccessResp(r, w, nil)
	return
}

func Update(r *restful.Request, w *restful.Response)  {

}

func Delete(r *restful.Request, w *restful.Response)  {

}

// get storage of the config object.
func getStorage() storage.Storage {
	bs := storage.GetStorage(storageType, redisAddr, redisPasswd, redisDb)
	return bs
}



