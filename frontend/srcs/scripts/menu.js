window.onload = checkToken;

function checkToken() {
    const token = localStorage.getItem('token');
    if (!token) {
        // alert('No token found. Please login.');
        window.location.href = '/';
    }
}
document.getElementById("burger").onclick = function () {
    let burger = document.querySelector(".js-burger");
    let nav = document.querySelector(".js-nav");

    nav.classList.toggle("_active");
    burger.classList.toggle("_active");
}

function redirectionPage(path) {
    window.location.href = path;
}
