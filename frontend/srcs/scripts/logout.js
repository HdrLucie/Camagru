// Fonction pour vérifier si le token est présent
window.onload = checkToken;

function checkToken() {
	console.log("Function check token")
	const token = localStorage.getItem('token');
	if (!token) {
		window.location.href = '/';
	}
}

document.getElementById("logout").onclick = async function () {
	const token = localStorage.getItem('token')
	console.log("Logout function\n" + token)
	try {
		const response = await fetch("/logout", {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
				"Authorization": `Bearer ${token}`,
			},
		});
        const responseText = await response.text();
        let data;
        if (responseText) {
            data = JSON.parse(responseText);
        }
		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`);
		}
		if (response.status === 200) {
			localStorage.clear()
			window.location.href = '/'
		} else {
			alert(`Error logout message: ${data.message}`);
		}
	} catch (error) {
		console.error("Error: ", error);
		alert(`Error logout: ${error.message}`);
	}
}