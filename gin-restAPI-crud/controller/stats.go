package controller

import (
	"net/http"
	"strconv"

	config "github.com/alfaa19/gin-restAPI-crud/config/mysql"
	"github.com/alfaa19/gin-restAPI-crud/helpers"
	"github.com/alfaa19/gin-restAPI-crud/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAll(c *gin.Context) {
	var stats []model.Stats

	if err := config.DB.Find(&stats).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ResponseError(c, http.StatusNotFound, err)
		} else {
			helpers.ResponseError(c, http.StatusInternalServerError, err)
		}

		return
	}

	helpers.ResponseSuccess(c, "", stats, http.StatusOK)
}

func GetById(c *gin.Context) {

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	var stats model.Stats

	if err := config.DB.First(&stats, idInt).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ResponseError(c, http.StatusNotFound, err)
		} else {
			helpers.ResponseError(c, http.StatusInternalServerError, err)
		}
		return
	}
	helpers.ResponseSuccess(c, "", stats, http.StatusOK)
}

func Create(c *gin.Context) {
	var stats model.Stats

	if err := c.ShouldBindJSON(&stats); err != nil {
		helpers.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	if err := config.DB.Create(&stats).Error; err != nil {
		helpers.ResponseError(c, http.StatusUnprocessableEntity, err)
		return
	}

	helpers.ResponseSuccess(c, "", stats, http.StatusCreated)
}

func Update(c *gin.Context) {

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		helpers.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	var stats model.Stats
	if err := config.DB.First(&stats, idInt).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ResponseError(c, http.StatusNotFound, err)
		} else {
			helpers.ResponseError(c, http.StatusInternalServerError, err)
		}
		return
	}

	// Memperbarui field-field yang diberikan dalam permintaan
	if err := config.DB.Model(&stats).Updates(updateData).Error; err != nil {
		helpers.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	helpers.ResponseSuccess(c, "", stats, http.StatusOK)
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)

	if err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	var stats model.Stats
	if err := config.DB.First(&stats, idInt).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helpers.ResponseError(c, http.StatusNotFound, err)
		} else {
			helpers.ResponseError(c, http.StatusInternalServerError, err)
		}
		return
	}

	if err := config.DB.Delete(&stats).Error; err != nil {
		helpers.ResponseError(c, http.StatusInternalServerError, err)
	}

	helpers.ResponseSuccess(c, "", map[string]interface{}{}, http.StatusOK)
}
