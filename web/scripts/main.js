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
