package server

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
	Clicks   int    `json:"clicks"`
}

func NewHandler(service *Service) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())

		// weak ass router
		if ctx.IsPost() && path == "/shorten" {
			shortenHandler(ctx, service)
		} else if ctx.IsGet() && len(path) > 3 && path[:3] == "/r/" { // "/r/{shortURL}"
			shortURL := path[3:]
			redirectHandler(ctx, service, shortURL)
		} else if ctx.IsGet() && len(path) > 7 && path[:7] == "/stats/" { // "/stats/{shortURL}"
			shortURL := path[7:]
			statsHandler(ctx, service, shortURL)
		} else {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			ctx.SetBody([]byte("Not found"))
		}
	}
}

func shortenHandler(ctx *fasthttp.RequestCtx, service *Service) {
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
}

func redirectHandler(ctx *fasthttp.RequestCtx, service *Service, shortURL string) {
	m, err := service.GetOriginalURL(shortURL)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetBody([]byte("Not found"))
	}

	_ = service.AddClick(shortURL)
	ctx.Redirect(m.OriginalURL, fasthttp.StatusMovedPermanently)
}

func statsHandler(ctx *fasthttp.RequestCtx, service *Service, shortURL string) {
	m, err := service.GetOriginalURL(shortURL)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetBody([]byte("Not found"))
	}

	resp := shortenResponse{
		ShortURL: fmt.Sprintf("http://%s/r/%s", ctx.Host(), m.ShortURL),
		Clicks:   m.Clicks,
	}

	respJSON, _ := json.Marshal(resp)
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(respJSON)
}
