document.addEventListener('DOMContentLoaded', () => {
    displayGallery();
});

async function getPictures() {
	console.log("Get all pictures");
    const token = localStorage.getItem('token');
    try {
        const response = await fetch("/getPictures", {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json",
            },
        });
        const pictures = await response.json();
		console.log(pictures);
        return pictures;
    } catch (error) {
        console.error("Erreur:", error);
        return null;
    }
}

async function displayGallery() {
	console.log("Display gallery");
    const pictures = await getPictures();
    const container = document.getElementById('galleryContainer');
    container.innerHTML = '';
    if (pictures && pictures.length > 0) {
        pictures.forEach(picture => {
            const img = document.createElement('img');
            img.src = "/backend/srcs/" + picture.path;
			img.alt = "/backend/srcs/" + picture.path;
			img.id = picture.id;
			console.log(img);
            container.appendChild(img);
        });
    } else {
        const message = document.createElement('p');
        message.textContent = 'No sticker available';
        container.appendChild(message);
    }
}


