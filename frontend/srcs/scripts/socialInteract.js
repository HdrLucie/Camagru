document.addEventListener('DOMContentLoaded', () => {
	checkToken();
    // loadUserData();
});

async function checkToken() {
    const token = localStorage.getItem('token');
    console.log("Function check token")
    console.log(token)
    if (!token) {
        // alert('No token found. Please login.');
        window.location.href = '/';
    }
}

async function getUser() {
    const token = localStorage.getItem('token');
    try {
        const response = await fetch("/getUser", {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json",
            },
        });
        const userData = await response.json();
		console.log(userData);
		return userData;
    } catch (error) {
        console.error("Erreur:", error);
        return null;
    }
}

document.getElementById("sendLikes").onclick = async function () {
	console.log("Likes button");
	const photoId = window.location.pathname.split("/").pop();
	console.log(photoId);
	const user = await getUser();
	console.log(user);
	try {
		const response = await fetch("/sendLikes", {
			method: "POST",
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify({
				"Username": user.username,
				"Id": user.id,
				"Photo": Number(photoId),
			})
		});
		const data = await response.json();
	} catch (error) {
		console.error("Error:", error);
	}
}
