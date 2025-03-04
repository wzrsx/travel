ymaps.ready(init);

function init() {
    const map = new ymaps.Map("map", {
        center: [55.76, 37.64], // Центр карты (Москва)
        zoom: 10
    });

    let startPoint = null;
    let endPoint = null;
    const waypoints = [];  // Контрольные точки

    // Обработчик клика по карте
    map.events.add('click', function (e) {
            const coords = e.get('coords');

            if (!startPoint) {
                // Установка начальной точки
                startPoint = new ymaps.Placemark(coords, {
                    hintContent: 'Начальная точка'
                });
                map.geoObjects.add(startPoint);
            } else if (!endPoint) {
                // Установка конечной точки
                endPoint = new ymaps.Placemark(coords, {
                    hintContent: 'Конечная точка'
                });
                map.geoObjects.add(endPoint);
            } else {
                // Добавление контрольной точки
                const waypoint = new ymaps.Placemark(coords, {
                    hintContent: 'Контрольная точка'
                });
                map.geoObjects.add(waypoint);
                waypoints.push(coords);
            }
        });

        //обработчик кнопки создать маршрут
        document.getElementById('create-route-form').addEventListener('submit', function (e) {
            e.preventDefault();
            const route_name = getElementById('route-name').Value
            const route_place = getElementById('route-name').Value
            const route_description = getElementById('route-description').Value

            if(route_name.length <= 1 || route_place.length <= 1){
                alert("Ошибка: поле имени или места маршрута должно быть не меньше двух символов.");
                return;
            }

            if (!startPoint || !endPoint) {
                alert('Пожалуйста, укажите начальную и конечную точки на карте.');
                return;
            }

            // Получаем координаты точек
            const startCoords = startPoint.geometry.getCoordinates().join(',');
            const endCoords = endPoint.geometry.getCoordinates().join(',');
            const waypointsCoords = waypoints.map(coord => coord.join(',')).join('|');

            // Генерация ссылки на маршрут
            const routeLink = generateRouteLink(startCoords, endCoords, waypointsCoords);
            alert('Ссылка на маршрут: ' + routeLink);

            // Сохранение ссылки в базу данных (пример)
            saveRoute(routeLink, route_name, route_place, route_description);
        });
}

function generateRouteLink(startCoords, endCoords, waypointsCoords) {
    const baseUrl = "https://yandex.ru/maps/?mode=routes";
    const params = new URLSearchParams({
        rtext: `${startCoords}~${endCoords}`,
        rtt: 'auto' // Тип маршрута (auto, mt, pd)
    });

    if (waypointsCoords) {
        params.set('rtext', `${startCoords}~${waypointsCoords.split('|').join('~')}~${endCoords}`);
    }

    return `${baseUrl}&${params.toString()}`;
}

function saveRoute(routeLink, route_name, route_place, route_description) {
    fetch('/save-route', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            routeLink: routeLink,
            route_name: route_name,
            route_place: route_place,
            route_description: route_description
        })
    })
    .then(response => response.json())
    .then(data => {
        console.log('Ссылка сохранена:', data);
    })
    .catch(error => {
        console.error('Ошибка:', error);
    });
}