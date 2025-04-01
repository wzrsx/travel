let isOpen = false;
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
  // Закрытие при клике вне области
  document.addEventListener('click', function(event) {
    const dropdown = document.getElementById("profileDropdown");
    const profileElement = document.querySelector('.profile-dropdown');
    
    // Если клик был вне области профиля
    if (!profileElement.contains(event.target)) {
      document.getElementById("strelka").classList.remove("rotate-strelka");
      isOpen = false;
    }
    else if(dropdown.contains(event.target)){
      document.getElementById("strelka").classList.remove("rotate-strelka");
      isOpen = false;
    }
  });