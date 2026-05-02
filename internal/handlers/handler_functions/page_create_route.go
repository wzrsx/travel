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
			log.Printf("Template error: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Проверка авторизации
		if data.UserID == 0 {
			http.Redirect(w, r, "/?openLoginDialog=true", http.StatusSeeOther)
			return
		}

		// Обработка POST-запроса (создание маршрута)
		if r.Method == http.MethodPost {
			// Парсим multipart-форму (лимит 32 МБ)
			err := r.ParseMultipartForm(32 << 20)
			if err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}

			// Структура для данных формы
			type CredentialsCreateRoute struct {
				Yandex_route      string
				Route_name        string
				Route_place       string
				Route_description string
			}
			var creds CredentialsCreateRoute

			// Получаем текстовые поля из формы
			creds.Yandex_route = r.FormValue("routeLink")
			creds.Route_name = r.FormValue("route_name")
			creds.Route_place = r.FormValue("route_place")
			creds.Route_description = r.FormValue("route_description")

			// Валидация обязательных полей
			if creds.Yandex_route == "" || creds.Route_name == "" ||
				creds.Route_place == "" || creds.Route_description == "" {
				http.Error(w, "All text fields are required", http.StatusBadRequest)
				return
			}

			// Получаем загруженные файлы
			files := r.MultipartForm.File["photos"]
			if len(files) == 0 {
				http.Error(w, "At least one photo is required", http.StatusBadRequest)
				return
			}

			// Подготовка директории для загрузки
			uploadDir := "./web/photos"
			if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
				if err = os.MkdirAll(uploadDir, 0755); err != nil {
					log.Printf("Failed to create upload directory: %v", err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
					return
				}
			}

			// Переменная для пути к превью (объявляем ДО цикла, чтобы использовать после)
			var previewWebPath string

			// Обрабатываем файлы (первый успешный будет использован как превью)
			for i, fileHeader := range files {
				// Проверка размера файла (макс. 10 МБ)
				if fileHeader.Size > 10<<20 {
					log.Printf("File %s is too large (%d bytes)", fileHeader.Filename, fileHeader.Size)
					continue
				}

				// Проверка расширения
				ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
				allowedExtensions := map[string]bool{
					".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
				}
				if !allowedExtensions[ext] {
					log.Printf("Invalid file extension: %s", ext)
					continue
				}

				// Генерация уникального имени файла
				uniqueID := uuid.New().String()
				newFilename := uniqueID + ext
				filePath := filepath.Join(uploadDir, newFilename)

				// Открываем исходный файл
				srcFile, err := fileHeader.Open()
				if err != nil {
					log.Printf("Error opening uploaded file: %v", err)
					continue
				}

				// Создаём файл на диске
				dstFile, err := os.Create(filePath)
				if err != nil {
					srcFile.Close()
					log.Printf("Error creating file on disk: %v", err)
					continue
				}

				// Копируем содержимое
				_, err = io.Copy(dstFile, srcFile)
				srcFile.Close()
				dstFile.Close()

				if err != nil {
					log.Printf("Error saving file content: %v", err)
					// Удаляем повреждённый файл
					_ = os.Remove(filePath)
					continue
				}

				log.Printf("Successfully saved file: %s", filePath)

				// ✅ Если это первое успешно сохранённое фото — используем как превью
				if previewWebPath == "" {
					previewWebPath = "/photos/" + newFilename
					log.Printf("Preview path set: %s", previewWebPath)
				}

				// Если нужно только одно превью — выходим после первого успешного
				// Если хотите сохранять все фото — уберите этот break
				if i == 0 {
					break
				}
			}

			// Проверка: удалось ли сохранить хотя бы одно фото для превью
			if previewWebPath == "" {
				http.Error(w, "Failed to save route preview image", http.StatusInternalServerError)
				return
			}

			// ✅ Создаём маршрут с путём к превью
			// Убедитесь, что ConstructRoute принимает 6 параметров, включая photoPreviewPath
			route := queries.ConstructRoute(
				creds.Yandex_route,
				creds.Route_name,
				creds.Route_place,
				creds.Route_description,
				previewWebPath, // ← путь к превью для БД
				data.UserID,
			)

			err = route.CreateNewRoute(p.PoolConns)
			if err != nil {
				log.Printf("Error executing CreateNewRoute: %s", err.Error())
				http.Error(w, "Failed to create route in database", http.StatusInternalServerError)
				return
			}

			log.Printf("Route '%s' created successfully with preview: %s", creds.Route_name, previewWebPath)

			// Возвращаем успешный JSON-ответ
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Route created successfully",
			})
			return
		}

		// Обработка GET-запроса (отображение формы)
		err = tmpl.Execute(w, data)
		if err != nil {
			log.Printf("Failed to render template: %s", err.Error())
			http.Error(w, "Failed to render page", http.StatusInternalServerError)
		}
	}
	return http.HandlerFunc(fn)
}
