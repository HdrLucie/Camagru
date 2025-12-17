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
            img.src = picture.path;
			img.alt = picture.path;
			img.id = picture.id;
			img.style.cursor = 'pointer';
            img.addEventListener('click', function() {
                window.location.href = `/photo`;
            });
            container.appendChild(img);
        });
    } else {
        const message = document.createElement('p');
        message.textContent = 'No image available';
        container.appendChild(message);
    }
}

