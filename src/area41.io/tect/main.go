package main
/*-----------------------------------------------------------------------------
 **
 ** - TECT -
 **
 ** Copyright 2018 by Krogoth and the contributing authors
 **
 ** This program is free software; you can redistribute it and/or modify it
 ** under the terms of the GNU Affero General Public License as published by the
 ** Free Software Foundation, either version 3 of the License, or (at your option)
 ** any later version.
 **
 ** This program is distributed in the hope that it will be useful, but WITHOUT
 ** ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 ** FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
 ** for more details.
 **
 ** You should have received a copy of the GNU Affero General Public License
 ** along with this program. If not, see <http://www.gnu.org/licenses/>.
 **
 ** -----------------------------------------------------------------------------*/

import (
	"github.com/gin-gonic/gin"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"area41.io/tect/lib"
	"area41.io/tect/db"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func main() {

	//
	lib.InitConfig()

	// only log database actions when env is set to "dev"
	env := lib.GetStringFromConfig("env")
	bIsDevMode := env == "dev"

	db.InitDatabase(bIsDevMode)
	defer db.CloseDB()

	if bIsDevMode {

		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Group using gin.BasicAuth() middleware
	//authorized := r.Group("/", lib.BasicAuth())

	r.GET("/fw/:file", DeliverFirmwareFile)
	r.GET("/spiffs/:file", DeliverSpiffsFile)

	serveProteus := lib.GetBoolFromConfig("serveProteus")
	if serveProteus {

		wwwwRoot := lib.GetStringFromConfig("folder.www") + "/"

		r.StaticFile("/index.html", wwwwRoot + "index.html")
		r.StaticFile("/default.html", wwwwRoot + "index.html")
		r.StaticFile("/mntsrt.css", wwwwRoot + "mntsrt.css")
		r.StaticFile("/mntsrt-b.woff2", wwwwRoot + "mntsrt-b.woff2")
		r.StaticFile("/mntsrt-r.woff2", wwwwRoot + "mntsrt-r.woff2")
		r.StaticFile("/siimple.css", wwwwRoot + "siimple.css")
		r.StaticFile("/siimple-colors.css", wwwwRoot + "siimple-colors.css")
		r.StaticFile("/ttbl.json", wwwwRoot + "ttbl.json")
		r.StaticFile("/zepto.min.js", wwwwRoot + "zepto.min.js")

		// throttle to have the client show the loading info
		r.GET("/api/config", GetConfig)
		r.POST("/api/config", SetConfig)
		r.GET("/api/hwinfo", getHWInfo)
	}

	//e.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/assets/", assetHandler)))
	//e.GET("/", s.get403)

	host := lib.GetStringFromConfig("www.host")
	port := lib.GetStringFromConfig("www.port")

	fmt.Printf("TECT is running on %v:%v\n", host, port)

	r.Run(host + ":" + port) // listen and serve

}

func GetConfig(c *gin.Context) {

	var json = []byte(`
{
    "nickname": "Uberh4x0r"
}
`)

	c.String(http.StatusOK, "application/json; charset=utf-8", json)
}

func SetConfig(c *gin.Context) {

	c.String(http.StatusOK, "application/json; charset=utf-8", "")
}

func getHWInfo(c *gin.Context) {

	var json = []byte(`
{
    "hwinfo": {
        "restart_reason": "hard reset",
        "free_heap": "12"
	}
}
`)

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Writer.Write(json)
}

func DeliverFirmwareFile(c *gin.Context) {

	folder := lib.GetStringFromConfig("folder.firmware")

	deliverFile(c, folder)
}

func DeliverSpiffsFile(c *gin.Context) {

	folder := lib.GetStringFromConfig("folder.spiffs")

	deliverFile(c, folder)
}

func logAccess(c *gin.Context) {

	/*
	md5Header := c.GetHeader("x-esp8266-sketch-md5")

	x-ESP8266-STA-MAC: 60:01:94:7A:0C:75
	x-ESP8266-AP-MAC: 62:01:94:7A:0C:75
	x-ESP8266-free-space: 2846720
	x-ESP8266-sketch-size: 295408
	x-ESP8266-sketch-md5: 77855cef563bb0f168ec0d4a891f9272
	x-ESP8266-chip-size: 4194304
	x-ESP8266-sdk-version: 2.1.0(deb1901)
	x-ESP8266-mode: sketch
	*/
}

func deliverFile(c *gin.Context, folder string) {

	files, err := ioutil.ReadDir(folder)
	if err != nil {

		c.JSON(http.StatusUnavailableForLegalReasons, gin.H{
			"error": err,
		})

		return
	}

	filesCount := len(files)
	if filesCount <= 0 {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "File not found",
		})

		return
	}

	path := filepath.Join(folder, files[filesCount -1].Name())

	hash, err := lib.HashFileMD5(path)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})

		return
	}

	fmt.Printf("hash: %v\n", hash)

	md5Header := c.GetHeader("x-esp8266-sketch-md5")
	if md5Header == hash {

		fmt.Printf("hash did not change\n")

		c.JSON(http.StatusNotModified, gin.H{
			"error": "Hash did not change",
		})

		return
	}

	fmt.Printf("hash did change, delivering new firmware\n")

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("x-MD5", hash)
	c.File(path)
}


