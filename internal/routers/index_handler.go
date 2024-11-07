package routers

import (
	"net/http"
)

type Index struct {
}

func NewIndex() *Index {
	return &Index{}
}

func (t *Index) File(w http.ResponseWriter, r *http.Request) {

}
