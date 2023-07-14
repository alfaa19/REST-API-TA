package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	config "github.com/alfaa19/gin-restAPI-redis/config/database"
	"github.com/alfaa19/gin-restAPI-redis/helpers"
	"github.com/alfaa19/gin-restAPI-redis/model"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
)

func GetAll(c *gin.Context) {

	ctx := context.TODO()
	key := "stats"
	mycache := cache.New(&cache.Options{
		Redis:      config.RedisDB,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	var stats []model.Stats
	if err := mycache.Get(ctx, key, &stats); err != nil {
		if err := config.DB.Find(&stats).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				helpers.ResponseError(c, http.StatusNotFound, err)
			} else {
				helpers.ResponseError(c, http.StatusInternalServerError, err)
			}
			return
		}
		if err := mycache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   key,
			Value: stats,
			TTL:   3 * time.Minute,
		}); err != nil {
			helpers.ResponseError(c, http.StatusInternalServerError, err)
			return
		}
		helpers.ResponseSuccess(c, "from database", stats, http.StatusOK)
		return
	}
	helpers.ResponseSuccess(c, "from redis", stats, http.StatusOK)
}

func GetById(c *gin.Context) {

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err)
		return
	}
	ctx := context.TODO()
	key := fmt.Sprintf("stats:%d", idInt)
	mycache := cache.New(&cache.Options{
		Redis:      config.RedisDB,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
	var stats model.Stats
	if err := mycache.Get(ctx, key, &stats); err != nil {
		if err := config.DB.First(&stats, idInt).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				helpers.ResponseError(c, http.StatusNotFound, err)
			} else {
				helpers.ResponseError(c, http.StatusInternalServerError, err)
			}
			return
		}
		if err := mycache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   key,
			Value: stats,
			TTL:   3 * time.Minute,
		}); err != nil {
			helpers.ResponseError(c, http.StatusInternalServerError, err)
			return
		}
		helpers.ResponseSuccess(c, "from database", stats, http.StatusOK)
		return
	}
	helpers.ResponseSuccess(c, "from redis", stats, http.StatusOK)
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

	var existingStats []model.Stats
	ctx := context.TODO()
	key := "stats"
	mycache := cache.New(&cache.Options{
		Redis:      config.RedisDB,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	if err := mycache.Get(ctx, key, &existingStats); err != nil {
		if err == cache.ErrCacheMiss {
			existingStats = []model.Stats{}
		} else {
			helpers.ResponseError(c, http.StatusInternalServerError, err)
			return
		}
	}

	mergedStats := append(existingStats, stats)
	// Menyimpan data yang telah digabungkan kembali ke Redis
	if err := mycache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: mergedStats,
		TTL:   3 * time.Minute,
	}); err != nil {
		helpers.ResponseError(c, http.StatusInternalServerError, err)
		return
	}
	helpers.ResponseSuccess(c, "", stats, http.StatusOK)
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

	ctx := context.TODO()
	key := "stats"
	mycache := cache.New(&cache.Options{
		Redis:      config.RedisDB,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	var existingStats []model.Stats
	if err := mycache.Get(ctx, key, &existingStats); err != nil {
		existingStats = []model.Stats{} // Jika data belum ada, inisialisasikan dengan slice kosong
	}

	// Cari indeks data yang akan diperbarui
	var updatedIndex int = -1
	for i, s := range existingStats {
		if s.Id == stats.Id {
			updatedIndex = i
			break
		}
	}

	if updatedIndex != -1 {
		// Perbarui data di slice existingStats
		existingStats[updatedIndex] = stats

		// Simpan kembali data di Redis
		if err := mycache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   key,
			Value: existingStats,
			TTL:   3 * time.Minute,
		}); err != nil {
			helpers.ResponseError(c, http.StatusInternalServerError, err)
			return
		}
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

	ctx := context.TODO()
	key := "stats"
	mycache := cache.New(&cache.Options{
		Redis:      config.RedisDB,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	var existingStats []model.Stats
	if err := mycache.Get(ctx, key, &existingStats); err != nil {
		existingStats = []model.Stats{} // Jika data belum ada, inisialisasikan dengan slice kosong
	}

	// Cari indeks data yang akan dihapus
	var deletedIndex int = -1
	for i, s := range existingStats {
		if s.Id == stats.Id {
			deletedIndex = i
			break
		}
	}

	if deletedIndex != -1 {
		// Hapus data dari slice existingStats
		existingStats = append(existingStats[:deletedIndex], existingStats[deletedIndex+1:]...)

		// Simpan kembali data yang telah dihapus ke Redis
		if err := mycache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   key,
			Value: existingStats,
			TTL:   3 * time.Minute,
		}); err != nil {
			helpers.ResponseError(c, http.StatusInternalServerError, err)
			return
		}
	}

	helpers.ResponseSuccess(c, "", map[string]interface{}{}, http.StatusOK)
}
