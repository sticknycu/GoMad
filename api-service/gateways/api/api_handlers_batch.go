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
		} else {
			savedProducts = append(savedProducts, id)
		}
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
		if !notExist {
			errors = append(errors, "object not exist")
		} else {
			products = append(products, id)
		}
	}

	_ = resp.WriteAsJson(map[string][]domain.Product{
		"products:": products,
	})
}

func (api *API) updateProductMemoryBatch(req *restful.Request, resp *restful.Response) {
	var products []domain.ProductDiff
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

	for _, id := range products {
		existedProduct, exists, err := api.storage.Get(id.ID)
		if err != nil {
			log.Errorf("Failed to get product from storage, err=%v", err)
			_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to get product from store"))
			return
		}
		if !exists {
			log.Infof("Product %s not in store", id)
			_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("product not found"))
			return
		} else {
			newProduct := existedProduct
			newProduct.Stock = id.Diff.Stock
			newProduct.Price = id.Diff.Price
			newProduct.Tags = id.Diff.Tags

			_, _, err = api.storage.Save(newProduct)
		}
	}
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
