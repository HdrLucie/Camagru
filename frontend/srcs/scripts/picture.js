import { check_token } from './check-token.js';

document.addEventListener('DOMContentLoaded', async () => {
	const r = await check_token();
	if (r == false) {
		window.location.href = '/';
	}
	displayPicture();
});

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
		return userData;
    } catch (error) {
        console.error("Erreur:", error);
        return null;
    }
}

async function getPicture(pictureId) {
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
		return picture;
	} catch (error) {
		console.error("Erreur:", error);
		windo.location.href="/gallery/1"
		return null;
	}
}

async function displayPicture() {
	const pictureId = window.location.pathname.split("/").pop();
	const data = await getPicture(pictureId);
	const picture = data?.Picture;
	const user = data?.Usr;

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
