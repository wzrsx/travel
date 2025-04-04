document.addEventListener('DOMContentLoaded', function() {
    // Элементы управления
    const showAllCheckbox = document.getElementById('showAllRouts');
    const routeItems = document.querySelectorAll('.route-item');


    // Обработчик для кнопок "Перейти к маршруту"
    const routeButtons = document.querySelectorAll('.route-btn');
    routeButtons.forEach(button => {
        button.addEventListener('click', function() {
            const routeId = this.getAttribute('data-route-id');
            const routeName = encodeURIComponent(this.getAttribute('data-route-name'));
            const routePlace = encodeURIComponent(this.getAttribute('data-route-place'));
            const routeYandex = encodeURIComponent(this.getAttribute('data-route-yandex'));
            const routeDescription = encodeURIComponent(this.getAttribute('data-route-description'));
            const routeEstimation = encodeURIComponent(this.getAttribute('data-route-estimation'));

            const url = `/route?route_id=${routeId}&route_name=${routeName}&route_place=${routePlace}&routeLink=${routeYandex}&route_description=${routeDescription}&route_estimation=${routeEstimation}`;
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

    // Функция фильтрации маршрутов с плавной анимацией
    function filterRoutes() {
        const showAll = showAllCheckbox.checked;
        
        routeItems.forEach(item => {
            const filledStars = item.querySelectorAll('.rating .filled').length;
            const shouldShow = showAll || filledStars > 3;
            
            if (shouldShow) {
                item.style.display = 'flex';
            } else {
                item.style.display = 'none';
            }
        });
    }


    // Обработчик изменения чекбокса
    showAllCheckbox.addEventListener('change', filterRoutes);
    filterRoutes();
    
});
function openSignInDialog(){
    location.href = "http://localhost:5050/?openLoginDialog=true";
}