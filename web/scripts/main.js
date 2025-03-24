const signInDialog = document.getElementById('signInDialog');
const registrationDialog = document.getElementById('registrationDialog');
const forgetPassDialog = document.getElementById('forgetPassDialog');
const blurDiv = document.getElementById('blurDiv');
const EMAIL_REGEXP = /^(([^<>()[\].,;:\s@"]+(\.[^<>()[\].,;:\s@"]+)*)|(".+"))@(([^<>()[\].,;:\s@"]+\.)+[^<>()[\].,;:\s@"]{2,})$/iu;
//инпуты
const emailSignIn = document.getElementById('signInEmailInput');
const emailRegistration = document.getElementById('registrationEmailInput');
const passwordSignIn = document.getElementById('signInPasswordInput');
const passwordRegistration = document.getElementById('registrationPasswordInput');
const passwordrepeatRegistration = document.getElementById('registrationPasswordInputRepeat');
const usernameRegistration = document.getElementById('registrationNameInput');
//ошибки регистрации
const nameRegistrationError = document.getElementById('nameRegistrationError');
const emailRegistrationError = document.getElementById('emailRegistrationError');
const passwordRegistrationError = document.getElementById('passRegistrationError');
const repeatPasswordRegistrationError = document.getElementById('repeatPassRegistrationError');

function openSignInDialog() {
    if (registrationDialog.open) {
        registrationDialog.close(); // Закрываем диалог регистрации, если он открыт
    }
    if (!signInDialog.open) {
        blurDiv.classList.add('blur'); // Добавляем класс размытия только если диалог не открыт
        signInDialog.showModal(); // Открываем диалог авторизации
        // Сбрасываем цвет рамки поля email при открытии диалога если она пустая
        if(!emailSignIn.value.trim()){
            emailSignIn.style.borderColor = ''; // Убираем цвет рамки
        }
    }
}

window.onload = function() {
    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.has('openLoginDialog')) {
        openSignInDialog();
    }
};

function authorize(){
    // Сброс ошибок
    const emailSignInError = document.getElementById('emailSignInError');
    const passwordSignInError = document.getElementById('passwordSignInError');

    emailSignInError.style.display = 'none';
    passwordSignInError.style.display = 'none';

    const emailValue = emailSignIn.value.trim(); // Убираем пробелы
    const passwordValue = passwordSignIn.value.trim(); // Убираем пробелы

    if (!emailValue) {
        emailSignInError.innerText = 'Пожалуйста, введите почту.';
        emailSignInError.style.display = 'block';
        return;
    }

    if (!passwordValue) {
        passwordSignInError.innerText = 'Пожалуйста, введите пароль.';
        passwordSignInError.style.display = 'block';
        return;
    }
    if(!isEmailValid(emailSignIn.value)){
        emailSignInError.innerText = 'Неккоректный формат почты.';
        emailSignInError.style.display = 'block';
        return;
    }

    const data = {
        email: emailValue,
        password: passwordValue
    };
    
    console.log('Sending data:', data);

    fetch('http://localhost:5050/authorize', {
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
        location.href = "/";
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
        if(!emailRegistration.value.trim()){
            emailRegistration.style.borderColor = ''; // Убираем цвет рамки
        }
    }
}

function openForgetPassDialog() {
    // Закрываем диалог авторизации, если он открыт
    if (signInDialog.open) {
        signInDialog.close();
    }
    // Закрываем диалог регистрации, если он открыт
    if (registrationDialog.open) {
        registrationDialog.close();
    }
    // Добавляем класс размытия, если диалог восстановления пароля не открыт
    if (!forgetPassDialog.open) {
        blurDiv.classList.add('blur');
        forgetPassDialog.showModal(); // Открываем диалог восстановления пароля
    }
}

function registration(){
    //ДОБАВИТЬ СБРОС ОШИБОК
    const username = usernameRegistration.value;
    const email = emailRegistration.value;
    const password = passwordRegistration.value;
    const repeat_password = passwordrepeatRegistration.value;
    const specialCharRegex = /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/;

    if (!isEmailValid(email)) {
        emailRegistrationError.innerText = 'Ошибка: Некорректный формат email.';
        emailRegistrationError.style.display = 'block';
        return;
    }

    // Проверка длины пароля
    if (password.length < 5) {
        passwordRegistrationError.innerText = 'Ошибка: Пароль должен содержать не менее 5 символов.';
        passwordRegistrationError.style.display = 'block';
        return;
    }

    // Проверка наличия специального символа в пароле
    if (!specialCharRegex.test(password)) {
        passwordRegistrationError.innerText = 'Ошибка: Пароль должен содержать хотя бы один специальный символ.';
        passwordRegistrationError.style.display = 'block';
        return;
    }

    // Проверка совпадения паролей
    if (password !== repeat_password) {
        repeatPasswordRegistrationError.innerText = 'Ошибка: Пароли не совпадают.';
        repeatPasswordRegistrationError.style.display = 'block';
        return;
    }
    const data = {
        username: username,
        email: email,
        password: password
    };
    fetch("http://localhost:5050/registration",{
        method: 'POST',
        headers: {
            'Content-Type':'application/json'
        },
        body: JSON.stringify(data)
    })
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
    if (!registrationDialog.open && !forgetPassDialog.open) {
        blurDiv.classList.remove('blur'); // Удаляем класс размытия, если ни один диалог не открыт
    }
});

registrationDialog.addEventListener('close', () => {
    if (!signInDialog.open && !forgetPassDialog.open) {
        blurDiv.classList.remove('blur'); // Удаляем класс размытия, если ни один диалог не открыт
    }
});

forgetPassDialog.addEventListener('close', () => {
    if (!signInDialog.open && !registrationDialog.open) {
        blurDiv.classList.remove('blur'); // Удаляем класс размытия, если ни один диалог не открыт
    }
});

function onInput(event) {
    const field = event.target;
    if (isEmailValid(field.value)) {
        field.style.borderColor = 'green';
    } else {
        field.style.borderColor = 'red';
    }
}
function isEmailValid(value) {
    return EMAIL_REGEXP.test(value);
}
emailSignIn.addEventListener('input', onInput);
emailRegistration.addEventListener('input', onInput);