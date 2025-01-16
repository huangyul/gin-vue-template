package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/huangyul/gin-vue-template/internal/dto"
	"github.com/huangyul/gin-vue-template/internal/pkg/errno"
	"github.com/huangyul/gin-vue-template/internal/service"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	FileBase = "./static"
)

type FileHandler struct {
	svc *service.FileService
}

func NewFileHandler(svc *service.FileService) *FileHandler {
	return &FileHandler{
		svc: svc,
	}
}

func (f *FileHandler) RegisterRoutes(s *gin.Engine) {
	ug := s.Group("/file")
	ug.POST("/upload", f.Upload)
	ug.GET("/delete/:id", f.Delete)
	ug.POST("/list", f.List)
	ug.GET("/get-option", f.Option)
}

func (f *FileHandler) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		WriteErrno(ctx, errno.BadRequest.SetMessage("上传文件为空"))
		return
	}

	userId := ctx.MustGet("userId").(int64)
	fileName := filepath.Base(file.Filename)
	fileExt := filepath.Ext(fileName)
	newFileName := uuid.New().String() + fileExt
	link := filepath.Join(FileBase, newFileName)
	if er := ctx.SaveUploadedFile(file, link); er != nil {
		WriteErrno(ctx, errno.InternalServerError.SetMessage("保存文件失败："+er.Error()))
		return
	}

	err = f.svc.Save(ctx, file, link, userId)
	if err != nil {
		WriteErrno(ctx, errno.InternalServerError.SetMessage(err.Error()))
		return
	}

	WriteSuccessResponse(ctx, gin.H{
		"url": f.fullFilePath(ctx, newFileName),
	})
}

func (f *FileHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteErrno(ctx, errno.BadRequest.SetMessage("非法文件id"))
		return
	}
	uId := ctx.MustGet("userId").(int64)
	file, err := f.svc.Delete(ctx, id, uId)
	if err != nil {
		WriteErrno(ctx, errno.BadRequest.SetMessage(err.Error()))
		return
	}

	go func() {
		normalizeLink := filepath.ToSlash(file.Link)
		fileUrl := strings.TrimPrefix(normalizeLink, "static/")
		osUrl := filepath.Join(FileBase, fileUrl)
		er := os.Remove(filepath.Clean(osUrl))
		if er != nil {
			fmt.Printf("文件删除失败，文件名：%s, 文件link: %s, 失败原因：%s", file.FileName, file.Link, er.Error())
		}

	}()

	WriteSuccess(ctx)
}

func (f *FileHandler) List(ctx *gin.Context) {

	var r dto.FileListQueryParam
	if err := ctx.ShouldBind(&r); err != nil {
		WriteErrno(ctx, errno.BadRequest.SetMessage(err.Error()))
		return
	}
	res, total, err := f.svc.List(ctx, r)
	if err != nil {
		WriteErrno(ctx, errno.BadRequest.SetMessage(err.Error()))
		return
	}
	type resp struct {
		Id         int64  `json:"id"`
		FileName   string `json:"file_name"`
		Link       string `json:"link"`
		UploadUser string `json:"upload_user"`
		UploadTime string `json:"upload_time"`
	}
	var resps []resp
	for _, item := range res {
		resps = append(resps, resp{
			Id:         item.Id,
			FileName:   item.FileName,
			Link:       item.Link,
			UploadUser: item.Uploader,
			UploadTime: item.CreateAt.Format(time.DateTime),
		})
	}
	WriteSuccessResponse(ctx, dto.ListResp{
		Data:  resps,
		Total: total,
	})
}

func (f *FileHandler) Option(ctx *gin.Context) {
	uOptions, err := f.svc.GetOption(ctx)
	if err != nil {
		WriteErrno(ctx, errno.InternalServerError.SetMessage(err.Error()))
		return
	}
	WriteSuccessResponse(ctx, gin.H{
		"user": uOptions,
	})
}

func (f *FileHandler) fullFilePath(ctx *gin.Context, filename string) string {
	return "http://" + ctx.Request.Host + "/static/" + filename
}
