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
const emailForgetPass = document.getElementById('forgetPassEmailInput');
//ошибки регистрации
const nameRegistrationError = document.getElementById('nameRegistrationError');
const emailRegistrationError = document.getElementById('emailRegistrationError');
const passwordRegistrationError = document.getElementById('passRegistrationError');
const repeatPasswordRegistrationError = document.getElementById('repeatPassRegistrationError');
//ошибки входа
const emailSignInError = document.getElementById('emailSignInError');
const passwordSignInError = document.getElementById('passwordSignInError');
//ошибки "забыли пароль"
const emailForgetPassError = document.getElementById('emailForgetPassError');

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
    // Добавляем обработчики событий
    emailSignIn.addEventListener('input', onInput);
    emailRegistration.addEventListener('input', onInput);
    emailForgetPass.addEventListener('input', onInput);
};

function authorize(){
    // Сброс ошибок
    resetSignInErrors();

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
        if (response.ok) {
            // юзер есть
            location.href = "/";
        } else if (response.status === 401) {
            // email есть, пароль неверный
            passwordSignInError.innerText = 'Неверный пароль.';
            passwordSignInError.style.display = 'block';
            passwordSignIn.style.borderColor = 'red';
        } else if (response.status === 404) {
            // email нету в бд
            emailSignInError.innerText = 'Пользователь с таким email не найден.';
            emailSignInError.style.display = 'block';
            emailSignIn.style.borderColor = 'red';
        } else {
            // Обработка других ошибок
            throw new Error('Network response was not ok ' + response.statusText);
        }
    })
    .catch((error) => {
        console.log('Error:', error); 
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
    event.preventDefault(); // Предотвращаем стандартное поведение формы
    resetRegistrationErrors();
    const username = usernameRegistration.value.trim();
    const email = emailRegistration.value.trim();
    const password = passwordRegistration.value;
    const repeat_password = passwordrepeatRegistration.value;
    const specialCharRegex = /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/;

    if (!username) {
        nameRegistrationError.innerText = 'Пожалуйста, введите имя.';
        nameRegistrationError.style.display = 'block';
        return;
    }
    if (!email) {
        emailRegistrationError.innerText = 'Пожалуйста, введите почту.';
        emailRegistrationError.style.display = 'block';
        return;
    }

    if (!password) {
        passwordRegistrationError.innerText = 'Пожалуйста, введите пароль.';
        passwordRegistrationError.style.display = 'block';
        return;
    }
    if (!isEmailValid(email)) {
        emailRegistrationError.innerText = 'Неккоректный формат почты.';
        emailRegistrationError.style.display = 'block';
        return;
    }

    // Проверка длины пароля
    if (password.length < 5) {
        passwordRegistrationError.innerText = 'Пароль должен содержать не менее 5 символов.';
        passwordRegistrationError.style.display = 'block';
        return;
    }

    // Проверка наличия специального символа в пароле
    if (!specialCharRegex.test(password)) {
        passwordRegistrationError.innerText = 'Пароль должен содержать хотя бы один специальный символ.';
        passwordRegistrationError.style.display = 'block';
        return;
    }

    // Проверка совпадения паролей
    if (password !== repeat_password) {
        repeatPasswordRegistrationError.innerText = 'Пароли не совпадают.';
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
function resetPass(){
    emailForgetPassError.style.display = 'none';
    event.preventDefault();
    const email = emailForgetPass.value.trim();
    if(!email){
        emailForgetPassError.innerText = 'Пожалуйста, введите почту.';
        emailForgetPassError.style.display = 'block';
        return;
    }
    if (!isEmailValid(email)) {
        emailForgetPassError.innerText = 'Неккоректный формат почты.';
        emailForgetPassError.style.display = 'block';
        return;
    }
}
function showCodeInput() {
    const dialog = document.getElementById('forgetPassDialog');
    dialog.innerHTML = `
        <input type="text" class="input-field" id="confirmationCodeInput" placeholder="Введите код подтверждения">
        <button type="button" onclick="confirmCode()">Подтвердить</button>
    `;
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
    const errorMessageElement = event.target.nextElementSibling;
    if (isEmailValid(field.value)) {
        field.style.borderColor = 'green';
        errorMessageElement.style.display = 'none';
    } else {
        field.style.borderColor = 'red';
    }
}
function isEmailValid(value) {
    return EMAIL_REGEXP.test(value);
}

//сброс ошибок регистрации
function resetRegistrationErrors(){
    nameRegistrationError.style.display = 'none';
    emailRegistrationError.style.display = 'none';
    passwordRegistrationError.style.display = 'none';
    repeatPasswordRegistrationError.style.display = 'none';
}
//сброс ошибок входа
function resetSignInErrors(){
    emailSignInError.style.display = 'none';
    passwordSignInError.style.display = 'none';
}