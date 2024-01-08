package results

import (
	"net/http"
	"sync"

	"github.com/retail-ai-inc/beanq/helper/json"
	"github.com/retail-ai-inc/beanqui/internal/routers/consts"
)

type Result struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func (t *Result) reset() {
	t.Data = nil
	t.Msg = consts.SuccessMsg
	t.Code = consts.SuccessCode
}

func (t *Result) Json(w http.ResponseWriter, httpCode int) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	b, err := json.Marshal(t)
	if err != nil {
		return err
	}

	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}

var resultPool = sync.Pool{New: func() any {
	return &Result{
		Code: consts.SuccessCode,
		Msg:  consts.SuccessMsg,
		Data: nil,
	}
}}

func Get() (*Result, func()) {

	result := resultPool.Get().(*Result)
	return result, func() {
		result.reset()
		resultPool.Put(result)
	}
}
