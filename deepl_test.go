package deepl_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bounoable/deepl"
	mock_http "github.com/bounoable/deepl/http/mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("Client.Translate", func() {
	var (
		request           chan *http.Request
		server            *ghttp.Server
		mockDeeplResponse string
		mockDeeplHeader   int
	)

	BeforeEach(func() {
		request = make(chan *http.Request, 1)
		mockDeeplResponse = "{}"
		mockDeeplHeader = http.StatusOK

		server = ghttp.NewServer()
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/translate"),
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					err := r.ParseForm()
					Ω(err).ShouldNot(HaveOccurred())
					request <- r
					w.WriteHeader(mockDeeplHeader)
					w.Write([]byte(mockDeeplResponse))
				}),
			),
		)
	})

	AfterEach(func() {
		defer server.Close()
	})

	var (
		authKey string
		client  *deepl.Client

		sourceText       string
		targetLang       deepl.Language
		opts             []deepl.TranslateOption
		resultText       string
		resultSourceLang deepl.Language
		resultError      error
	)

	BeforeEach(func() {
		authKey = "an-auth-key"
		sourceText = "This is an example text."
		targetLang = deepl.German
		opts = nil
		resultText = ""
		resultSourceLang = ""
		resultError = nil
	})

	JustBeforeEach(func() {
		client = deepl.New(authKey, deepl.BaseURL(server.URL()))
		resultText, resultSourceLang, resultError = client.Translate(
			context.Background(),
			sourceText,
			targetLang,
			opts...,
		)
	})

	itAddsTheRequiredFields(&request, &targetLang)

	It("adds the source text", func(done Done) {
		req := <-request
		Ω(req.FormValue("text")).Should(Equal(sourceText))
		close(done)
	})

	itHandlesOptions(&request, &opts)
	itHandlesErrors(&request, &mockDeeplHeader, &resultError)

	When("deepl responds with no translations (is that even possible?)", func() {
		BeforeEach(func() {
			mockDeeplResponse = `{"translations": []}`
		})

		It("returns an empty text", func() {
			Ω(resultText).Should(BeEmpty())
		})

		It("returns an empty language", func() {
			Ω(resultSourceLang).Should(Equal(deepl.Language("")))
		})

		It("returns an error", func() {
			Ω(resultError).Should(HaveOccurred())
		})
	})

	When("deepl responds with a translation", func() {
		BeforeEach(func() {
			mockDeeplResponse = `{"translations": [
				{
					"detected_source_language": "EN",
					"text": "Dies ist ein Beispieltext."
				}
			]}`
		})

		It("returns the translated text", func(done Done) {
			<-request
			Ω(resultText).Should(Equal("Dies ist ein Beispieltext."))
			close(done)
		})

		It("returns the detected source language", func(done Done) {
			<-request
			Ω(resultSourceLang).Should(Equal(deepl.English))
			close(done)
		})
	})
})

var _ = Describe("Client.TranslateMany", func() {
	var (
		request           chan *http.Request
		server            *ghttp.Server
		mockDeeplResponse string
		mockDeeplHeader   int
	)

	BeforeEach(func() {
		request = make(chan *http.Request, 1)
		mockDeeplResponse = "{}"
		mockDeeplHeader = http.StatusOK

		server = ghttp.NewServer()
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/translate"),
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					err := r.ParseForm()
					Ω(err).ShouldNot(HaveOccurred())
					request <- r
					w.WriteHeader(mockDeeplHeader)
					w.Write([]byte(mockDeeplResponse))
				}),
			),
		)
	})

	AfterEach(func() {
		defer server.Close()
	})

	var (
		authKey string
		client  *deepl.Client

		sourceTexts        []string
		targetLang         deepl.Language
		opts               []deepl.TranslateOption
		resultTranslations []deepl.Translation
		resultError        error
	)

	BeforeEach(func() {
		authKey = "an-auth-key"
		sourceTexts = nil
		targetLang = deepl.German
		opts = nil
		resultTranslations = nil
		resultError = nil
	})

	JustBeforeEach(func() {
		client = deepl.New(authKey, deepl.BaseURL(server.URL()))
		resultTranslations, resultError = client.TranslateMany(
			context.Background(),
			sourceTexts,
			targetLang,
			opts...,
		)
	})

	itAddsTheRequiredFields(&request, &targetLang)

	It("adds the source texts", func(done Done) {
		req := <-request
		Ω(req.Form["text"]).Should(Equal(sourceTexts))
		close(done)
	})

	itHandlesOptions(&request, &opts)
	itHandlesErrors(&request, &mockDeeplHeader, &resultError)

	When("the user provides texts", func() {
		BeforeEach(func() {
			sourceTexts = []string{
				"This is an example.",
				"C'est un autre texte.",
			}

			mockDeeplResponse = `{"translations": [
				{"detected_source_language": "EN", "text": "Dies ist ein Beispiel."},
				{"detected_source_language": "FR", "text": "Dies ist ein anderer Text."}
			]}`
		})

		It("returns the translations", func(done Done) {
			<-request
			Ω(resultTranslations).Should(Equal([]deepl.Translation{
				{DetectedSourceLanguage: "EN", Text: "Dies ist ein Beispiel."},
				{DetectedSourceLanguage: "FR", Text: "Dies ist ein anderer Text."},
			}))
			close(done)
		})
	})
})

func itAddsTheRequiredFields(request *chan *http.Request, targetLang *deepl.Language) {
	It("adds the auth key", func(done Done) {
		req := <-*request
		Ω(req.FormValue("auth_key")).Should(Equal("an-auth-key"))
		close(done)
	})

	It("uses the correct content-type", func(done Done) {
		req := <-*request
		Ω(req.Header.Get("Content-Type")).Should(Equal("application/x-www-form-urlencoded"))
		close(done)
	})

	It("adds the target lang", func(done Done) {
		req := <-*request
		Ω(req.FormValue("target_lang")).Should(Equal(string(*targetLang)))
		close(done)
	})
}

func itHandlesOptions(request *chan *http.Request, opts *[]deepl.TranslateOption) {
	Context("with SourceLang() option", func() {
		BeforeEach(func() {
			*opts = append(*opts, deepl.SourceLang(deepl.English))
		})

		It("adds the source lang", func(done Done) {
			req := <-*request
			Ω(req.FormValue("source_lang")).Should(Equal(string(deepl.English)))
			close(done)
		})
	})

	Context("with SplitSentences() option", func() {
		BeforeEach(func() {
			*opts = append(*opts, deepl.SplitSentences(deepl.SplitNone))
		})

		It("adds the split_sentences option", func(done Done) {
			req := <-*request
			Ω(req.FormValue("split_sentences")).Should(Equal(deepl.SplitNone.Value()))
			close(done)
		})
	})

	Context("with PreserveFormatting() option", func() {
		BeforeEach(func() {
			*opts = append(*opts, deepl.PreserveFormatting(true))
		})

		It("adds the preserve_formatting option", func(done Done) {
			req := <-*request
			Ω(req.FormValue("preserve_formatting")).Should(Equal("1"))
			close(done)
		})
	})

	Context("with Formatlity() option", func() {
		BeforeEach(func() {
			*opts = append(*opts, deepl.Formality(deepl.LessFormal))
		})

		It("adds the formality option", func(done Done) {
			req := <-*request
			Ω(req.FormValue("formality")).Should(Equal(deepl.LessFormal.Value()))
			close(done)
		})
	})
}

func itHandlesErrors(request *chan *http.Request, mockDeeplHeader *int, resultError *error) {
	Context("errors", func() {
		codes := []int{
			http.StatusBadRequest,
			http.StatusForbidden,
			http.StatusNotFound,
			http.StatusRequestEntityTooLarge,
			http.StatusTooManyRequests,
			456, // quota exceeded. character limit reached
			http.StatusServiceUnavailable,
		}

		for _, code := range codes {
			code := code
			Describe(http.StatusText(code), func() {
				BeforeEach(func() {
					*mockDeeplHeader = code
				})

				It("returns an error with a code", func(done Done) {
					<-*request
					var deeplError deepl.Error
					Ω(errors.As(*resultError, &deeplError)).Should(BeTrue())
					Ω(deeplError.Code).Should(Equal(code))
					close(done)
				})
			})
		}
	})
}

func TestClient_Translate_withCustomHTTPClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpClient := mock_http.NewMockClient(ctrl)

	client := deepl.New("an-auth-key", deepl.HTTPClient(httpClient))

	clientCalled := make(chan struct{})

	httpClient.EXPECT().
		Do(gomock.Any()).
		DoAndReturn(func(*http.Request) (*http.Response, error) {
			close(clientCalled)
			return httptest.NewRecorder().Result(), nil
		})

	client.Translate(context.Background(), "This is an example text.", deepl.German)

	assert.Eventually(t, func() bool {
		_, open := <-clientCalled
		return !open
	}, time.Second, time.Millisecond*100)
}

func TestClient_HTTPClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	httpClient := mock_http.NewMockClient(ctrl)

	client := deepl.New("an-auth-key", deepl.HTTPClient(httpClient))
	assert.Same(t, httpClient, client.HTTPClient())
}
