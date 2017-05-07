package httpRender

import "net/http"

func (t *renderImpl) JSON(obj interface{}) {
	t.WriteResponse(nil, http.StatusOK, obj)
}
