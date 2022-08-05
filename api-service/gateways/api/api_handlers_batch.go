package api

import (
	"encoding/json"
	"exam-api/domain"
	"github.com/emicklei/go-restful/v3"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func (api *API) createProductMemoryBatch(req *restful.Request, resp *restful.Response) {
	var products []domain.Product
	body, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		log.Errorf("Failed to read product, err=%v", err)
		_ = resp.WriteError(http.StatusBadRequest, err)
		return
	}
	err = json.Unmarshal(body, &products)
	if err != nil {
		log.Errorf("Failed to unmarshal, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, err)
		return
	}

	var errors []string
	var savedProducts []string

	for i := range products {
		id, alreadyExists, err := api.storage.Save(products[i])
		if err != nil {
			errors = append(errors, "internal error")
		}
		if alreadyExists {
			errors = append(errors, "object already exist")
		}
		savedProducts = append(savedProducts, id)
	}

	_ = resp.WriteAsJson(map[string][]string{
		"id object inserted": savedProducts,
	})
}

func (api *API) getProductMemoryBatch(req *restful.Request, resp *restful.Response) {
	panic("TODO")
}

func (api *API) updateProductMemoryBatch(req *restful.Request, resp *restful.Response) {
	panic("TODO")
}

func (api *API) deleteProductMemoryBatch(req *restful.Request, resp *restful.Response) {
	panic("TODO")
}
