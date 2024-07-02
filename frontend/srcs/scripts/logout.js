// Fonction pour vérifier si le token est présent
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

document.getElementById("burger").onclick = function () {
    let burger = document.querySelector(".js-burger");
    let nav = document.querySelector(".js-nav");

    nav.classList.toggle("_active");
    burger.classList.toggle("_active");
}

document.getElementById("logout").onclick = async function () {
   alert('Logout')
    const token = localStorage.getItem('token')
    localStorage.clear()
    window.location.href = '/'
}