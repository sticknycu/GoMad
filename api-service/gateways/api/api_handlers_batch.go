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
		id, alreadyExists, _ := api.storage.Save(products[i])
		if alreadyExists {
			errors = append(errors, "object already exist")
		}
		savedProducts = append(savedProducts, id)
	}

	_ = resp.WriteAsJson(map[string][]string{
		"ids objects inserted": savedProducts,
	})
}

func (api *API) getProductMemoryBatch(req *restful.Request, resp *restful.Response) {
	var ids []string
	body, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		log.Errorf("Failed to read ids, err=%v", err)
		_ = resp.WriteError(http.StatusBadRequest, err)
		return
	}
	err = json.Unmarshal(body, &ids)
	if err != nil {
		log.Errorf("Failed to unmarshal, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, err)
		return
	}

	var errors []string
	var products []domain.Product

	for i := range ids {
		id, notExist, _ := api.storage.Get(ids[i])
		if notExist {
			errors = append(errors, "object not exist")
		}
		products = append(products, id)
	}

	_ = resp.WriteAsJson(map[string][]domain.Product{
		"products:": products,
	})
}

func (api *API) updateProductMemoryBatch(req *restful.Request, resp *restful.Response) {
	panic("TODO")
}

func (api *API) deleteProductMemoryBatch(req *restful.Request, resp *restful.Response) {
	var ids []string
	body, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		log.Errorf("Failed to read ids, err=%v", err)
		_ = resp.WriteError(http.StatusBadRequest, err)
		return
	}
	err = json.Unmarshal(body, &ids)
	if err != nil {
		log.Errorf("Failed to unmarshal, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, err)
		return
	}

	var errors []string
	var products []string

	for i := range ids {
		deleted, _ := api.storage.Delete(ids[i])
		if !deleted {
			errors = append(errors, "object not exist")
		} else {
			products = append(products, ids[i])
		}
	}

	_ = resp.WriteAsJson(map[string][]string{
		"Ids objects deleted:": products,
	})
}
