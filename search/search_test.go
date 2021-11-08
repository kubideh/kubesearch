package search

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	randomdata "github.com/Pallinder/go-randomdata"
)

var _ = Describe("Search", func() {

	Describe("Query the server for an object that doesn't exist", func() {
		Context("Just searching for a non-existent Pod", func() {
			It("should return an empty response", func() {
				server := httptest.NewServer(http.DefaultServeMux)
				defer server.Close()

				params := url.Values{}
				params.Add("query", randomdata.SillyName())

				response, err := http.Get(fmt.Sprintf("%s/v1/search?query=%s", server.URL, params.Encode()))
				Expect(err).ShouldNot(HaveOccurred())

				body, _ := io.ReadAll(response.Body)
				response.Body.Close()

				Expect(response.StatusCode).To(Equal(http.StatusOK))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=utf-8"))
				Expect(string(body)).To(Equal(`{}`))
			})
		})
	})

	Describe("Asking the server to say Hello", func() {
		Context("Just saying Hello", func() {
			It("should return an empty response", func() {
				server := httptest.NewServer(http.DefaultServeMux)
				defer server.Close()

				response, err := http.Get(fmt.Sprintf("%s/v1/search", server.URL))
				Expect(err).ShouldNot(HaveOccurred())

				body, _ := io.ReadAll(response.Body)
				response.Body.Close()

				Expect(response.StatusCode).To(Equal(http.StatusOK))
				Expect(response.Header.Get("Content-Type")).To(Equal("application/json; charset=utf-8"))
				Expect(string(body)).To(Equal(`{}`))
			})
		})
	})

	Describe("Asking the server to do something else", func() {
		Context("Testing that the server won't accept other URI paths", func() {
			It("should return 404", func() {
				server := httptest.NewServer(http.DefaultServeMux)
				defer server.Close()

				response, err := http.Get(fmt.Sprintf("%s/%s", server.URL, randomdata.SillyName()))
				Expect(err).ShouldNot(HaveOccurred())

				body, _ := io.ReadAll(response.Body)
				response.Body.Close()

				Expect(response.StatusCode).To(Equal(http.StatusNotFound))
				Expect(response.Header.Get("Content-Type")).To(Equal("text/plain; charset=utf-8"))
				Expect(string(body)).To(Equal("404 page not found\n"))
			})
		})
	})
})
