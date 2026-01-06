document.addEventListener('DOMContentLoaded', () => {
	checkToken();
	displayPicture();
});

function checkToken() {
	const token = localStorage.getItem('token');
	if (!token) {
		window.location.href = '/';
	}
	return token
}

async function getPicture() {
	const token = localStorage.getItem('token');
	try {
		const response = await fetch("/getPicture", {
			method: "GET",
			headers: {
				"Authorization": `Bearer ${token}`,
				"Content-Type": "application/json",
			},
		});
		const picture = await response.json();
		console.log(picture);
		return picture;
	} catch (error) {
		console.error("Erreur:", error);
		return null;
	}
}

async function displayPicture() {
	console.log("Display picture");
	const picture = await getPicture();
	const container = document.getElementById('photo');
	container.innerHTML = '';
	if (picture && picture.path > 0) {
		const img = document.createElement('img');
		img.src = picture.path;
		img.alt = picture.path;
		img.id = picture.id;
		img.style.cursor = 'pointer';
		container.appendChild(img);
	} else {
        const message = document.createElement('p');
        message.textContent = 'No image available';
        container.appendChild(message);
    }
}
