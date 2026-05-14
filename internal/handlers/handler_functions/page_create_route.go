package handler_functions

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"travel/internal/repositories/pool_conections"
	"travel/internal/repositories/queries"

	"github.com/google/uuid"
)

func CreateRouteHandler(p *pool_conections.Pool_conections) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tmpl, data, err := check_username(w, r, "web/pages/create_route.html")
		if err != nil {
			log.Printf("Template error: %s", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if data.UserID == 0 {
			http.Redirect(w, r, "/?openLoginDialog=true", http.StatusSeeOther)
			return
		}

		if r.Method != http.MethodPost {
			err = tmpl.Execute(w, data)
			if err != nil {
				log.Printf("Render template: %s", err)
				http.Error(w, "Failed to render page", http.StatusInternalServerError)
			}
			return
		}

		// Парсим форму
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Текстовые поля
		creds := struct {
			Yandex_route      string
			Route_name        string
			Route_place       string
			Route_description string
		}{
			Yandex_route:      r.FormValue("routeLink"),
			Route_name:        r.FormValue("route_name"),
			Route_place:       r.FormValue("route_place"),
			Route_description: r.FormValue("route_description"),
		}

		if creds.Yandex_route == "" || creds.Route_name == "" ||
			creds.Route_place == "" || creds.Route_description == "" {
			http.Error(w, "All text fields are required", http.StatusBadRequest)
			return
		}

		// Файлы
		files := r.MultipartForm.File["photos"]
		if len(files) == 0 {
			http.Error(w, "At least one photo is required", http.StatusBadRequest)
			return
		}

		// Директория загрузки
		uploadDir := "./web/photos"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			if err = os.MkdirAll(uploadDir, 0755); err != nil {
				log.Printf("Mkdir %s: %v", uploadDir, err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		var previewWebPath string
		var savedPhotoPaths []string

		// Обрабатываем ВСЕ файлы
		for _, fileHeader := range files {
			if fileHeader.Size > 10<<20 {
				log.Printf("File too large: %s", fileHeader.Filename)
				continue
			}

			ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
			if !map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}[ext] {
				log.Printf("Invalid extension: %s", ext)
				continue
			}

			uniqueID := uuid.New().String()
			newFilename := uniqueID + ext
			filePath := filepath.Join(uploadDir, newFilename)

			src, err := fileHeader.Open()
			if err != nil {
				log.Printf("Open uploaded file: %v", err)
				continue
			}

			dst, err := os.Create(filePath)
			if err != nil {
				src.Close()
				log.Printf("Create file %s: %v", filePath, err)
				continue
			}

			_, err = io.Copy(dst, src)
			src.Close()
			dst.Close()
			if err != nil {
				log.Printf("Copy file content: %v", err)
				_ = os.Remove(filePath)
				continue
			}

			webPath := "/photos/" + newFilename
			log.Printf("Saved: %s → %s", fileHeader.Filename, webPath)

			if previewWebPath == "" {
				previewWebPath = webPath
			}
			savedPhotoPaths = append(savedPhotoPaths, webPath)
		}

		if previewWebPath == "" {
			http.Error(w, "Failed to save preview image", http.StatusInternalServerError)
			return
		}

		// Создаём маршрут и получаем ID
		route := queries.ConstructRoute(
			creds.Yandex_route,
			creds.Route_name,
			creds.Route_place,
			creds.Route_description,
			previewWebPath,
			data.UserID,
		)

		routeID, err := route.CreateNewRouteWithID(p.PoolConns)
		if err != nil {
			log.Printf("Create route: %v", err)
			http.Error(w, "Failed to create route", http.StatusInternalServerError)
			return
		}
		log.Printf("Route created: id=%d, name=%s", routeID, creds.Route_name)

		// Сохраняем ВСЕ фото в таблицу photos
		for _, photoPath := range savedPhotoPaths {
			if err := queries.InsertPhoto(p.PoolConns, routeID, photoPath); err != nil {
				log.Printf("Insert photo %s: %v", photoPath, err)
				// Не прерываем весь запрос, продолжаем с остальными
			}
		}
		log.Printf("Saved %d photo(s) to DB for route %d", len(savedPhotoPaths), routeID)

		// Ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message":  "Route created successfully",
			"route_id": fmt.Sprintf("%d", routeID),
		})
	}
	return http.HandlerFunc(fn)
}
