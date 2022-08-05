package api

import (
	"encoding/json"
	"exam-api/domain"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func (api *API) createProductSingle(req *restful.Request, resp *restful.Response) {
	body := req.Request.Body
	if body == nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, "nil body"))
		return

	}
	defer body.Close()
	var err error
	if err != nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()))
		return
	}
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	var product domain.Product
	unmarshalErr := json.Unmarshal(data, &product)
	if unmarshalErr != nil {
		log.Printf("Failed to unmarshal Product, err=%v\n", unmarshalErr)
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()))
		return
	}
	id, alreadyExists, err := api.storage.Save(product)
	if err != nil || alreadyExists {
		log.Errorf("Failed to save product in database, err=%v", err)
		_ = resp.WriteError(http.StatusConflict, fmt.Errorf("Failed to save product to database"))
		return
	}
	log.Printf("Product with id %d has been saved to database", id)

	respData := make(map[string]string)
	respData["message"] = "Status OK"
	errWriteResponse := resp.WriteAsJson(respData)
	if errWriteResponse != nil {
		log.Errorf("An error has occured while writing status code, err=%v", errWriteResponse)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("An error has occured while writing status code"))
		return
	}
	return
}

func (api *API) getProductSingle(req *restful.Request, resp *restful.Response) {
	id := req.QueryParameter("id")
	if id == "" {
		log.Infof("No id provided in request")
		_ = resp.WriteError(http.StatusBadRequest, fmt.Errorf("id must be provided"))
		return
	}
	product, exists, err := api.storage.Get(id)
	if err != nil {
		log.Errorf("Failed to get product from database, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to get product from database"))
		return
	}
	if !exists {
		log.Infof("Product %s not in database", id)
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}
	_ = resp.WriteAsJson(product)
}

func (api *API) updateProductSingle(req *restful.Request, resp *restful.Response) {
	id := req.QueryParameter("id")
	if id == "" {
		log.Infof("No id provided in request")
		_ = resp.WriteError(http.StatusBadRequest, fmt.Errorf("id must be provided"))
		return
	}

	productDiff := domain.Product{}
	err := req.ReadEntity(productDiff)
	if err != nil {
		log.Printf("[ERROR] Failed to read user, err=%v", err)
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	_, err = api.storage.Update(id, productDiff)
	if err != nil {
		log.Printf("[ERROR] User does not exist to database")
		resp.WriteError(http.StatusConflict, fmt.Errorf("user does not exists"))
		return
	}

	log.Printf("Product with id %d has been updated to database", id)

	respData := make(map[string]string)
	respData["message"] = "Status OK"
	errWriteResponse := resp.WriteAsJson(respData)
	if errWriteResponse != nil {
		log.Errorf("An error has occured while writing status code, err=%v", errWriteResponse)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("An error has occured while writing status code"))
		return
	}
	return
}

func (api *API) deleteProductSingle(req *restful.Request, resp *restful.Response) {
	id := req.QueryParameter("id")
	if id == "" {
		log.Infof("No id provided in request")
		_ = resp.WriteError(http.StatusBadRequest, fmt.Errorf("id must be provided"))
		return
	}

	_, err := api.storage.Delete(id)
	if err != nil {
		log.Printf("[ERROR] User does not exist to database")
		resp.WriteError(http.StatusConflict, fmt.Errorf("user does not exists"))
		return
	}

	log.Printf("Product with id %d has been deleted from database", id)

	respData := make(map[string]string)
	respData["message"] = "Status OK"
	errWriteResponse := resp.WriteAsJson(respData)
	if errWriteResponse != nil {
		log.Errorf("An error has occured while writing status code, err=%v", errWriteResponse)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("An error has occured while writing status code"))
		return
	}
	return
}
