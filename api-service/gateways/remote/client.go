package remote

import (
	"bytes"
	"encoding/json"
	"exam-api/domain"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// Implement the following client to connect to the remote storage server
const (
	baseUrl         = "http://localhost:8081/"
	storeProductUrl = baseUrl + "/store/product"
)

type Client struct {
	client http.Client
}

func NewClient(client http.Client) *Client {
	return &Client{client: client}
}

func (c *Client) Save(product domain.Product) (string, bool, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(product)
	if err != nil {
		log.Errorf("An error has occured while encoding product, err=%v\n", err)
	}

	resp, err := http.NewRequest(http.MethodPost, storeProductUrl, &buf)

	if err != nil {
		log.Errorf("An error has occured while making a request to store service, err= %v\n", err)
		return "", false, err
	}

	body := resp.Body
	if body == nil {
		log.Errorf("[ERROR] Couldn't read request body, err=%v\n", err)
		return "", false, nil

	}
	defer body.Close()
	if err != nil {
		log.Errorf("[ERROR] Couldn't read request body, err=%v\n", err)
		return "", false, nil
	}
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Errorf("[ERROR] Couldn't read request body, err=%v\n", err)
		return "", false, nil
	}

	var receivedProduct domain.Product
	unmarshalErr := json.Unmarshal(data, &receivedProduct)
	if unmarshalErr != nil {
		log.Printf("Failed to unmarshal Product, err=%v\n", unmarshalErr)
		return "", false, nil
	}

	return receivedProduct.GetHash(), true, nil
}

func (c *Client) Get(id string) (domain.Product, bool, error) {
	resp, err := http.Get(storeProductUrl + "?id=" + id)
	if err != nil {
		log.Errorf("An error has occured while making a request to store service, err= %v\n", err)
		return domain.Product{}, false, err
	}

	body := resp.Request.Body
	if body == nil {
		log.Printf("[ERROR] Couldn't read request body, err=%v\n", err)
		return domain.Product{}, false, nil

	}
	defer body.Close()
	if err != nil {
		log.Printf("[ERROR] Couldn't read request body, err=%v\n", err)
		return domain.Product{}, false, nil
	}
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("[ERROR] Couldn't read request body, err=%v\n", err)
		return domain.Product{}, false, nil
	}

	var product domain.Product
	unmarshalErr := json.Unmarshal(data, &product)
	if unmarshalErr != nil {
		log.Printf("Failed to unmarshal Product, err=%v\n", unmarshalErr)
		return domain.Product{}, false, nil
	}

	return product, true, nil
}

func (c *Client) Update(id string, diff domain.Product) (bool, error) {
	var buf bytes.Buffer
	var newDiff = struct {
		Price int      `json:"price"`
		Stock int      `json:"stock"`
		Tags  []string `json:"tags"`
	}{
		Price: diff.Price,
		Stock: diff.Stock,
		Tags:  diff.Tags,
	}
	newProductDiff := domain.ProductDiff{
		ID:   id,
		Diff: newDiff,
	}
	err := json.NewEncoder(&buf).Encode(newProductDiff)
	if err != nil {
		log.Errorf("An error has occured while encoding product, err=%v\n", err)
	}

	resp, err := http.NewRequest(http.MethodPatch, storeProductUrl, &buf)

	if err != nil {
		log.Errorf("An error has occured while making a request to store service, err= %v\n", err)
		return false, err
	}

	body := resp.Body
	if body == nil {
		log.Errorf("[ERROR] Couldn't read request body, err=%v\n", err)
		return false, nil

	}
	defer body.Close()
	if err != nil {
		log.Errorf("[ERROR] Couldn't read request body, err=%v\n", err)
		return false, nil
	}
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Errorf("[ERROR] Couldn't read request body, err=%v\n", err)
		return false, nil
	}

	var receivedProduct domain.Product
	unmarshalErr := json.Unmarshal(data, &receivedProduct)
	if unmarshalErr != nil {
		log.Printf("Failed to unmarshal Product, err=%v\n", unmarshalErr)
		return false, nil
	}

	return true, nil

}

func (c *Client) Delete(id string) (bool, error) {
	resp, err := http.NewRequest(http.MethodDelete, storeProductUrl+"?id="+id, nil)

	if err != nil {
		log.Errorf("An error has occured while making a request to store service, err= %v\n", err)
		return false, err
	}

	body := resp.Body
	if body == nil {
		log.Errorf("[ERROR] Couldn't read request body, err=%v\n", err)
		return false, nil

	}
	defer body.Close()
	if err != nil {
		log.Errorf("[ERROR] Couldn't read request body, err=%v\n", err)
		return false, nil
	}
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Errorf("[ERROR] Couldn't read request body, err=%v\n", err)
		return false, nil
	}

	var receivedProduct domain.Product
	unmarshalErr := json.Unmarshal(data, &receivedProduct)
	if unmarshalErr != nil {
		log.Printf("Failed to unmarshal Product, err=%v\n", unmarshalErr)
		return false, nil
	}

	return true, nil
}
