package handler

import (
	"hash/crc32"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/speps/go-hashids"
)

func CoderUrlHandler(ctx *gin.Context) {
	decodeUrl := ctx.Query("decodeUrl")

	if decodeUrl == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "decodeUrl parameter is required",
		})
		return
	}

	decodeUrl, err := url.QueryUnescape(decodeUrl)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid URL encoding",
		})
		return
	}

	hashed := urlEncode(decodeUrl)

	ctx.JSON(http.StatusOK, gin.H{
		"hashedUrl": hashed,
	})
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
	ctx.JSON(http.StatusOK, gin.H{
		"decodedUrl": "url",
	})
}
