document.addEventListener('DOMContentLoaded', function() {
    const routeButtons = document.querySelectorAll('.route-btn');

    routeButtons.forEach(button => {
        button.addEventListener('click', function() {
            const routeId = this.getAttribute('data-route-id');
            const route_name = this.getAttribute('data-route-name');
            const route_place = this.getAttribute('data-route-place');
            const routeLink = this.getAttribute('data-route-yandex');
            const route_description = this.getAttribute('data-route-description');
            const route_estimation = this.getAttribute('data-route-estimation');

            fetch(`/route?route_id=${routeId}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    route_name: route_name,
                    route_place: route_place,
                    routeLink: routeLink,
                    route_description: route_description,
                    route_estimation: route_estimation
                })
            })
            .then(response => {
                if (response.redirected) {
                    // Перенаправление на новую страницу
                    window.location.href = response.url;
                } else if (!response.ok) {
                    throw new Error('Ошибка при перенаправлении: ' + response.status);
                }
            })
            .catch(error => {
                console.error('Ошибка при отправке запроса:', error);
            });
        });
    });
});