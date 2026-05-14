package handler_functions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"travel/internal/repositories/pool_conections"
	"travel/internal/repositories/queries"
)

// AddReviewResponse — структура ответа для фронтенда
type AddReviewResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func AddReviewHandler(p *pool_conections.Pool_conections) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// 🔐 Только POST
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(AddReviewResponse{
				Success: false,
				Error:   "Метод не разрешён",
			})
			return
		}

		// 🔐 Проверка авторизации через ВАШУ функцию check_username_data
		userData, err := check_username_data(r)
		if err != nil {
			log.Printf("Auth error: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(AddReviewResponse{
				Success: false,
				Error:   "Ошибка проверки сессии",
			})
			return
		}

		// Проверяем, что пользователь авторизован (UserID == 0 — не авторизован)
		if userData.UserID == 0 || userData.Username == "" {
			log.Printf("Unauthorized: userData=%+v", userData)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(AddReviewResponse{
				Success: false,
				Error:   "Пользователь не авторизован",
			})
			return
		}

		log.Printf("✅ Authorized: userID=%d, username=%s", userData.UserID, userData.Username)

		// 📥 Парсинг FormData
		err = r.ParseMultipartForm(10 << 20) // 10 MB лимит
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(AddReviewResponse{
				Success: false,
				Error:   "Ошибка парсинга формы",
			})
			return
		}

		// 📥 Извлечение и валидация данных
		routeIDStr := r.FormValue("route_id")
		description := r.FormValue("description")
		estimationStr := r.FormValue("estimation")

		routeID, err := strconv.Atoi(routeIDStr)
		if err != nil || routeID <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(AddReviewResponse{
				Success: false,
				Error:   "Неверный ID маршрута",
			})
			return
		}

		if len(description) < 10 || len(description) > 1000 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(AddReviewResponse{
				Success: false,
				Error:   "Текст отзыва: от 10 до 1000 символов",
			})
			return
		}

		estimation, err := strconv.ParseFloat(estimationStr, 64)
		if err != nil || estimation < 1 || estimation > 5 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(AddReviewResponse{
				Success: false,
				Error:   "Оценка должна быть от 1 до 5",
			})
			return
		}

		// 🚫 Проверка: не оставлял ли пользователь уже отзыв на этот маршрут
		exists, err := queries.CheckReviewExists(routeID, userData.UserID, p.PoolConns)
		if err != nil {
			log.Printf("Error checking review existence: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(AddReviewResponse{
				Success: false,
				Error:   "Ошибка проверки отзыва",
			})
			return
		}
		if exists {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(AddReviewResponse{
				Success: false,
				Error:   "Вы уже оставляли отзыв на этот маршрут",
			})
			return
		}

		// 📝 Создание и сохранение отзыва
		// Date оставим пустым — он установится в БД как NOW() в CreateReview
		review := queries.ConstructReview(
			userData.Username, // из контекста через check_username_data
			description,       // из формы
			estimation,        // из формы (float64)
			"",                // date: пустая строка → в БД будет время создания
		)

		// Используем ваш метод CreateReview
		if err := review.CreateReview(routeID, userData.UserID, p.PoolConns); err != nil {
			log.Printf("Error creating review: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(AddReviewResponse{
				Success: false,
				Error:   "Не удалось сохранить отзыв",
			})
			return
		}

		// 🔄 Обновляем средний рейтинг маршрута
		if err := queries.UpdateRouteEstimation(routeID, p.PoolConns); err != nil {
			log.Printf("Warning: failed to update route estimation: %v", err)
			// Не прерываем ответ — отзыв уже сохранён
		}

		// ✅ Успех
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(AddReviewResponse{
			Success: true,
			Message: "Отзыв успешно добавлен",
		})
	})
}
