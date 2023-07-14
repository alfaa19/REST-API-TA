package statscontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	databases "github.com/alfaa19/go-restapi-redis/database"
	"github.com/alfaa19/go-restapi-redis/helpers"
	"github.com/alfaa19/go-restapi-redis/models"
	"github.com/go-redis/cache/v8"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var (
	ResponseSuccess = helpers.ResponseSuccess
	ResponseError   = helpers.ResponseError
)

func GetAll(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()
	key := "stats"
	mycache := cache.New(&cache.Options{
		Redis:      databases.RedisDB,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
	var stats []models.Stats
	if err := mycache.Get(ctx, key, &stats); err != nil {
		if err := databases.DB.Find(&stats).Error; err != nil {
			ResponseError(w, http.StatusInternalServerError, err)
			return
		}
		if err := mycache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   key,
			Value: stats,
			TTL:   3 * time.Minute,
		}); err != nil {
			ResponseError(w, http.StatusInternalServerError, err)
			return
		}
		ResponseSuccess(w, http.StatusCreated, stats, "from database")
		return
	}

	ResponseSuccess(w, http.StatusOK, stats, "from redis")

}

func GetOne(w http.ResponseWriter, r *http.Request) {

	ctx := context.TODO()
	var stats models.Stats
	vars := mux.Vars(r)
	id := vars["id"]
	key := fmt.Sprintf("stats:%s", id)
	mycache := cache.New(&cache.Options{
		Redis:      databases.RedisDB,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
	if err := mycache.Get(ctx, key, &stats); err != nil {
		intID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			ResponseError(w, http.StatusBadRequest, err)
			return
		}

		if err := databases.DB.First(&stats, intID).Error; err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				ResponseError(w, http.StatusNotFound, err)
				return
			default:
				ResponseError(w, http.StatusInternalServerError, err)
				return
			}
		}

		if err := mycache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   key,
			Value: stats,
			TTL:   3 * time.Minute,
		}); err != nil {
			ResponseError(w, http.StatusInternalServerError, err)
			return
		}

		ResponseSuccess(w, http.StatusCreated, stats, "from database")
		return
	}

	ResponseSuccess(w, http.StatusOK, stats, "from redis")
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

	ctx := context.TODO()
	key := "stats"
	mycache := cache.New(&cache.Options{
		Redis:      databases.RedisDB,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	var existingStats []models.Stats
	if err := mycache.Get(ctx, key, &existingStats); err != nil {
		existingStats = []models.Stats{} // Jika data belum ada, inisialisasikan dengan slice kosong
	}

	// Menggabungkan data baru dengan data yang sudah ada sebelumnya
	mergedStats := append(existingStats, stats)
	// Menyimpan data yang telah digabungkan kembali ke Redis
	if err := mycache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: mergedStats,
		TTL:   3 * time.Minute,
	}); err != nil {
		ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	ResponseSuccess(w, http.StatusCreated, stats, "")

}

func Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err)
		return
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

	// Update data di Redis juga
	ctx := context.TODO()
	key := "stats"
	mycache := cache.New(&cache.Options{
		Redis:      databases.RedisDB,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	var existingStats []models.Stats
	if err := mycache.Get(ctx, key, &existingStats); err != nil {
		existingStats = []models.Stats{} // Jika data belum ada, inisialisasikan dengan slice kosong
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
			ResponseError(w, http.StatusInternalServerError, err)
			return
		}
	}

	ResponseSuccess(w, http.StatusOK, stats, "")
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

	// Hapus data dari Redis juga
	ctx := context.TODO()
	key := "stats"
	mycache := cache.New(&cache.Options{
		Redis:      databases.RedisDB,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	var existingStats []models.Stats
	if err := mycache.Get(ctx, key, &existingStats); err != nil {
		existingStats = []models.Stats{} // Jika data belum ada, inisialisasikan dengan slice kosong
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
			ResponseError(w, http.StatusInternalServerError, err)
			return
		}
	}

	ResponseSuccess(w, http.StatusOK, map[string]interface{}{}, "")
}
