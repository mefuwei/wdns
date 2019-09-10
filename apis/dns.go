// restful api add delete update select

package apis

import (
	"github.com/emicklei/go-restful"
	"github.com/mefuwei/dns/storage"
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

}

func Update(r *restful.Request, w *restful.Response)  {

}

func Delete(r *restful.Request, w *restful.Response)  {

}

func getStorage() storage.Storage {
	bs := storage.GetStorage(storageType, redisAddr, redisPasswd, redisDb)
	return bs
}



