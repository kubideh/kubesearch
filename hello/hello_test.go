package hello

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	randomdata "github.com/Pallinder/go-randomdata"
)

var _ = Describe("Hello", func() {

	Describe("Asking the server to say Hello", func() {
		Context("Just saying Hello", func() {
			It("should print an HTML-friendly Hello message", func() {
				server := httptest.NewServer(http.DefaultServeMux)
				defer server.Close()

				response, err := http.Get(fmt.Sprintf("%s/v1/hello", server.URL))
				Expect(err).ShouldNot(HaveOccurred())

				body, _ := io.ReadAll(response.Body)
				response.Body.Close()

				statusCode := response.StatusCode
				contentType := response.Header.Get("Content-Type")
				helloMessage := string(body)

				Expect(statusCode).To(Equal(http.StatusOK))
				Expect(contentType).To(Equal("text/html; charset=utf-8"))
				Expect(helloMessage).To(Equal("<html><body>Hello World!</body></html>"))
			})
		})
	})

	Describe("Asking the server to say something else", func() {
		Context("Testing that the server won't accept other URI paths", func() {
			It("should return 404", func() {
				server := httptest.NewServer(http.DefaultServeMux)
				defer server.Close()

				response, err := http.Get(fmt.Sprintf("%s/%s", server.URL, randomdata.SillyName()))
				Expect(err).ShouldNot(HaveOccurred())

				body, _ := io.ReadAll(response.Body)
				response.Body.Close()

				statusCode := response.StatusCode
				contentType := response.Header.Get("Content-Type")
				helloMessage := string(body)

				Expect(statusCode).To(Equal(http.StatusNotFound))
				Expect(contentType).To(Equal("text/plain; charset=utf-8"))
				Expect(helloMessage).To(Equal("404 page not found\n"))
			})
		})
	})
})
