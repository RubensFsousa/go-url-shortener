package handler

import (
	"errors"
	"hash/crc32"
	"net/http"
	"strings"

	"github.com/RubensFsousa/go-url-shortener/config"
	"github.com/RubensFsousa/go-url-shortener/models"
	"github.com/gin-gonic/gin"
	"github.com/speps/go-hashids"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	logger *config.Logger
)

type CodedUrlRequest struct {
	Url string `json:"url" binding:"required"`
}
type DecoderUrlRequest struct {
	Hash string `json:"hash" binding:"required"`
}

func InitializeHandler() {
	logger = config.GetLogger("handler")
	db = config.GetPSQL()
}

// CoderUrlHandler shortens a given URL.
// @Summary         Shorten URL
// @Description     Accepts an original URL as a query parameter and returns a shortened URL.
// @Tags            URL
// @Accept          json
// @Produce         json
// @Param           decodedUrl  query    string  true  "Original URL to be shortened"
// @Success         201         {string} string  "Shortened URL created successfully"
// @Failure         400         {string} string  "Invalid parameter or URL already shortened"
// @Failure         500         {string} string  "Internal error while processing the URL"
// @Router          /url/codeUrl [post]
func CoderUrlHandler(ctx *gin.Context) {
	var request CodedUrlRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		logger.Errorf("error to bind request %v", err)
		sendError(ctx, http.StatusBadRequest, "bind request error")
	}

	url := request.Url

	if strings.TrimSpace(url) == "" {
		sendError(ctx, http.StatusBadRequest, "url parameter is required")
		return
	}

	hashed := urlEncode(url)

	var existingUrl models.Url
	if err := db.Where("decoded_url = ?", url).First(&existingUrl).Error; err == nil {
		logger.Warnf("URL already exists: %v", url)
		sendError(ctx, http.StatusBadRequest, "URL already shortened")
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Errorf("Database error during URL lookup: %v", err)
		sendError(ctx, http.StatusInternalServerError, "Database error")
		return
	}

	newUrl := models.Url{
		DecodedUrl: url,
		CodedUrl:   hashed,
	}

	if err := db.Create(&newUrl).Error; err != nil {
		logger.Errorf("error registering url: %v", err)
		sendError(ctx, http.StatusInternalServerError, "error registering url")
		return
	}

	sendSuccess(ctx, http.StatusCreated, hashed)
}

func urlEncode(url string) string {
	hd := hashids.NewData()
	hd.Salt = "7bcff2aa"
	hd.MinLength = 8
	h, _ := hashids.NewWithData(hd)

	urlHash := int(crc32.ChecksumIEEE([]byte(url)))

	hash, _ := h.Encode([]int{urlHash})
	return hash
}

// DecoderUrlHandler decodes a shortened URL back to its original form.
// @Summary         Decode URL
// @Description     Takes a shortened URL (path parameter) and returns the original URL.
// @Tags            URL
// @Accept          json
// @Produce         json
// @Param           codedUrl  path      string  true  "Shortened URL to decode"
// @Success         200       {string} string  "Original URL"
// @Failure         400       {string} string  "Error finding the URL or URL not found"
// @Router          /url/decodeUrl [get]
func DecoderUrlHandler(ctx *gin.Context) {
	var request DecoderUrlRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		logger.Errorf("error to bind request %v", err)
		sendError(ctx, http.StatusBadRequest, "bind request error")
	}

	hash := request.Hash

	if hash == "" {
		sendError(ctx, http.StatusBadRequest, "hash parameter is required")
		return
	}

	url := models.Url{}
	if err := db.Where("coded_url = ?", hash).First(&url).Error; err != nil {
		logger.Errorf("error to find url: %v", err)
		sendError(ctx, http.StatusBadRequest, "error url not find")
		return
	}

	sendSuccess(ctx, http.StatusOK, url.DecodedUrl)
}
