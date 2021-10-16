package hello

import (
	"io"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hello", func() {

	Describe("Asking the server to say Hello", func() {
		Context("Just saying Hello", func() {
			It("should print an HTML-friendly Hello message", func() {
				request := httptest.NewRequest(http.MethodGet, "http://localhost/v1/search", nil)
				writer := httptest.NewRecorder()

				Handler(writer, request)

				response := writer.Result()
				body, _ := io.ReadAll(response.Body)

				statusCode := response.StatusCode
				contentType := response.Header.Get("Content-Type")
				helloMessage := string(body)

				Expect(statusCode).To(Equal(http.StatusOK))
				Expect(contentType).To(Equal("text/html; charset=utf-8"))
				Expect(helloMessage).To(Equal("<html><body>Hello World!</body></html>"))
			})
		})
	})
})
