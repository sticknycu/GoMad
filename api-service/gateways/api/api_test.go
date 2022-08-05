package api

import (
	"github.com/emicklei/go-restful/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var api *API

var _ = Describe("API Handler", func() {
	var (
		ws        *restful.WebService
		container *restful.Container
		recorder  *httptest.ResponseRecorder
	)

	BeforeEach(func() {

		ws = new(restful.WebService)
		container = restful.NewContainer()
		container.Add(ws)
		recorder = httptest.NewRecorder()
		api.RegisterRoutes(ws)
	})

	Context("When getting a product", func() {
		It("should fail if id is mentioned in query parameters", func() {
			request, _ := http.NewRequest("GET", "/store/memory/product/single", nil)

			container.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
		})
	})
})
