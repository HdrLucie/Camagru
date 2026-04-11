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
		console.log(pictures);
        return pictures;
    } catch (error) {
        console.error("Erreur:", error);
        return null;
    }
}

async function displayGallery() {
	const data = await getPictures();
    const container = document.getElementById('galleryContainer');
    container.innerHTML = '';

    const pictures = data?.pictures;
    const isLast = data?.last;
    if (pictures && pictures.length > 0) {
        pictures.forEach(picture => {
            const div = document.createElement('div');
            div.className = 'gallery-item';
            const img = document.createElement('img');
            img.className = 'gallery-image';
            img.src = picture.path;
            img.alt = picture.path;
            img.id = picture.id;
            img.style.cursor = 'pointer';
            div.addEventListener('click', function() {
                window.location.href = `/photo/${picture.id}`;
            });

            const galleryItemInfo = document.createElement('div');
            galleryItemInfo.className = 'gallery-item-info';

            const ul = document.createElement('ul');

            const liLikes = document.createElement('li');
            liLikes.className = 'gallery-item-likes';

            const spanLikes = document.createElement('span');
            spanLikes.className = 'visually-hidden';
            spanLikes.textContent = 'Likes:';

            const iHeart = document.createElement('i');
            iHeart.className = 'fas fa-heart';
            iHeart.setAttribute('aria-hidden', 'true');

            liLikes.appendChild(spanLikes);
            liLikes.appendChild(iHeart);
            liLikes.append(` ${picture.likes ?? 0}`);

            const liComments = document.createElement('li');
            liComments.className = 'gallery-item-comments';

            const spanComments = document.createElement('span');
            spanComments.className = 'visually-hidden';
            spanComments.textContent = 'Comments:';

            const iComment = document.createElement('i');
            iComment.className = 'fas fa-comment';
            iComment.setAttribute('aria-hidden', 'true');

            liComments.appendChild(spanComments);
            liComments.appendChild(iComment);
            liComments.append(` ${picture.comments ?? 0}`);

            ul.appendChild(liLikes);
            ul.appendChild(liComments);
            galleryItemInfo.appendChild(ul);
            div.appendChild(img);
            div.appendChild(galleryItemInfo);
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

async function nextPage() {
	const data = await getPictures();
	const isLast = data?.last;
	const next = document.getElementById('nextBtn');
	const page = parseInt(window.location.pathname.split("/").pop());
	if (isLast == true)
		return ;
	window.location.href = `/gallery/${page + 1}`;
}
