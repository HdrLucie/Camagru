document.addEventListener('DOMContentLoaded', () => {
    displayGallery();
});

document.getElementById("burger").onclick = function () {
    let burger = document.querySelector(".js-burger");
    let nav = document.querySelector(".js-nav");

    nav.classList.toggle("_active");
    burger.classList.toggle("_active");
}

function redirectionPage(path) {
    window.location.href = path;
}

async function getPictures() {
	const page = window.location.pathname.split("/").pop();
    const token = localStorage.getItem('token');
    try {
        const response = await fetch(`/getPictures/${page}`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
        });
        const pictures = await response.json();
        return pictures;
    } catch (error) {
        console.error("Erreur:", error);
        return null;
    }
}

async function displayGallery() {
    const pictures = await getPictures();
    const container = document.getElementById('galleryContainer');
    container.innerHTML = '';
    if (pictures && pictures.length > 0) {
        pictures.forEach(picture => {
			const div = document.createElement('div');
			div.className = 'gallery-item';
            const img = document.createElement('img');
			img.className = 'gallery-item'
            img.src = picture.path;
			img.alt = picture.path;
			img.id = picture.id;
			img.style.cursor = 'pointer';
            img.addEventListener('click', function() {
                window.location.href = `/photo/${picture.id}`;
            });
			div.appendChild(img);
            container.appendChild(div);
        });
    } else {
        const message = document.createElement('p');
        message.textContent = 'No image available';
        container.appendChild(message);
    }
}

function prevPage() {
	const next = document.getElementById('prevBtn');
	const page = parseInt(window.location.pathname.split("/").pop());
	
	if (page == 1) return;
	window.location.href = `/gallery/${page - 1}`;
}

function nextPage() {
	const next = document.getElementById('nextBtn');
	const page = parseInt(window.location.pathname.split("/").pop());
	
	window.location.href = `/gallery/${page + 1}`;
}
