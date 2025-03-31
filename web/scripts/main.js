const signInDialog = document.getElementById("signInDialog");
const registrationDialog = document.getElementById("registrationDialog");
const forgetPassDialog = document.getElementById("forgetPassDialog");
const changePassDialog = document.getElementById("changePassDialog");
const blurDiv = document.getElementById("blurDiv");
const EMAIL_REGEXP =
  /^(([^<>()[\].,;:\s@"]+(\.[^<>()[\].,;:\s@"]+)*)|(".+"))@(([^<>()[\].,;:\s@"]+\.)+[^<>()[\].,;:\s@"]{2,})$/iu;
//инпуты
const emailSignIn = document.getElementById("signInEmailInput");
const emailRegistration = document.getElementById("registrationEmailInput");
const passwordSignIn = document.getElementById("signInPasswordInput");
const passwordRegistration = document.getElementById(
  "registrationPasswordInput"
);
const passwordrepeatRegistration = document.getElementById(
  "registrationPasswordInputRepeat"
);
const usernameRegistration = document.getElementById("registrationNameInput");
const emailForgetPass = document.getElementById("forgetPassEmailInput");
const confirmationCode = document.getElementById("confirmationCode");
const newPass = document.getElementById("newPassInput");
const changePassLast = document.getElementById("changePassLastPassInput");
const changePassNew = document.getElementById("changePassNewPassInput");
//ошибки регистрации
const nameRegistrationError = document.getElementById("nameRegistrationError");
const emailRegistrationError = document.getElementById(
  "emailRegistrationError"
);
const passwordRegistrationError = document.getElementById(
  "passRegistrationError"
);
const repeatPasswordRegistrationError = document.getElementById(
  "repeatPassRegistrationError"
);
//ошибки входа
const emailSignInError = document.getElementById("emailSignInError");
const passwordSignInError = document.getElementById("passwordSignInError");
//ошибки "забыли пароль"
const emailForgetPassError = document.getElementById("emailForgetPassError");
const codeForgetPassError = document.getElementById("codeForgetPassError");
const passForgetPassError = document.getElementById("passForgetPassError");
//ошибки изменения пароля
const lastPassChangePassError = document.getElementById("lastPassChangePassError");
const newPassChangePassError = document.getElementById("newPassChangePassError");

let isOpen = false;
function openSignInDialog() {
  if (registrationDialog.open) {
    registrationDialog.close(); // Закрываем диалог регистрации, если он открыт
  }
  if (!signInDialog.open) {
    blurDiv.classList.add("blur"); // Добавляем класс размытия только если диалог не открыт
    signInDialog.showModal(); // Открываем диалог авторизации
    // Сбрасываем цвет рамки поля email при открытии диалога если она пустая
    if (!emailSignIn.value.trim()) {
      emailSignIn.style.borderColor = ""; // Убираем цвет рамки
    }
  }
}

window.onload = function () {
  const urlParams = new URLSearchParams(window.location.search);
  if (urlParams.has("openLoginDialog")) {
    openSignInDialog();
  }
  // Добавляем обработчики событий
  emailSignIn.addEventListener("input", onInput);
  emailRegistration.addEventListener("input", onInput);
  emailForgetPass.addEventListener("input", onInput);
};

function authorize() {
  // Сброс ошибок
  resetSignInErrors();

  const emailValue = emailSignIn.value.trim(); // Убираем пробелы
  const passwordValue = passwordSignIn.value.trim(); // Убираем пробелы

  if (!emailValue) {
    showError(emailSignInError, "Пожалуйста, введите почту.");
    return;
  }

  if (!passwordValue) {
    showError(passwordSignInError, "Пожалуйста, введите пароль.");
    return;
  }
  if (!isEmailValid(emailSignIn.value)) {
    showError(emailSignInError, "Неккоректный формат почты.");
    return;
  }

  const data = {
    email: emailValue,
    password: passwordValue,
  };

  console.log("Sending data:", data);

  fetch("http://localhost:5050/authorize", {
    method: "POST", // Метод запроса
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data), // Преобразуем объект в JSON-строку
  })
    .then((response) => {
      if (response.ok) {
        // юзер есть
        location.href = "/";
      } else if (response.status === 401) {
        // email есть, пароль неверный
        showError(passwordSignInError, "Неверный пароль");
        passwordSignIn.style.borderColor = "red";
      } else if (response.status === 404) {
        // email нету в бд
        showError(emailSignInError, "Пользователь с таким email не найден.");
        emailSignIn.style.borderColor = "red";
      } else {
        // Обработка других ошибок
        throw new Error("Network response was not ok " + response.statusText);
      }
    })
    .catch((error) => {
      console.log("Error:", error);
    });
}

function openRegistrationDialog() {
  if (signInDialog.open) {
    signInDialog.close(); // Закрываем диалог авторизации, если он открыт
  }
  if (!registrationDialog.open) {
    blurDiv.classList.add("blur"); // Добавляем класс размытия только если диалог не открыт
    registrationDialog.showModal(); // Открываем диалог регистрации
    if (!emailRegistration.value.trim()) {
      emailRegistration.style.borderColor = ""; // Убираем цвет рамки
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
    blurDiv.classList.add("blur");
    forgetPassDialog.showModal(); // Открываем диалог восстановления пароля
    if (!emailForgetPass.value.trim()) {
      emailForgetPass.style.borderColor = ""; // Убираем цвет рамки
    }
  }
}

function registration() {
  event.preventDefault(); // Предотвращаем стандартное поведение формы
  resetRegistrationErrors();
  const username = usernameRegistration.value.trim();
  const email = emailRegistration.value.trim();
  const password = passwordRegistration.value;
  const repeat_password = passwordrepeatRegistration.value;

  if (!username) {
    showError(nameRegistrationError, "Пожалуйста, введите имя.");
    return;
  }
  if (!email) {
    showError(emailRegistrationError, "Пожалуйста, введите почту.");
    return;
  }

  if (!isEmailValid(email)) {
    showError(emailRegistrationError, "Неккоректный формат почты.");
    return;
  }
  if (!isPassValid(password, passwordRegistrationError, passwordRegistration)) {
    return;
  }
  if (!repeat_password) {
    showError(repeatPasswordRegistrationError, "Пожалуйста, повторите пароль.");
    return;
  }
  // Проверка совпадения паролей
  if (password !== repeat_password) {
    showError(repeatPasswordRegistrationError, "Пароли не совпадают.");
    return;
  }
  const data = {
    username: username,
    email: email,
    password: password,
  };
  fetch("http://localhost:5050/registration", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  })
    .then((response) => {
      if (response.ok) {
        data = {
          email: email,
        };
        fetch("http://localhost:5050/send_to_email/pass_code", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(data),
        }) 
        .then((response) => {
          if(!response.ok){
            showError(repeatPasswordRegistrationError, "Ошибка отправки письма.");
            console.log();
            return;
          }
          // ОТКРЫТИЕ ОКНА С ВВОДОМ КОДА ДЛЯ РЕГИСТРАЦИИ.

        })
        return;
      }
      if (response.status == 409) {
        showError(
          emailRegistrationError,
          "Пользователь с таким email уже зарегистрирован."
        );
        emailRegistration.style.borderColor = "red";
        return;
      } 

      if (response.status == 400) {
        throw new Error("Network response was not ok " + response.statusText);
      }
    })
    .catch((error) => {
      console.log("Error: ", error);
    });
}
function resetPass() {
  emailForgetPassError.style.display = "none";
  event.preventDefault();
  const email = emailForgetPass.value.trim();
  if (!email) {
    showError(emailForgetPassError, "Пожалуйста, введите почту.");
    return;
  }
  if (!isEmailValid(email)) {
    showError(emailForgetPassError, "Неккоректный формат почты.");
    return;
  }
  const data = {
    email: email,
  };
  fetch("http://localhost:5050/check/email", {
    method: "POST", // Метод запроса
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data), // Преобразуем объект в JSON-строку
  }).then((response) => {
    if (response.ok) {
      fetch("http://localhost:5050/send_to_email/pass_code", {
        method: "POST", // Метод запроса
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data), // Преобразуем объект в JSON-строку
      }).then((response) => {
        if (!response.ok) {
          showError(emailForgetPassError, "Ошибка отправки письма");
          return;
        }
        showCodeInput();
      });
    } else if (response.status == 409) {
      showError(
        emailForgetPassError,
        "Пользователя с такой почтой не существует."
      );
    }
  });
}
function showCodeInput() {
  emailForgetPass.style.display = "none";
  document.getElementById("supportTextforgetPassForm").innerText =
    "Введите код подтверждения";
  confirmationCode.style.display = "flex";
  document.getElementById("resetPassButton").style.display = "none";
}

function showPassInput(email) {
  confirmationCode.style.display = "none";
  document.getElementById("supportTextforgetPassForm").innerText =
    "Введите новый пароль";
  newPass.style.display = "block";

  setNewPassButton = document.getElementById("setNewPassButton");
  setNewPassButton.style.display = "block";

  setNewPassButton.addEventListener("click", (event) => {
    event.preventDefault(); // Предотвращаем стандартное поведение

    const data = {
      email: email,
      password: newPass.value,
    };
    setNewPass(data, event);
  });
}
function setNewPass(data, event) {
  if (event) {
    event.preventDefault();
  }

  if (!isPassValid(newPass.value, passForgetPassError, newPass)) {
    console.log("asdasdasd");
    return;
  }
  fetch("http://localhost:5050/change/password", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  }).then((response) => {
    switch (response.status) {
      case 409:
        console.log("Код не был введен.");
        break;
      case 423:
        showError(passForgetPassError, "Время кода было просрочено.");
        console.log("Время кода было просрочено.");
        break;
      case 408:
        console.log("Попытка сменить пароль при блокировке.");
        break;
      default:
        forgetPassDialog.close();
        break;
    }
  });
}
// Закрытие диалогов и удаление размытия
signInDialog.addEventListener("close", () => {
  if (!registrationDialog.open && !forgetPassDialog.open) {
    blurDiv.classList.remove("blur"); // Удаляем класс размытия, если ни один диалог не открыт
  }
});

registrationDialog.addEventListener("close", () => {
  if (!signInDialog.open && !forgetPassDialog.open) {
    blurDiv.classList.remove("blur"); // Удаляем класс размытия, если ни один диалог не открыт
  }
});

forgetPassDialog.addEventListener("close", () => {
  if (!signInDialog.open && !registrationDialog.open) {
    blurDiv.classList.remove("blur"); // Удаляем класс размытия, если ни один диалог не открыт
  }
});

function onInput(event) {
  const field = event.target;
  const errorMessageElement = event.target.nextElementSibling;
  if (isEmailValid(field.value)) {
    field.style.borderColor = "green";
    errorMessageElement.style.display = "none";
  } else {
    field.style.borderColor = "red";
  }
}
function isEmailValid(value) {
  return EMAIL_REGEXP.test(value);
}

//сброс ошибок регистрации
function resetRegistrationErrors() {
  nameRegistrationError.style.display = "none";
  emailRegistrationError.style.display = "none";
  passwordRegistrationError.style.display = "none";
  repeatPasswordRegistrationError.style.display = "none";
}
//сброс ошибок входа
function resetSignInErrors() {
  emailSignInError.style.display = "none";
  passwordSignInError.style.display = "none";
}
function resetChangePassErrors() {
  newPassChangePassError.style.display = "none";
  lastPassChangePassError.style.display = "none";
  changePassNew.style.borderColor = 'black';
  changePassLast.style.borderColor = 'black';
}
//проверка пароля
function isPassValid(value, field, input) {
  const specialCharRegex = /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/;
  if (!value) {
    showError(field, "Пожалуйста, введите пароль.");
    input.style.borderColor = "red";
    return false;
  }
  // Проверка длины пароля
  if (value.length < 5) {
    showError(field, "Пароль должен содержать не менее 5 символов.");
    input.style.borderColor = "red";
    return false;
  }

  // Проверка наличия специального символа в пароле
  if (!specialCharRegex.test(value)) {
    showError(
      field,
      "Пароль должен содержать хотя бы один специальный символ."
    );
    input.style.borderColor = "red";
    return false;
  }
  return true;
}
function showError(field, text) {
  field.innerText = text;
  field.style.display = "block";
  field.style.color = "red";
}

//переход к некст инпуту кода
function moveToNext(currentInput, nextInputId) {
  // Если текущее поле не пустое, перемещаем фокус на следующее поле
  if (currentInput.value.length >= 1 && nextInputId) {
    document.getElementById(nextInputId).focus();
  }
}

//проверка на число
function isNumberKey(evt) {
  const charCode = evt.which ? evt.which : evt.keyCode;
  // Разрешаем только цифры (0-9)
  if (charCode < 48 || charCode > 57) {
    evt.preventDefault(); // Запрещаем ввод
    return false;
  }
  return true;
}
//проверяем 6 инпутов
function checkAllFilled() {
  const inputs = document.querySelectorAll("#confirmationCode .code-input"); //получаем все инпуты внутри блока
  const allFilled = Array.from(inputs).every(
    (input) => input.value.length === 1
  );
  resetInputStyles(inputs);
  if (allFilled) {
    const code = Array.from(inputs)
      .map((input) => input.value)
      .join("");
    const email = emailForgetPass.value.trim();
    console.log("email:", email);
    body = {
      email: email,
      code: code,
    };
    fetch("http://localhost:5050/check/pass_code", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    })
      .then((response) => {
        if (response.status === 423) {
          showBlockInput(inputs);
          return response.json(); // Преобразуем ответ в JSON
        }
        if (response.status === 409) {
          applyErrorStyles(inputs);
          throw new Error("Invalid code");
        }
        if (response.ok) {
          resetInputStyles(inputs);
          showPassInput(email);
        }
        return Promise.reject();
      })
      .then((data) => {
        if (!data?.unlock_at) {
          throw new Error("Не получено время разблокировки");
        }

        const unlockTime = new Date(data.unlock_at);
        if (isNaN(unlockTime.getTime())) {
          throw new Error(`Некорректный формат времени: ${data.unlock_at}`);
        }

        const updateTimer = () => {
          const now = new Date();
          const diff = unlockTime - now;
          console.log(now);
          if (diff <= 0) {
            clearInterval(interval);
            showCanField(
              codeForgetPassError,
              "Можете пробовать снова.",
              inputs
            );
            return;
          }

          const totalSeconds = Math.floor(diff / 1000);
          const minutes = Math.floor(totalSeconds / 60);
          const seconds = totalSeconds % 60;

          const timeString = `${minutes.toString().padStart(2, "0")}:${seconds
            .toString()
            .padStart(2, "0")}`;
          showError(
            codeForgetPassError,
            `Слишком большое количество попыток, попробуйте снова через: ${timeString}`
          );
        };

        // Первое обновление сразу
        updateTimer();

        // Затем каждую секунду
        const interval = setInterval(updateTimer, 1000);

        return { interval, unlockTime };
      });
  }

  return allFilled;
}

function resetInputStyles(inputs) {
  inputs.forEach((input) => {
    input.style.borderColor = "";
    input.style.background = "";
  });
  codeForgetPassError.style.display = "none";
}

function applyErrorStyles(inputs) {
  showError(codeForgetPassError, "Неверный код.");
  inputs.forEach((input) => {
    input.style.borderColor = "red";
    input.style.background = "#ffcccc";
  });
}
function resetErrorStyles(inputs) {
  codeForgetPassError.style.display = "none";
  inputs.forEach((input) => {
    input.style.borderColor = "green";
    input.style.background = "#e9f7d6";
  });
}

function showCanField(field, text, inputs) {
  field.innerText = text;
  field.style.color = "green";
  field.style.display = "block";

  inputs.forEach((input) => {
    input.style.borderColor = "black";
    input.style.backgroundColor = "#e9f7d6";
    input.style.color = "black";
    input.disabled = false; // Блокируем ввод
    input.style.opacity = "1"; // Затемняем
    input.style.cursor = "auto"; // Меняем курсор
  });
}
function showBlockInput(inputs) {
  inputs.forEach((input) => {
    input.style.borderColor = "#ff0000"; // Яркий красный цвет
    input.disabled = true; // Блокируем ввод (используем input, а не inputs)
    input.style.opacity = "0.6"; // Оптимальное затемнение
    input.style.cursor = "not-allowed"; // Курсор "недоступно"
    input.style.backgroundColor = "#f5f5f5"; // Серый фон
    input.style.color = "#999"; // Серый текст
  });
}
function handleKeyDown(event, inputElement) {
  if (event.key === "Backspace") {
    const currentIndex = parseInt(inputElement.id.replace("inputCode", ""));
    if (inputElement.value === "") {
      // только если input уже пустой
      inputElement.value = "";
      const previousInput = document.getElementById(
        "inputCode" + (currentIndex - 1)
      );
      if (previousInput) {
        previousInput.focus();
      }
      checkAllFilled();
    } else {
      if (currentIndex == 6) {
        const inputs = document.querySelectorAll(
          "#confirmationCode .code-input"
        );
        resetErrorStyles(inputs);
      }
      inputElement.value = ""; // если есть значение - просто очищаем, не переходим
    }
  }
}
//обработка вставки
confirmationCode.addEventListener("paste", function (e) {
  e.preventDefault(); // Отменяем стандартное поведение

  // Получаем текст из буфера обмена
  const pastedText = (e.clipboardData || window.clipboardData).getData("text");

  // Оставляем только цифры и обрезаем до 6 символов
  const digitsOnly = pastedText.replace(/\D/g, "").substring(0, 6);

  // Если длина не 6 символов, игнорируем
  if (digitsOnly.length !== 6) return;

  // Заполняем поля посимвольно
  for (let i = 0; i < 6; i++) {
    const input = document.getElementById(`inputCode${i + 1}`);
    if (input) {
      input.value = digitsOnly[i];
    }
  }

  // Вызываем проверку заполненности
  checkAllFilled();
  document.getElementById("inputCode6").focus();
});
function openProfil() {
  if(!isOpen){
    document.getElementById("strelka").classList.add("rotate-strelka");
    isOpen = true;
  }
  else{
    document.getElementById("strelka").classList.remove("rotate-strelka");
    isOpen = false;
  }
}

function toggleProfileMenu() {
  event.preventDefault();
  openProfil();
  document.getElementById("profileDropdown").classList.toggle("show");
}

// Закрытие меню при клике вне его области
window.onclick = function(event) {
  if (!event.target.matches('#profilUser') && !event.target.matches('#strelka')) {
      var dropdowns = document.getElementsByClassName("dropdown-content");
      for (var i = 0; i < dropdowns.length; i++) {
          var openDropdown = dropdowns[i];
          if (openDropdown.classList.contains('show')) {
              openDropdown.classList.remove('show');
          }
      }
  }
}

function changePass(){
  changePassDialog.showModal();
  blurDiv.classList.add("blur");
}
changePassDialog.addEventListener("close", () => {
  blurDiv.classList.remove("blur");
});
function changePassSave(){
    event.preventDefault();
    resetChangePassErrors();
    const lastPassValue = changePassLast.value;
    const newPassValue = changePassNew.value;
    if(!isPassValid(lastPassValue, lastPassChangePassError, changePassLast)){
      return;
    }
    if(!isPassValid(newPassValue, newPassChangePassError, changePassNew)){
      return;
    }
    //проверить с текущим паролем юзера
    if (lastPassValue !== user_pass) {
      showError(lastPassChangePassError, "Пароль не соовпадает с текущим.");
      changePassLast.style.borderColor = 'red';
      return;
    }
    //если совпадает то отправляем update
}