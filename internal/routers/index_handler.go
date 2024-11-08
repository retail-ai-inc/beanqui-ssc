package routers

type Index struct {
}

func NewIndex() *Index {
	return &Index{}
}

func (t *Index) Home(ctx *BeanContext) error {
	return nil
}
