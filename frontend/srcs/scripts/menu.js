window.onload = checkToken;

function checkToken() {
    console.log("Function check token")
    const token = localStorage.getItem('token');
    console.log(token)
    if (!token) {
        // alert('No token found. Please login.');
        window.location.href = '/';
    }
}
console.log('La page entière est complètement chargée');
document.getElementById("burger").onclick = function () {
    let burger = document.querySelector(".js-burger");
    let nav = document.querySelector(".js-nav");

    nav.classList.toggle("_active");
    burger.classList.toggle("_active");
}

function redirectionPage(path) {
    window.location.href = path;
}
