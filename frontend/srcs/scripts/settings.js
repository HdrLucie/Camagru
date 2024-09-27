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

document.getElementById("signUp").onclick = async function () {
	console.log("\n\nupdate name\n\n")
	var tmpUser = document.getElementById("Username")
	const token = localStorage.getItem('token')

	var login = tmpUser.value

    const response = await fetch("/editUsername", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            "username": login,
            "JWT": token
        })
    });
}