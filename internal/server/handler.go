package server

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
)

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
	Clicks   int    `json:"clicks"`
}

var reqCount uint32 = 0

func NewHandler(service *Service) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())

		// weak ass router
		if ctx.IsPost() && path == "/shorten" {
			log.Printf("[%d] POST on /shorten\n", reqCount)
			shortenHandler(ctx, service, reqCount)
			reqCount += 1
		} else if ctx.IsGet() && len(path) > 3 && path[:3] == "/r/" { // "/r/{shortURL}"
			shortURL := path[3:]
			log.Printf("[%d] GET on /r/%s\n", reqCount, shortURL)
			redirectHandler(ctx, service, shortURL, reqCount)
			reqCount += 1
		} else if ctx.IsGet() && len(path) > 7 && path[:7] == "/stats/" { // "/stats/{shortURL}"
			shortURL := path[7:]
			log.Printf("[%d] GET on /stats/%s\n", reqCount, shortURL)
			statsHandler(ctx, service, shortURL, reqCount)
			reqCount += 1
		} else if ctx.IsGet() && path == "/random" {
			log.Printf("[%d] GET on /random", reqCount)
			shortURL, err := service.GetRandomURL()
			if err != nil {
				log.Printf("[%d] GET: %v (this most likely occurs because your database is empty)\n", reqCount, err)
			} else {
				statsHandler(ctx, service, shortURL, reqCount)
			}
			reqCount += 1
		} else {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetBody([]byte("Not found"))
		}
	}
}

func shortenHandler(ctx *fasthttp.RequestCtx, service *Service, reqID uint32) {
	var req shortenRequest
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil || req.URL == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte("Invalid request body"))
	}

	m, err := service.ShortenURL(req.URL)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte("Failed to shorten URL"))
	}

	resp := shortenResponse{
		ShortURL: fmt.Sprintf("http://%s/r/%s", ctx.Host(), m.ShortURL),
		Clicks:   m.Clicks,
	}

	respJSON, _ := json.Marshal(resp)
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(respJSON)
	log.Printf("[%d] POST: %s\n", reqID, respJSON)
}

func redirectHandler(ctx *fasthttp.RequestCtx, service *Service, shortURL string, reqID uint32) {
	m, err := service.GetOriginalURL(shortURL)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetBody([]byte("Not found"))
	}

	_ = service.AddClick(shortURL)
	ctx.Redirect(m.OriginalURL, fasthttp.StatusMovedPermanently)
	log.Printf("[%d] GET: Redirecting to %s\n", reqID, m.OriginalURL)
}

func statsHandler(ctx *fasthttp.RequestCtx, service *Service, shortURL string, reqID uint32) {
	m, err := service.GetOriginalURL(shortURL)
	if err != nil || m == nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetBody([]byte("Not found"))
		return
	}

	resp := shortenResponse{
		ShortURL: fmt.Sprintf("http://%s/r/%s", string(ctx.Host()), m.ShortURL),
		Clicks:   m.Clicks,
	}

	respJSON, _ := json.Marshal(resp)
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(respJSON)
	log.Printf("[%d] GET: %s\n", reqID, respJSON)
}
