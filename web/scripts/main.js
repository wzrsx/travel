const signInDialog = document.getElementById('signInDialog');
const registrationDialog = document.getElementById('registrationDialog');
const blurDiv = document.getElementById('blurDiv');
// const IP = process.env.IP; перемнные среды

function openSignInDialog() {
    if (registrationDialog.open) {
        registrationDialog.close(); // Закрываем диалог регистрации, если он открыт
    }
    if (!signInDialog.open) {
        blurDiv.classList.add('blur'); // Добавляем класс размытия только если диалог не открыт
        signInDialog.showModal(); // Открываем диалог авторизации
    }
}

function authorize(){
    const email = document.getElementById('signInEmailInput').value;
    const password = document.getElementById('signInPasswordInput').value;

    const data = {
        email: email,
        password: password
    };
    
    console.log('Sending data:', data);

    fetch('http://localhost:8080/authorize', {
        method: 'POST', // Метод запроса
        headers: {
            'Content-Type': 'application/json' 
        },
        body: JSON.stringify(data) // Преобразуем объект в JSON-строку
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        location.reload();
    })
    .catch((error) => {
        console.log('Error:', error); // Обработка ошибок
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
function registration(){
    const email = document.getElementById('registrationEmailInput').value;
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    const password = document.getElementById('registrationPasswordInput').value;
    const repeat_password = document.getElementById('registrationPasswordInputRepeat').value;
    const specialCharRegex = /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/;
    if (!emailRegex.test(email)){
        //condition returning error with incorrect email format
        return;
    }
    if (!password.length >= 5){
        //condition returning error size of password
        return;
    }
    if (!specialCharRegex.test(password)){
        //condition returning error special simbol absence
        return;
    }
    if (password != repeat_password) {
        //condition returning error not matching passwords 
        return;
    }
    const data = {
        email: email,
        password: password
    };
    fetch("http://localhost:8080/registration"),{
        method: 'POST',
        headers: {
            'Content-Type':'application/json'
        },
        body: JSON.stringify(data)
    }
    .then (response =>{
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        location.reload();
    })
    .catch ((error)=>{
        console.log('Error: ', error);
    })
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
