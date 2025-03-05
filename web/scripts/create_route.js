ymaps.ready(init);

function init() {
    const map = new ymaps.Map('map', {
        center: [55.76, 37.64], // Центр карты (Москва)
        zoom: 10
    });

    let startPoint = null;
    let endPoint = null;
    let waypoints = [];
    let multiRoute = null;

    // Функция для создания маршрута
    function updateRoute() {
        if (multiRoute) {
            map.geoObjects.remove(multiRoute);
        }

        if (startPoint && endPoint) {
            const referencePoints = [
                startPoint.geometry.getCoordinates(),
                ...waypoints.map(w => w.geometry.getCoordinates()),
                endPoint.geometry.getCoordinates()
            ];

            // Создание мультимаршрута с параметрами для пешего передвижения
            multiRoute = new ymaps.multiRouter.MultiRoute({
                referencePoints: referencePoints,
                params: {
                    routingMode: 'pedestrian' // Указываем тип маршрута как пешеходный
                }
            }, {
                boundsAutoApply: true // Автоматически подгонять карту под маршрут
            });

            // Добавляем маршрут на карту
            map.geoObjects.add(multiRoute);

            // Обработка кликов на точках маршрута
            multiRoute.model.events.add('requestsuccess', function () {
                multiRoute.getWayPoints().each(function (point) {
                    point.events.add('click', function () {
                        const coords = point.geometry.getCoordinates();

                        // Проверка на удаление начальной точки
                        if (startPoint && startPoint.geometry.getCoordinates().join(',') === coords.join(',')) {
                            map.geoObjects.remove(startPoint);
                            if (waypoints.length > 0) {
                                startPoint = waypoints[0];
                                map.geoObjects.add(startPoint);
                                waypoints.shift();
                            } else {
                                startPoint = null;
                            }
                            updateRoute();
                        }

                        // Проверка на удаление конечной точки
                        if (endPoint && endPoint.geometry.getCoordinates().join(',') === coords.join(',')) {
                            map.geoObjects.remove(endPoint);
                            if (waypoints.length > 0) {
                                endPoint = waypoints[waypoints.length - 1];
                                map.geoObjects.add(endPoint);
                                waypoints.pop();
                            } else {
                                endPoint = null;
                            }
                            updateRoute();
                        }

                        // Проверка на удаление контрольной точки
                        const waypointToRemove = waypoints.find(w => w.geometry.getCoordinates().join(',') === coords.join(','));
                        if (waypointToRemove) {
                            map.geoObjects.remove(waypointToRemove);
                            waypoints = waypoints.filter(w => w !== waypointToRemove);
                            updateRoute();
                        }
                    });
                });
            });
        }
    }

    // Обработчик клика на карте
    map.events.add('click', function (e) {
        const coords = e.get('coords');

        if (!startPoint) {
            // Установка начальной точки
            startPoint = new ymaps.Placemark(coords, {
                hintContent: 'Начальная точка'
            }, {
                preset: 'islands#greenDotIcon', // Зеленая точка
                draggable: true // Разрешить перемещение
            });
            map.geoObjects.add(startPoint);

            // Удаление начальной точки при клике
            startPoint.events.add('click', function () {
                map.geoObjects.remove(startPoint);
                startPoint = null;
                updateRoute();
            });
        } else if (!endPoint) {
            // Установка конечной точки
            endPoint = new ymaps.Placemark(coords, {
                hintContent: 'Конечная точка'
            }, {
                preset: 'islands#redDotIcon', // Красная точка
                draggable: true // Разрешить перемещение
            });
            map.geoObjects.add(endPoint);

            // Удаление конечной точки при клике
            endPoint.events.add('click', function () {
                map.geoObjects.remove(endPoint);
                endPoint = null;
                updateRoute();
            });

            // Создание маршрута
            updateRoute();
        } else {
            // Добавление контрольной точки
            const waypoint = new ymaps.Placemark(coords, {
                hintContent: 'Контрольная точка'
            }, {
                preset: 'islands#blueDotIcon', // Синяя точка
                draggable: true // Разрешить перемещение
            });
            map.geoObjects.add(waypoint);
            waypoints.push(waypoint);

            // Обновление маршрута при перемещении контрольной точки
            waypoint.events.add('dragend', function () {
                updateRoute();
            });

            // Удаление контрольной точки при клике
            waypoint.events.add('click', function () {
                map.geoObjects.remove(waypoint);
                waypoints = waypoints.filter(w => w !== waypoint);
                updateRoute();
            });

            // Создание маршрута
            updateRoute();
        }
    });

    // Обработчик кнопки "Создать маршрут"
    document.getElementById('create-route-form').addEventListener('submit', function (e) {
        e.preventDefault();

        const route_name = document.getElementById('route-name').value;
        const route_place = document.getElementById('route-place').value;
        const route_description = document.getElementById('route-description').value;

        if (route_name.length <= 1 || route_place.length <= 1) {
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
        const waypointsCoords = waypoints.map(w => w.geometry.getCoordinates().join(',')).join('|');

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
        rtt: 'pd' // Пешеходный маршрут
    });

    if (waypointsCoords) {
        params.set('rtext', `${startCoords}~${waypointsCoords.split('|').join('~')}~${endCoords}`);
    }
    const url = `${baseUrl}&${params.toString()}`;

    const urlSafeString = encodeURIComponent(url);

    return urlSafeString

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
        location.reload()
    })
    .catch(error => {
        console.error('Ошибка:', error);
    });
}