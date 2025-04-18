package handler_functions

import (
	"encoding/json"
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

// Обработчик для страницы "Создать маршрут"
func CreateRouteHandler(p *pool_conections.Pool_conections) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tmpl, data, err := check_username(w, r, "web/pages/create_route.html")
		if err != nil {
			log.Printf("%s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if data.UserID == 0 {
			http.Redirect(w, r, "/?openLoginDialog=true", http.StatusSeeOther)
			return
		}

		if r.Method == http.MethodPost {
			// Проверяем Content-Type
			err := r.ParseMultipartForm(32 << 20) // 32MB максимум
			if err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}

			type CredentialsCreateRoute struct {
				Yandex_route      string
				Route_name        string
				Route_place       string
				Route_description string
				Route_preview     string
			}
			var creds CredentialsCreateRoute

			// Получаем текстовые данные из формы
			creds.Yandex_route = r.FormValue("routeLink")
			creds.Route_name = r.FormValue("route_name")
			creds.Route_place = r.FormValue("route_place")
			creds.Route_description = r.FormValue("route_description")
			creds.Route_preview = r.FormValue("route_description")

			// Проверяем обязательные поля
			if creds.Yandex_route == "" || creds.Route_name == "" || creds.Route_place == "" || creds.Route_description == "" || creds.Route_preview == "" {
				http.Error(w, "All fields are required", http.StatusBadRequest)
				return
			}

			files := r.MultipartForm.File["photos"]

			// Здесь можно обработать загруженные файлы, например:
			for _, fileHeader := range files {
				file, err := fileHeader.Open()
				if err != nil {
					log.Printf("Error opening file: %v", err)
					continue
				}
				defer file.Close()
				uploadDir := "./web/photos"
				if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
					err = os.Mkdir(uploadDir, 0755)
					if err != nil {
						log.Printf("Failed to create upload directory: %v", err)
						http.Error(w, "Internal server error", http.StatusInternalServerError)
						return
					}
				}
				var savedFiles []string // Для хранения путей сохраненных файлов

				for _, fileHeader := range files {
					// Ограничиваем размер файла (например, 5MB)
					if fileHeader.Size > 10<<20 {
						log.Printf("File %s is too large", fileHeader.Filename)
						continue
					}

					// Проверяем расширение файла
					ext := filepath.Ext(fileHeader.Filename)
					allowedExtensions := map[string]bool{
						".jpg":  true,
						".jpeg": true,
						".png":  true,
						".gif":  true,
					}
					if !allowedExtensions[strings.ToLower(ext)] {
						log.Printf("Invalid file extension: %s", ext)
						continue
					}

					// Генерируем уникальное имя файла
					uniqueID := uuid.New().String()
					newFilename := uniqueID + ext
					filePath := filepath.Join(uploadDir, newFilename)

					// Открываем и сохраняем файл
					file, err := fileHeader.Open()
					if err != nil {
						log.Printf("Error opening file: %v", err)
						continue
					}
					defer file.Close()

					// Создаем новый файл на диске
					dst, err := os.Create(filePath)
					if err != nil {
						log.Printf("Error creating file: %v", err)
						continue
					}
					defer dst.Close()

					// Копируем содержимое файла
					if _, err := io.Copy(dst, file); err != nil {
						log.Printf("Error saving file: %v", err)
						continue
					}

					savedFiles = append(savedFiles, filePath)
					log.Printf("Successfully saved file: %s", filePath)
				}
			}

			// Создаем маршрут
			route := queries.ConstructRoute(creds.Yandex_route, creds.Route_name, creds.Route_place, creds.Route_description, data.UserID)
			err = route.CreateNewRoute(p.PoolConns)
			if err != nil {
				log.Printf("Error query CreateNewRoute: %s", err.Error())
				http.Error(w, "Failed to create route", http.StatusInternalServerError)
				return
			}

			// Возвращаем успешный ответ
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Route created successfully"})
			return
		}

		// Обработка GET-запроса
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			log.Printf("Failed to render template: %s", err.Error())
		}
	}
	return http.HandlerFunc(fn)
}
