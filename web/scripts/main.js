const signInDialog = document.getElementById('signInDialog');
const registrationDialog = document.getElementById('registrationDialog');
const blurDiv = document.getElementById('blurDiv');

function openSignInDialog() {
    if (registrationDialog.open) {
        registrationDialog.close(); // Закрываем диалог регистрации, если он открыт
    }
    if (!signInDialog.open) {
        blurDiv.classList.add('blur'); // Добавляем класс размытия только если диалог не открыт
        signInDialog.showModal(); // Открываем диалог авторизации
    }
}

function Authorize(){
    const email = document.getElementById('signInEmailInput').value;
    const password = document.getElementById('signInPasswordInput').value;

    const data = {
        email: email,
        password: password
    };
    
    console.log('Sending data:', data);

    fetch('https://127.0.0.1/authorize', {
        method: 'POST', // Метод запроса
        headers: {
            'Content-Type': 'application/json' // Указываем, что отправляем JSON
        },
        body: JSON.stringify(data) // Преобразуем объект в JSON-строку
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        return response.json(); // Преобразуем ответ в JSON
    })
    .then(data => {
        console.log('Success:', data); // Обработка успешного ответа
    })
    .catch((error) => {
        console.error('Error:', error); // Обработка ошибок
    });
}

function openRegistrationDialog() {
    if (signInDialog.open) {
        signInDialog.close(); // Закрываем диалог авторизации, если он открыт
    }
    if (!registrationDialog.open) {
        blurDiv.classList.add('blur'); // Добавляем класс размытия только если диалог не открыт
        registrationDialog.showModal(); // Открываем диалог регистрации
    }
}

// Закрытие диалогов и удаление размытия
signInDialog.addEventListener('close', () => {
    if (!registrationDialog.open) {
        blurDiv.classList.remove('blur'); // Удаляем класс размытия, если ни один диалог не открыт
    }
});

registrationDialog.addEventListener('close', () => {
    if (!signInDialog.open) {
        blurDiv.classList.remove('blur'); // Удаляем класс размытия, если ни один диалог не открыт
    }
});
