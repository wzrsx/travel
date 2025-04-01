const route_name = document.getElementById('route-name');
const route_place = document.getElementById('route-place');
const route_description = document.getElementById('route-description');
window.onload = function () {
    route_name.addEventListener("input", onInput);
    route_place.addEventListener("input", onInput);
    route_description.addEventListener("input", onInput);
}
let myMap; // Делаем карту глобальной переменной
ymaps.ready(init);

function init() {
    myMap = new ymaps.Map('map', {
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
            myMap.geoObjects.remove(multiRoute);
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
            myMap.geoObjects.add(multiRoute);

            // Обработка кликов на точках маршрута
            multiRoute.model.events.add('requestsuccess', function () {
                multiRoute.getWayPoints().each(function (point) {
                    point.events.add('click', function () {
                        const coords = point.geometry.getCoordinates();

                        // Проверка на удаление начальной точки
                        if (startPoint && startPoint.geometry.getCoordinates().join(',') === coords.join(',')) {
                            myMap.geoObjects.remove(startPoint);
                            if (waypoints.length > 0) {
                                startPoint = waypoints[0];
                                myMap.geoObjects.add(startPoint);
                                waypoints.shift();
                            } else {
                                startPoint = null;
                            }
                            updateRoute();
                        }

                        // Проверка на удаление конечной точки
                        if (endPoint && endPoint.geometry.getCoordinates().join(',') === coords.join(',')) {
                            myMap.geoObjects.remove(endPoint);
                            if (waypoints.length > 0) {
                                endPoint = waypoints[waypoints.length - 1];
                                myMap.geoObjects.add(endPoint);
                                waypoints.pop();
                            } else {
                                endPoint = null;
                            }
                            updateRoute();
                        }

                        // Проверка на удаление контрольной точки
                        const waypointToRemove = waypoints.find(w => w.geometry.getCoordinates().join(',') === coords.join(','));
                        if (waypointToRemove) {
                            myMap.geoObjects.remove(waypointToRemove);
                            waypoints = waypoints.filter(w => w !== waypointToRemove);
                            updateRoute();
                        }
                    });
                });
            });
        }
    }

    // Обработчик клика на карте
    myMap.events.add('click', function (e) {
        const coords = e.get('coords');

        if (!startPoint) {
            // Установка начальной точки
            startPoint = new ymaps.Placemark(coords, {
                hintContent: 'Начальная точка'
            }, {
                preset: 'islands#greenDotIcon', // Зеленая точка
                draggable: true // Разрешить перемещение
            });
            myMap.geoObjects.add(startPoint);

            // Удаление начальной точки при клике
            startPoint.events.add('click', function () {
                myMap.geoObjects.remove(startPoint);
                startPoint = null;
                updateRoute();
            });

            updateRoute();
        } else if (!endPoint) {
            // Установка конечной точки
            endPoint = new ymaps.Placemark(coords, {
                hintContent: 'Конечная точка'
            }, {
                preset: 'islands#redDotIcon', // Красная точка
                draggable: true // Разрешить перемещение
            });
            myMap.geoObjects.add(endPoint);

            // Удаление конечной точки при клике
            endPoint.events.add('click', function () {
                myMap.geoObjects.remove(endPoint);
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
            myMap.geoObjects.add(waypoint);
            waypoints.push(waypoint);

            // Обновление маршрута при перемещении контрольной точки
            waypoint.events.add('dragend', function () {
                updateRoute();
            });

            // Удаление контрольной точки при клике
            waypoint.events.add('click', function () {
                myMap.geoObjects.remove(waypoint);
                waypoints = waypoints.filter(w => w !== waypoint);
                updateRoute();
            });

            // Создание маршрута
            updateRoute();
        }
        // После добавления точек скрываем ошибку
        if (startPoint && endPoint) {
            showError(document.getElementById('map'), '', document.getElementById('mapRouteError'));
            document.getElementById('map').style.border = 'none';
            route_name.scrollIntoView({ behavior: 'smooth', block: 'center' });
        }
    });

    // Обработчик кнопки "Создать маршрут"
    document.getElementById('create-route-form').addEventListener('submit', function (e) {
        e.preventDefault();

        if(!validateField(route_place, document.getElementById('placeRouteError'))){
            return; 
        }
        if (!startPoint || !endPoint) {
            showError(document.getElementById('map'), 'Пожалуйста, укажите начальную и конечную точки на карте.', document.getElementById('mapRouteError'));
            return;
        }else{
            showError(document.getElementById('map'), '', document.getElementById('mapRouteError'));
        }
        if(!validateField(route_name, document.getElementById('nameRouteError'))){
            return; 
        }
    
        if(!validateField(route_description, document.getElementById('descriptionRouteError'))){
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
    .then(response => {
        if (!response.ok) {
            return response.text().then(text => {
                throw new Error(text);
            });
        }
        return response.json();
    })
    .then(data => {
        console.log('Ссылка сохранена:', data);
        alert('Маршрут успешно создан!');
        location.reload();
    })
    .catch(error => {
        console.error('Ошибка:', error);
        alert('Ошибка при сохранении маршрута: ' + error.message);
    });
}

const scrollable = document.getElementById('scrollable');
let isDown = false;
//Слайдер фоток
scrollable.addEventListener('mousedown', (e) => {
    isDown = true;
    e.preventDefault();
    scrollable.classList.add('active');
    startX = e.pageX - scrollable.offsetLeft;
    scrollLeft = scrollable.scrollLeft;
});
scrollable.addEventListener('mouseleave', () => {
    isDown = false;
    scrollable.classList.remove('active');
});

scrollable.addEventListener('mouseup', () => {
    isDown = false;
    scrollable.classList.remove('active');
});
scrollable.addEventListener('mousemove', (e) => {
    if (!isDown) return; 
    e.preventDefault();
    const x = e.pageX - scrollable.offsetLeft;
    const walk = (x - startX) * 2; 
    scrollable.scrollLeft = scrollLeft - walk;
});

scrollable.addEventListener('wheel', (e) => {
    e.preventDefault(); 
    scrollable.scrollLeft += e.deltaY; 
});

document.addEventListener('DOMContentLoaded', function() {
    const fileInput = document.getElementById('route-photo');
    const photosContainer = document.getElementById('scrollable');
    let filesQueue = 0; // Счетчик файлов в обработке
    
    fileInput.addEventListener('change', function(e) {
        if (this.files && this.files.length > 0) {
            filesQueue = this.files.length; // Устанавливаем общее количество файлов
            
            Array.from(this.files).forEach(file => {
                if (file.type.startsWith('image/')) {
                    const reader = new FileReader();
                    
                    reader.onload = function(e) {
                        createPhotoItem(e.target.result);
                        filesQueue--; // Уменьшаем счетчик после загрузки
                        
                        // Обновляем счетчик только когда все файлы обработаны
                        if (filesQueue === 0) {
                            updateCounter();
                        }
                    }
                    
                    reader.onerror = function() {
                        filesQueue--; // Учитываем ошибки загрузки
                        if (filesQueue === 0) {
                            updateCounter();
                        }
                    }
                    
                    reader.readAsDataURL(file);
                } else {
                    filesQueue--; // Пропускаем не-изображения
                }
            });
        }
    });
    
    function createPhotoItem(imageSrc) {
        const photoItem = document.createElement('div');
        photoItem.className = 'photo-item';
        
        photoItem.innerHTML = `
            <img src="${imageSrc}" alt="" class="thumbnail">
            <span class="delete-photo">&times;</span>
        `;
        
        photosContainer.appendChild(photoItem);
        
        photoItem.querySelector('.delete-photo').addEventListener('click', function() {
            photoItem.remove();
            updateCounter();
        });
    }
    
    function updateCounter() {
        const totalPhotos = photosContainer.querySelectorAll('.photo-item').length;
        document.querySelector('.custom-file-input').textContent = 
            totalPhotos > 0 ? `Выбрано ${totalPhotos} фото` : 'Выберите файлы';
    }
    
    // Инициализация существующих элементов
    document.querySelectorAll('#scrollable .photo-item .delete-photo').forEach(btn => {
        btn.addEventListener('click', function() {
            this.closest('.photo-item').remove();
            updateCounter();
        });
    });
    
    updateCounter();
});

function validateField(field, errorMessageElement) {
    const value = field.value;
    let isValid = true;
    let errorMessage = '';
    
    if (!value) {
        errorMessage = 'Это поле обязательно для заполнения';
        isValid = false;
    }
    showError(field, errorMessage, errorMessageElement);
    
    return isValid;
}
function showError(field, message, errorElement) {
    if (message) {
        //Показываем ошибку
        errorElement.style.display = 'block';
        errorElement.innerText = message;
        field.style.border = '3px solid #752828';
        //прокрутка к полю
        field.scrollIntoView({ behavior: 'smooth', block: 'center' });
    } else {
        // Убираем ошибку
        errorElement.style.display = 'none';
        field.style.borderColor = '#506C56';
    }
}
function onInput(event) {
    const field = event.target;
    const errorMessageElement = event.target.nextElementSibling;
    if (field.value) {
        showError(field, '', errorMessageElement);
        if(field.value.length > 1){
            findCity(field.value);
            field.classList.add('delete-border-radius');
        }
    } else {
        showError(field, 'Это поле обязательно для заполнения', errorMessageElement);
        document.getElementById('suggestions-container').style.display = 'none';
        field.classList.remove('delete-border-radius');
    }
  }
  /*Подсказки при выборе места */
function findCity(query){
    var url = "http://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/address";
    var token = "6edd25e46a970c5f63c88a772177b0e4cf5a57b7";

    var options = {
        method: "POST",
        mode: "cors",
        headers: {
            "Content-Type": "application/json",
            "Accept": "application/json",
            "Authorization": "Token " + token
        },
        body: JSON.stringify({
            query: query,
            "from_bound": { "value": "country" },
            "to_bound": { "value": "settlement" },
            "locations": [{ "country": "Россия" }]
        })
    }

    fetch(url, options)
    .then(response => response.json())
    .then(result => 
    {
        if (!result.suggestions) return;
        // Фильтруем и обрабатываем подсказки
        console.log(result.suggestions);
        const filtered = result.suggestions.filter(suggestion => {
            return suggestion.data && suggestion.data.settlement_type !== "тер";
        });
        showSuggestions(filtered);
    })
    .catch(error => console.log("Ошибка:", error));
}
function showSuggestions(suggestions) {
    const container = document.getElementById('suggestions-container');
    container.innerHTML = '';
    
    if (suggestions.length === 0) {
        container.style.display = 'none';
        return;
    }
    
    suggestions.forEach(suggestion => {
        const div = document.createElement('div');
        div.className = 'suggestion-item';
        div.textContent = suggestion.value;

        div.addEventListener('click', () => {
            document.getElementById('route-place').value = suggestion.value;
            container.style.display = 'none';
            route_place.classList.remove('delete-border-radius');
            document.getElementById('map').scrollIntoView({ behavior: 'smooth', block: 'center' });
            handleSuggestionSelect(suggestion);
        });
        container.appendChild(div);
    });
    
    container.style.display = 'block';
}
//при выборе места из списка -> установить карту в эту точку
function handleSuggestionSelect(selectedSuggestion) {
    const coords = [
        parseFloat(selectedSuggestion.data.geo_lat),
        parseFloat(selectedSuggestion.data.geo_lon)
    ];
    console.log(coords);
    if (!myMap) {
        console.error('Карта не инициализирована');
        return;
    }
    myMap.setCenter(coords, 10);
}