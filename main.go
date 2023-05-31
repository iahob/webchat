package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"webchat/connection"
	"webchat/db"
	"webchat/user"
)

//go:embed static
var assets embed.FS

var (
	mysqlDSN string
)

func init() {
	flag.StringVar(&mysqlDSN, "dsn", "", "mysql dsn")
}

func main() {
	flag.Parse()
	err := db.Init(mysqlDSN)
	if err != nil {
		fmt.Printf("init mysql by %s error %s", mysqlDSN, err.Error())
		return
	}
	engine := gin.New()
	engine.Use(static.Serve("", EmbedFolder(assets, "static")))

	engine.POST("/upload", func(ctx *gin.Context) {
		file, err := ctx.FormFile("image")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Bad Request",
			})
			return
		}
		err = ctx.SaveUploadedFile(file, "static/assert/image/"+file.Filename)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Image uploaded successfully",
			"url":     "static/assert/image/" + file.Filename,
		})
	})
	engine.GET("/ws", connection.Upgrade)
	engine.POST("/auth", func(ctx *gin.Context) {
		data := &struct {
			Name string `json:"name"`
			Pwd  string `json:"pwd"`
		}{}
		err := json.NewDecoder(ctx.Request.Body).Decode(data)
		if err != nil {
			return
		}
		err = user.Login(data.Name, data.Pwd)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "login fail",
				"code":    200,
			})
			return
		}
		token := uuid.New().String()
		um := user.UserModel{
			Token: token,
			Name:  data.Name,
			Pwd:   data.Pwd,
		}
		um.Update()
		ctx.JSON(http.StatusOK, gin.H{
			"message": "login success",
			"code":    100,
			"token":   token,
		})
	})
	// 开启协程启动connection服务管理中心
	go connection.DefaultH.Run()

	// 启动http服务
	err = engine.Run(":80")
	if err != nil {
		return
	}
}

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	if err != nil {
		return false
	}
	return true
}

func EmbedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem {
	efs, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(efs),
	}
}
