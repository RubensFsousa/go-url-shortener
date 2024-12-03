package handler

import (
	"errors"
	"hash/crc32"
	"net/http"
	"os"
	"strconv"
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

type CoderUrlRequest struct {
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
// @Param           request body CoderUrlRequest true     "Request body"
// @Success         201         {object} SuccessResponse  "Shortened URL created successfully"
// @Failure         400         {object} ErrorResponse    "Invalid parameter or URL already shortened"
// @Failure         500         {object} ErrorResponse    "Internal error while processing the URL"
// @Router          /url/codeUrl [post]
func CoderUrlHandler(ctx *gin.Context) {
	var request CoderUrlRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		logger.Errorf("error to bind request %v", err)
		sendError(ctx, http.StatusBadRequest, "bind request error")
	}

	url := request.Url

	if strings.TrimSpace(url) == "" {
		sendError(ctx, http.StatusBadRequest, "url parameter is required")
		return
	}

	hashed := urlEncode(url, ctx)

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

func urlEncode(url string, ctx *gin.Context) string {
	hd := hashids.NewData()
	hd.Salt = os.Getenv("HASH_SALT")
	minHashSize, err := strconv.Atoi(os.Getenv("MIN_HASH_SIZE"))
	if err != nil {
		logger.Error("error to cast MIN_HASH_SIZE env to int")
		sendError(ctx, http.StatusInternalServerError, "error on hashing")
	}

	hd.MinLength = minHashSize
	h, _ := hashids.NewWithData(hd)

	urlHash := int(crc32.ChecksumIEEE([]byte(url)))

	hash, _ := h.Encode([]int{urlHash})
	return hash
}

// DecoderUrlHandler decodes a shortened URL back to its original form.
// @Summary         Decode URL
// @Description     Takes a shortened URL and returns the original URL.
// @Tags            URL
// @Accept          json
// @Produce         json
// @Param           hash         path      string  true  "Shortened URL to decode"
// @Success         200         {object} SuccessResponse  "Original URL"
// @Failure         400         {object} ErrorResponse    "Error URL not found"
// @Failure         500         {object} ErrorResponse    "Internal error while searching the URL"
// @Router          /url/decodeUrl/{hash} [get]
func DecoderUrlHandler(ctx *gin.Context) {
	hash := ctx.Param("hash")

	if strings.TrimSpace(hash) == "" {
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
