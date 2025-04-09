document.addEventListener('DOMContentLoaded', function() {
    const routeButtons = document.querySelectorAll('.route-btn');

    routeButtons.forEach(button => {
        button.addEventListener('click', function() {
            // Получаем данные из атрибутов кнопки
            const routeId = this.getAttribute('data-route-id');
            const routeName = encodeURIComponent(this.getAttribute('data-route-name'));
            const routePlace = encodeURIComponent(this.getAttribute('data-route-place'));
            const routeYandex = encodeURIComponent(this.getAttribute('data-route-yandex'));
            const routeDescription = encodeURIComponent(this.getAttribute('data-route-description'));
            const routeEstimation = encodeURIComponent(this.getAttribute('data-route-estimation'));

            // Формируем URL с параметрами
            const url = `/route?route_id=${routeId}&route_name=${routeName}&route_place=${routePlace}&routeLink=${routeYandex}&route_description=${routeDescription}&route_estimation=${routeEstimation}`;

            // Перенаправляем пользователя
            window.location.href = url;
        });
    });
    
    const reviewsButtons = document.querySelectorAll('.reviews-btn');
    reviewsButtons.forEach(button => {
        button.addEventListener('click', function() {
            const routeId = this.getAttribute('data-route-id');

            const url = `/route/reviews?route_id=${routeId}`;
            window.location.href = url;
        });
    });
});

  