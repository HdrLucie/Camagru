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

document.getElementById("setUsername").onclick = async function () {
	console.log("\n\nupdate name\n\n")
	var tmpUser = document.getElementById("username")
    // console.log(tmpUser)
	const token = localStorage.getItem('token')

	var login = tmpUser.value
    console.log("login :" + login)
    const response = fetch("/editUsername", {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${token}`,
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            "login": login,
        })
    });
}

document.getElementById("setPassword").onclick = async function () {
	console.log("\n\nupdate name\n\n")
	var tmpUser = document.getElementById("password")
    // console.log(tmpUser)
	const token = localStorage.getItem('token')

	var password = tmpUser.value
    console.log("password :" + password)
    const response = fetch("/editPassword", {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${token}`,
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            "password": password,
        })
    });
}

document.getElementById("setEmail").onclick = async function () {
	console.log("\n\nupdate mail\n\n")
	var tmpUser = document.getElementById("email")
    // console.log(tmpUser)
	const token = localStorage.getItem('token')

	var email = tmpUser.value
    console.log("Email :" + email)
    const response = fetch("/editEmail", {
        method: "POST",
        headers: {
            "Authorization": `Bearer ${token}`,
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            "email": email,
        })
    });
}