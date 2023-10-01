package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io"
	_ "nats-streaming-consumer/docs"
	"nats-streaming-consumer/internal/entity"
	"net/http"
)

type ModelProvider interface {
	SelectAllModels() ([]entity.Model, error)
	SelectModelByUid(uid string) (entity.Model, error)
}

type ModelPublisher interface {
	PublishModel(model entity.Model) error
}

type Handler struct {
	provider  ModelProvider
	publisher ModelPublisher
}

func NewHandler(provider ModelProvider, publisher ModelPublisher) *Handler {
	return &Handler{
		provider:  provider,
		publisher: publisher,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	engine := gin.New()

	engine.GET("/", h.welcome)
	engine.GET("/models", h.modelList)
	engine.GET("/model", h.modelById)
	engine.POST("/send-model", h.sendModel)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return engine
}

func (h *Handler) welcome(context *gin.Context) {
	_, _ = io.WriteString(context.Writer, "Welcome")
}

// @Summary      	get models
// @Description  	method to select models
// @Accept       	json
// @Consume      	json
// @Success         200 {string} string "Ok"
// @Router       	/models [get]
func (h *Handler) modelList(context *gin.Context) {
	models, err := h.provider.SelectAllModels()
	if err != nil {
		logrus.Errorf("Failed to select models: %s", err.Error())
		context.JSON(http.StatusBadRequest, "Failed to select models")
		return
	}
	context.JSON(http.StatusOK, models)
}

// @Summary      	get model by uid
// @Description  	method to select model by uid
// @Accept       	json
// @Consume      	json
// @Param        	uid	query	string  true  "uid заказа"  Format(email)
// @Success         200 {string} string "Ok"
// @Router       	/model [get]
func (h *Handler) modelById(context *gin.Context) {
	uid, hasParameter := context.GetQuery("uid")
	if !hasParameter {
		logrus.Errorf("No query parameter: uid")
		context.JSON(http.StatusBadRequest, "No query parameter: uid")
		return
	}

	model, err := h.provider.SelectModelByUid(uid)
	if err != nil {
		logrus.Errorf("Failed to select models: %s", err.Error())
		context.JSON(http.StatusBadRequest, "Failed to select models")
		return
	}
	context.JSON(http.StatusOK, model)
}

// @Summary      	send model
// @Description  	method to send model
// @Accept       	json
// @Consume      	json
// @Param 			model body entity.Model true "Объект model"
// @Success         200 {string} string "Ok"
// @Router       	/send-model [post]
func (h *Handler) sendModel(context *gin.Context) {
	var model entity.Model

	if err := context.BindJSON(&model); err != nil {
		logrus.Errorf("Failed to bind Json: %s", err.Error())
		context.JSON(http.StatusBadRequest, "Invalid data")
		return
	}
	err := h.publisher.PublishModel(model)
	if err != nil {
		logrus.Errorf("Failed to publish model: %s", err.Error())
		context.JSON(http.StatusBadRequest, "Failed to publish model")
		return
	}
	context.JSON(http.StatusOK, "Ok")
}
