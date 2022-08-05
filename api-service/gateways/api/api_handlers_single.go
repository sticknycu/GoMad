package api

import (
	"encoding/json"
	"exam-api/domain"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (api *API) createProductMemorySingle(req *restful.Request, resp *restful.Response) {
	product := &domain.Product{}
	err := req.ReadEntity(product)
	if err != nil {
		log.Errorf("Failed to read product, err=%v", err)
		_ = resp.WriteError(http.StatusBadRequest, err)
		return
	}
	id, alreadyExists, err := api.storage.Save(*product)
	if err != nil {
		log.Errorf("Failed to save product in storage, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to save product"))
		return
	}

	if alreadyExists {
		log.Infof("Product %s already in store", id)
		_ = resp.WriteError(http.StatusConflict, fmt.Errorf("product already exists"))
		return
	}

	log.Infof("Product %s saved in store", id)

	_ = resp.WriteAsJson(map[string]string{
		"id": id,
	})
}

func (api *API) getProductMemorySingle(req *restful.Request, resp *restful.Response) {
	id := req.QueryParameter("id")
	if id == "" {
		log.Infof("No id provided in request")
		_ = resp.WriteError(http.StatusBadRequest, fmt.Errorf("id must be provided"))
		return
	}
	product, exists, err := api.storage.Get(id)
	if err != nil {
		log.Errorf("Failed to get product from storage, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to get product from store"))
		return
	}
	if !exists {
		log.Infof("Product %s not in store", id)
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}
	_ = resp.WriteAsJson(product)
}

func (api *API) updateProductMemorySingle(req *restful.Request, resp *restful.Response) {
	id := req.QueryParameter("id")
	if id == "" {
		log.Infof("No id provided in request")
		_ = resp.WriteError(http.StatusBadRequest, fmt.Errorf("id must be provided"))
		return
	}
	_, exists, err := api.storage.Get(id)
	if err != nil {
		log.Errorf("Failed to get product from storage, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to get product from store"))
		return
	}
	if !exists {
		log.Infof("Product %s not in store", id)
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}

	var p domain.Product
	err = json.NewDecoder(req.Request.Body).Decode(&p)
	if err != nil {
		log.Infof("Product diff not found in body")
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("product diff not found"))
		return
	}

	ok, _ := api.storage.Update(id, p)
	if !ok {
		log.Errorf("Failed to update product in storage, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to update product from store"))
		return
	}

	log.Infof("Product %s updated in store", id)
	_ = resp.WriteAsJson(map[string]string{
		"successfully updated product with id:": id,
	})
	fmt.Println(api.storage)

}

func (api *API) deleteProductMemorySingle(req *restful.Request, resp *restful.Response) {
	id := req.QueryParameter("id")
	if id == "" {
		log.Infof("No id provided in request")
		_ = resp.WriteError(http.StatusBadRequest, fmt.Errorf("id must be provided"))
		return
	}
	_, exists, err := api.storage.Get(id)
	if err != nil {
		log.Errorf("Failed to get product from storage, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to get product from store"))
		return
	}
	if !exists {
		log.Infof("Product %s not in store", id)
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}

	ok, _ := api.storage.Delete(id)
	if !ok {
		log.Infof("Error on deleting process")
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("error on deleting process"))
		return
	}

	log.Infof("Product %s deleted from store", id)
	_ = resp.WriteAsJson(map[string]string{
		"successfully deleted product with id:": id,
	})

}
