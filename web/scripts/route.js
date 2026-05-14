  // ==========================================
  // 1. ИНИЦИАЛИЗАЦИЯ СЛАЙДЕРОВ
  // ==========================================
  function initSliders() {
    // Слайдер фотографий маршрута (.itcss-1)
    const slider1 = document.querySelector(".itcss-1");
    if (slider1) {
      const items = slider1.querySelectorAll(".itcss__item");
      if (items.length > 1) {
        new ItcSimpleSlider(".itcss-1", {
          loop: true,
          autoplay: true,
          swipe: true,
        });
      } else if (items.length === 1) {
        const item = items[0];
        item.style.opacity = 1;
        item.style.transform = "translate3D(0px, 0px, 0.1px)";
      }
    }

    // Слайдер отзывов (.itcss-2)
    const slider2 = document.querySelector(".itcss-2");
    if (slider2) {
      new ItcSimpleSlider(".itcss-2", {
        loop: false,
        autoplay: true,
        swipe: true,
      });
    }
  }

  // ==========================================
  // 2. ЯНДЕКС.КАРТА И МАРШРУТ
  // ==========================================
  function initYandexMap() {
    // Проверяем, что объект ymaps загружен
    if (typeof ymaps === 'undefined') {
      console.warn("Yandex Maps API не загружен");
      return;
    }

    // Проверяем наличие данных о маршруте
    if (typeof yandexRoute === 'undefined' || !yandexRoute.points) {
      console.warn("Данные о маршруте не переданы");
      return;
    }

    ymaps.ready(function () {
      // Разбираем ссылку
      const urlParts = yandexRoute.points.split("?");
      if (urlParts.length < 2) {
        console.error("Неверный формат ссылки маршрута");
        return;
      }
      
      const urlParams = new URLSearchParams(urlParts[1]);
      const rtext = urlParams.get("rtext");
      const rtt = urlParams.get("rtt");

      if (!rtext) {
        console.error("Не удалось извлечь координаты из ссылки.");
        return;
      }

      // Парсим координаты точек
      const points = rtext.split("~").map((point) => {
        const [lat, lon] = point.split(",");
        return [parseFloat(lat), parseFloat(lon)];
      });

      if (points.length < 2) {
        console.error("Для построения маршрута нужно минимум две точки.");
        return;
      }

      // Создаем карту
      const map = new ymaps.Map("map", {
        center: points[0],
        zoom: 10,
        controls: ["zoomControl", "geolocationControl"]
      });

      // Создаем маршрут
      const multiRoute = new ymaps.multiRouter.MultiRoute(
        {
          referencePoints: points,
          params: {
            routingMode: "pedestrian",
          },
        },
        {
          boundsAutoApply: true,
        }
      );

      // Обработчик успешного построения маршрута
      multiRoute.model.events.add("requestsuccess", function () {
        console.log("Маршрут успешно построен.");
        const activeRoute = multiRoute.getActiveRoute();

        if (!activeRoute) {
          console.error("Активный маршрут не найден.");
          return;
        }

        const properties = activeRoute.properties.getAll();

        if (!properties.distance || !properties.duration) {
          console.error("Данные о расстоянии или времени отсутствуют.");
          return;
        }

        // Форматируем расстояние
        const distanceMeters = properties.distance.value;
        const distanceKm = (distanceMeters / 1000).toFixed(2);

        // Форматируем время
        const durationSeconds = properties.duration.value;
        const hours = Math.floor(durationSeconds / 3600);
        const minutes = Math.floor((durationSeconds % 3600) / 60);
        const timeString = `${hours.toString().padStart(2, "0")}:${minutes.toString().padStart(2, "0")}`;

        // Обновляем DOM
        const distanceElement = document.querySelector(".about.distance-route .about-route-text p:last-child");
        const timeElement = document.querySelector(".about.time-route .about-route-text p:last-child");

        if (distanceElement) distanceElement.textContent = `${distanceKm} км`;
        if (timeElement) timeElement.textContent = timeString;
      });

      // Обработчик ошибки
      multiRoute.model.events.add("requestfail", function (error) {
        console.error("Ошибка при построении маршрута:", error);
      });

      // Добавляем маршрут на карту
      map.geoObjects.add(multiRoute);
    });
  }

  // ==========================================
  // 3. МОДАЛЬНОЕ ОКНО ОТЗЫВА
  // ==========================================
  function initReviewDialog() {
    const dialog = document.getElementById("review-dialog");
    if (!dialog) return; // Если диалога нет на странице, выходим

    const openBtn = document.querySelector(".btn-route.estimate");
    const closeBtn = document.getElementById("closeReviewDialog");
    const cancelBtn = document.getElementById("cancelReview");
    const form = document.getElementById("review-form");
    const textarea = document.getElementById("review-text");
    const charCount = document.getElementById("charCount");
    const fileInput = document.getElementById("review-photos");
    const fileNames = document.getElementById("file-names");

    // Открытие диалога
    if (openBtn) {
      openBtn.addEventListener("click", (e) => {
        e.preventDefault();
        
        // Проверка авторизации
        const profileLink = document.getElementById("profilUser");
        if (profileLink && profileLink.textContent.trim() === "Войти") {
          alert("Пожалуйста, войдите в систему, чтобы оставить отзыв.");
          return;
        }
        dialog.showModal();
      });
    }

    // Закрытие диалога
    const closeDialog = () => {
      dialog.close();
      if (form) form.reset();
      if (charCount) charCount.textContent = "0";
      if (fileNames) fileNames.textContent = "Файлы не выбраны";
      const errorSpan = document.getElementById("rating-error");
      if (errorSpan) errorSpan.style.display = "none";
    };

    if (closeBtn) closeBtn.addEventListener("click", closeDialog);
    if (cancelBtn) cancelBtn.addEventListener("click", closeDialog);
    
    // Закрытие по клику вне окна
    dialog.addEventListener("click", (e) => {
      const rect = dialog.getBoundingClientRect();
      const isInDialog = (rect.top <= e.clientY && e.clientY <= rect.top + rect.height &&
                          rect.left <= e.clientX && e.clientX <= rect.left + rect.width);
      if (!isInDialog) {
        closeDialog();
      }
    });

    // Счетчик символов
    if (textarea && charCount) {
      textarea.addEventListener("input", () => {
        charCount.textContent = textarea.value.length;
      });
    }

    // Отображение имен файлов
    if (fileInput && fileNames) {
      fileInput.addEventListener("change", () => {
        if (fileInput.files.length > 0) {
          const names = Array.from(fileInput.files).map(f => f.name).join(", ");
          fileNames.textContent = names.length > 30 ? names.substring(0, 30) + "..." : names;
        } else {
          fileNames.textContent = "Файлы не выбраны";
        }
      });
    }

    // Отправка формы
    if (form) {
      form.addEventListener("submit", async (e) => {
        e.preventDefault();

        // Валидация рейтинга
        const rating = document.querySelector('input[name="rating"]:checked');
        const errorSpan = document.getElementById("rating-error");
        
        if (!rating) {
          if (errorSpan) errorSpan.style.display = "block";
          return;
        }
        if (errorSpan) errorSpan.style.display = "none";

        // Подготовка данных
        const formData = new FormData(form);
        formData.append("estimation", rating.value);

        const submitBtn = form.querySelector('button[type="submit"]');
        if (submitBtn) {
          const originalBtnText = submitBtn.textContent;
          submitBtn.disabled = true;
          submitBtn.textContent = "Отправка...";
        }

        try {
          const response = await fetch("/add-review", {
            method: "POST",
            body: formData
          });

            if (response.ok) {
              const result = await response.json().catch(() => ({}));
              
              // Показываем уведомление (опционально)
              if (result.message) {
                  alert(result.message);
              }
              
              closeDialog();
              
              // 🔁 Перезагружаем страницу с чистым route_id, чтобы OpenRoutePage подгрузил свежие данные из БД
              const routeId = new URLSearchParams(window.location.search).get('route_id');
              if (routeId) {
                  window.location.href = `/route?route_id=${routeId}`;
              } else {
                  window.location.reload();
              }
          } else {
            const result = await response.json().catch(() => ({}));
            alert("Ошибка: " + (result.error || "Не удалось отправить отзыв"));
          }
        } catch (err) {
          console.error("Ошибка сети:", err);
          alert("Произошла ошибка при соединении с сервером");
        } finally {
          if (submitBtn) {
            submitBtn.disabled = false;
            submitBtn.textContent = "Отправить отзыв";
          }
        }
      });
    }
  }

  // ==========================================
  // 4. ТОЧКА ВХОДА
  // ==========================================
  document.addEventListener("DOMContentLoaded", () => {
    initSliders();
    initReviewDialog();
    // Карту инициализируем, если есть контейнер и скрипт Яндекс.Карт уже загружен
    if (document.getElementById("map")) {
      initYandexMap();
    }
  });

  // Глобальная функция для открытия диалога (если нужно вызвать из HTML)
  window.openReviewDialog = function() {
    const dialog = document.getElementById("review-dialog");
    if (dialog) dialog.showModal();
  };
