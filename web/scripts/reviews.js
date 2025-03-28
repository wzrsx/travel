const scrollable = document.getElementById('scrollable');
const modal = document.getElementById('modal');
const modalImg = document.getElementById('modal-img');
const closeBtn = document.getElementById('closeButtonModalImg');
const blurDiv = document.getElementById('blurDiv');
const nextBtn = document.getElementById('nextButtonModalImg');
const prevBtn = document.getElementById('prevButtonModalImg');

let isDown = false;
let startX;
let scrollLeft;
let currentIndexImg = 0;

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

// Работа с img
function showImage(index) {
    const thumbnails = document.querySelectorAll('.thumbnail');
    if (index < 0) {
        currentIndex = thumbnails.length - 1; 
    } else if (index >= thumbnails.length) {
        currentIndex = 0; 
    } else {
        currentIndex = index; 
    }
    modalImg.src = thumbnails[currentIndex].src; 
}
const thumbnails = document.querySelectorAll('.thumbnail');
thumbnails.forEach((thumbnail, index) => {
    thumbnail.addEventListener('click', function() {
        blurDiv.classList.add('blur');
        modal.showModal();
        showImage(index);
    });
});
// Обработчики событий для стрелок
nextBtn.addEventListener('click', function() {
    showImage(currentIndex + 1); // Переход к следующему изображению
});

prevBtn.addEventListener('click', function() {
    showImage(currentIndex - 1); // Переход к предыдущему изображению
});


closeBtn.addEventListener('click', function() {
    modal.close();
    blurDiv.classList.remove('blur');
});

// Закрытие модального окна при нажатии вне изображения
modal.addEventListener('click', function(event) {
    if (event.target === modal) {
        modal.close();
        blurDiv.classList.remove('blur');
    }
});
modal.addEventListener('close', () => {
    blurDiv.classList.remove('blur'); 
});
let isOpen = false;

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
  