package handler

import (
	"hash/crc32"
	"net/http"
	"net/url"

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

func InitializeHandler() {
	logger = config.GetLogger("handler")
	db = config.GetPSQL()
}

func CoderUrlHandler(ctx *gin.Context) {
	decodedUrl := ctx.Query("decodedUrl")

	if decodedUrl == "" {
		sendError(ctx, http.StatusBadRequest, "decodedUrl parameter is required")
		return
	}

	decodedUrl, err := url.QueryUnescape(decodedUrl)
	if err != nil {
		sendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	hashed := urlEncode(decodedUrl)

	url := models.Url{
		DecodedUrl: decodedUrl,
		CodedUrl:   hashed,
	}

	if existsUrl := db.Where("decoded_url = ?", decodedUrl).First(&url); existsUrl != nil {
		logger.Errorf("exist url with %v", decodedUrl)
		sendError(ctx, http.StatusBadRequest, "url already shorted")
		return
	}

	if err := db.Create(&url).Error; err != nil {
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

func DecoderUrlHandler(ctx *gin.Context) {
	codedUrl := ctx.Param("codedUrl")

	if codedUrl == "" {
		sendError(ctx, http.StatusBadRequest, "codedUrl parameter is required")
		return
	}

	url := models.Url{}
	if err := db.Where("coded_url = ?", codedUrl).First(&url).Error; err != nil {
		logger.Errorf("error to find url: %v", err)
		sendError(ctx, http.StatusBadRequest, "error url not find")
		return
	}

	sendSuccess(ctx, http.StatusOK, url.DecodedUrl)
}
