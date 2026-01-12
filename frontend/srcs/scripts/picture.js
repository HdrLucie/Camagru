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

async function getPicture(pictureId) {
	console.log("Id" + pictureId)
	const token = localStorage.getItem('token');
	try {
		const response = await fetch(`/getPicture/${pictureId}`, {
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
	const pictureId = window.location.pathname.split("/").pop();
	const picture = await getPicture(pictureId);
	console.log(picture.path);
	const container = document.getElementById('photo');
	container.innerHTML = '';
	if (picture.path != '') {
		const img = document.createElement('img');
		img.src = "/" + picture.path;
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
