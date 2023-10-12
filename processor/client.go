package processor

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/openapi"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/infrastructure/log"
	"studio.sunist.work/sunist-c/ceobebot-qqchanel/processor/message"
)

const (
	sandbox = "https://sandbox.api.sgroup.qq.com"
	prod    = "https://api.sgroup.qq.com"
)

var (
	httpClient = &http.Client{}
)

type Client interface {
	ReplyMessage(ctx Context, reply message.Message)
	SendMessage(ctx Context, content message.Message)
}

type client struct {
	api openapi.OpenAPI
}

func (c *client) ReplyMessage(ctx Context, reply message.Message) {
	msg := reply.Build()
	msg.MsgID = ctx.GetPayload().Message.ID
	c.executeMessageSend(ctx, reply.Type(), msg)
}

func (c *client) SendMessage(ctx Context, content message.Message) {
	msg := content.Build()
	c.executeMessageSend(ctx, content.Type(), msg)
}

func (c *client) sendMessage(ctx Context, message *dto.MessageToCreate) {
	payload := ctx.GetPayload()
	result, err := c.api.PostMessage(ctx.GetContext(), payload.Message.ChannelID, message)
	c.replyFailedHooker(ctx, result, err)
}

func (c *client) sendImage(ctx Context, message *dto.MessageToCreate) {
	result, err := sendImage(ctx, message)
	c.replyFailedHooker(ctx, result, err)
}

func (c *client) executeMessageSend(ctx Context, messageType message.Type, msg *dto.MessageToCreate) {
	switch messageType {
	case message.TextMessageType:
		c.sendMessage(ctx, msg)
	case message.ImageMessageType:
		c.sendImage(ctx, msg)
	}
}

func (c *client) replyFailedHooker(ctx Context, result *dto.Message, err error) {
	payload := ctx.GetPayload()
	if err != nil {
		logger.Error(log.NewFieldsWithError(err))
		if systemConfig.TestMode {
			_, _ = c.api.PostMessage(ctx.GetContext(), payload.Message.ChannelID, message.NewTextMessage().Text(err.Error()).Build())
		} else {
			_, _ = c.api.PostMessage(ctx.GetContext(), payload.Message.ChannelID, message.NewTextMessage().Text("回复被夹掉了，换个话题吧").Build())
		}
	}

	logger.Info(log.NewFieldsWithMessage("message sent").With("content", result.Content))
}

func NewClient(api openapi.OpenAPI) Client {
	return &client{api: api}
}

// sendImage 傻逼腾讯没给 sdk 搞上传图片的接口，只能自己实现
func sendImage(ctx Context, msg *dto.MessageToCreate) (result *dto.Message, err error) {
	var image io.Reader
	var getImageErr error

	// 获取图片二进制，如果是链接则下载，否则读取本地文件
	if strings.HasPrefix(msg.Image, "https://") || strings.HasPrefix(msg.Image, "http://") {
		log.Debug(log.NewFieldsWithMessage("download image").With("url", msg.Image))
		image, getImageErr = downloadImage(msg.Image)
	} else {
		log.Debug(log.NewFieldsWithMessage("load image").With("path", msg.Image))
		image, getImageErr = loadFsImage(msg.Image)
	}
	if getImageErr != nil {
		log.Error(log.NewFieldsWithError(getImageErr))
		return &dto.Message{}, getImageErr
	}

	// 构建请求地址
	var url = ""
	if systemConfig.TestMode {
		url = fmt.Sprintf("%s/channels/%s/messages", sandbox, ctx.GetPayload().Message.ChannelID)
	} else {
		url = fmt.Sprintf("%s/channels/%s/messages", prod, ctx.GetPayload().Message.ChannelID)
	}

	// 构建请求体
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("msg_id", msg.MsgID)
	formFile, createFormFileErr := writer.CreateFormFile("file_image", generateImageName(msg.Image))
	if _, createFormFileErr = io.Copy(formFile, image); createFormFileErr != nil {
		log.Error(log.NewFieldsWithError(createFormFileErr))
		return &dto.Message{}, createFormFileErr
	}
	if closeErr := writer.Close(); closeErr != nil {
		log.Error(log.NewFieldsWithError(closeErr))
		return &dto.Message{}, closeErr
	}
	req, buildRequestErr := http.NewRequest(http.MethodPost, url, payload)
	if buildRequestErr != nil {
		log.Error(log.NewFieldsWithError(buildRequestErr))
		return &dto.Message{}, buildRequestErr
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bot %d.%s", systemConfig.AppID, systemConfig.AppToken))
	req.Header.Add("Content-Type", writer.FormDataContentType())

	// 执行请求
	res, execRequestErr := httpClient.Do(req)
	if execRequestErr != nil {
		log.Error(log.NewFieldsWithError(execRequestErr))
		return &dto.Message{}, execRequestErr
	}

	// 解析返回结果，因为已经发送完了，即使有错误也不能处理，打日志即可
	body, readBodyErr := io.ReadAll(res.Body)
	if readBodyErr != nil {
		log.Error(log.NewFieldsWithError(readBodyErr))
	}

	type ErrorResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	var errorResponse ErrorResponse
	if unmarshalErr := json.Unmarshal(body, &errorResponse); unmarshalErr != nil {
		log.Error(log.NewFieldsWithError(unmarshalErr))
	} else if errorResponse.Code != 0 && errorResponse.Message != "" {
		// 发送失败，实际有错误
		log.Error(log.NewFieldsWithMessage("send image failed").With("code", errorResponse.Code).With("message", errorResponse.Message))
		logrus.WithField("status", "send message failed").Error(string(body))
		return &dto.Message{}, fmt.Errorf("send image failed: %s", errorResponse.Message)
	}

	log.Info(log.NewFieldsWithMessage("send image success").With("result", string(body)))
	logrus.WithField("status", "send message success").Info(string(body))
	_ = res.Body.Close()
	if unmarshalErr := json.Unmarshal(body, &result); unmarshalErr != nil {
		log.Error(log.NewFieldsWithError(unmarshalErr))
	}
	return result, nil
}

func downloadImage(url string) (image io.Reader, err error) {
	resp, downloadErr := http.Get(url)
	if downloadErr != nil {
		return nil, downloadErr
	}

	bytesOfImg, readImgErr := io.ReadAll(resp.Body)
	if readImgErr != nil {
		return nil, readImgErr
	}

	if systemConfig.TestMode {
		f, of := os.Create("./data/chat/cache/" + generateImageName(url))
		if of != nil {
			panic(of)
		}

		_, _ = f.Write(bytesOfImg)
		_ = f.Close()
	}

	return bytes.NewBuffer(bytesOfImg), nil
}

func loadFsImage(path string) (image io.Reader, err error) {
	info, getInfoErr := os.Stat(path)
	if getInfoErr != nil {
		return nil, getInfoErr
	}

	if info.IsDir() {
		return nil, fmt.Errorf("path %s is a directory", path)
	}

	if file, openErr := os.Open(path); openErr != nil {
		return nil, openErr
	} else {
		return file, nil
	}
}

func generateImageName(image string) string {
	hash := md5.Sum([]byte(image))
	return fmt.Sprintf("%x.png", hash)
}
