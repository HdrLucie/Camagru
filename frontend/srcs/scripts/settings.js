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

function redirectionPage(path) {
    alert('Redirection page')
    window.location.href = path;
}

// document.getElementById("setUsername").onclick = async function () {
// 	console.log("\n\nupdate name\n\n")
// 	var tmpUser = document.getElementById("username")
//     // console.log(tmpUser)
// 	const token = localStorage.getItem('token')

// 	var login = tmpUser.value
//     console.log("login :" + login)
//     const response = fetch("/editUsername", {
//         method: "POST",
//         headers: {
//             "Authorization": `Bearer ${token}`,
//             "Content-Type": "application/json",
//         },
//         body: JSON.stringify({
//             "login": login,
//         })
//     });
// }

// document.getElementById("setPassword").onclick = async function () {
// 	console.log("\n\nupdate name\n\n")
// 	var tmpUser = document.getElementById("password")
//     // console.log(tmpUser)
// 	const token = localStorage.getItem('token')

// 	var password = tmpUser.value
//     console.log("password :" + password)
//     const response = fetch("/editPassword", {
//         method: "POST",
//         headers: {
//             "Authorization": `Bearer ${token}`,
//             "Content-Type": "application/json",
//         },
//         body: JSON.stringify({
//             "password": password,
//         })
//     });
// }

// document.getElementById("setEmail").onclick = async function () {
// 	console.log("\n\nupdate mail\n\n")
// 	var tmpUser = document.getElementById("email")
//     // console.log(tmpUser)
// 	const token = localStorage.getItem('token')

// 	var email = tmpUser.value
//     console.log("Email :" + email)
//     const response = fetch("/editEmail", {
//         method: "POST",
//         headers: {
//             "Authorization": `Bearer ${token}`,
//             "Content-Type": "application/json",
//         },
//         body: JSON.stringify({
//             "email": email,
//         })
//     });
// }

document.addEventListener('DOMContentLoaded', function() {
    const showPhotos = document.getElementById('show-photos');
    const showSettings = document.getElementById('show-settings');
    const photosList = document.getElementById('photos-list');
    const settings = document.getElementById('settings');

    function showContent(contentToShow) {
        [photosList, settings].forEach(content => {
            content.classList.remove('active');
        });
        contentToShow.classList.add('active');
    }

    showPhotos.addEventListener('click', function(e) {
        e.preventDefault();
        showContent(photosList);
    });

    showSettings.addEventListener('click', function(e) {
        e.preventDefault();
        showContent(settings);
    });

    // La liste des photos est affichée par défaut
    showContent(photosList);
});