package statscontroller

import (
	"encoding/json"
	"net/http"
	"strconv"

	databases "github.com/alfaa19/go-restapi-crud/database"
	"github.com/alfaa19/go-restapi-crud/helpers"
	"github.com/alfaa19/go-restapi-crud/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var ResponseSuccess = helpers.ResponseSuccess
var ResponseError = helpers.ResponseError

func GetAll(w http.ResponseWriter, r *http.Request) {
	var stats []models.Stats

	if err := databases.DB.Find(&stats).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	ResponseSuccess(w, http.StatusOK, stats)

}

func GetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err)
	}

	var stats models.Stats
	if err := databases.DB.First(&stats, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ResponseError(w, http.StatusNotFound, err)
			return
		default:
			ResponseError(w, http.StatusInternalServerError, err)
			return
		}
	}

	ResponseSuccess(w, http.StatusOK, stats)
}

func Create(w http.ResponseWriter, r *http.Request) {

	var stats models.Stats

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&stats); err != nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	defer r.Body.Close()

	if err := databases.DB.Create(&stats).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	ResponseSuccess(w, http.StatusCreated, stats)

}

func Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err)
	}

	var stats models.Stats
	if err := databases.DB.First(&stats, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			ResponseError(w, http.StatusNotFound, err)
			return
		default:
			ResponseError(w, http.StatusInternalServerError, err)
			return
		}
	}

	var updateStats map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateStats); err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	if err := databases.DB.Model(&stats).Updates(updateStats).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	ResponseSuccess(w, http.StatusOK, stats)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
	}

	var stats models.Stats

	if err := databases.DB.First(&stats, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {

			ResponseError(w, http.StatusNotFound, err)
			return
		} else {
			ResponseError(w, http.StatusInternalServerError, err)
			return
		}
	}

	if err := databases.DB.Delete(&stats).Error; err != nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	ResponseSuccess(w, http.StatusOK, map[string]interface{}{})

}
