document.addEventListener('DOMContentLoaded', () => {
    displayGallery();
    const container = document.getElementById('galleryContainer');

    container.addEventListener('click', (e) => {
        const btn = e.target.closest('.btn-love');
        if (btn) {
            const photoId = btn.dataset.id;
			const counter = document.querySelector(
				`.like[data-id="${photoId}"]`
			);
			if (!counter) return;
			let count = parseInt(counter.textContent);
			if (btn.classList.contains('liked')) {
				btn.classList.remove('liked');
				count--;
			} else {
				btn.classList.add('liked');
				count++;
			}
			counter.textContent = `${count} likes`;
			sendLikes(photoId);
		}

    });
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

async function getPhotoUserData(pictureId) {
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
		return picture.Usr;
	} catch (error) {
		console.error("Erreur:", error);
		return null;
}
}

async function getPictures() {
	var page = window.location.pathname.split("/").pop();
	if (!page)
		window.location.href = window.location.pathname + "1";
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

async function sendLikes(photoId) {
	const user = await getUser();
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

async function getLikes(pId) {
    const token = localStorage.getItem('token');
    try {
        const response = await fetch(`/getLikes/${pId}`, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json",
            },
        });
        const likes = await response.json();

        const user = await getUser();
        const hasLiked = likes.some(like => like.uId === user.id);
        const heart = document.querySelector(
            `.btn-love[data-id="${pId}"]`
        );
        if (hasLiked) {
            heart.classList.add('liked');
        } else {
            heart.classList.remove('liked');
        }
        return likes;
    } catch (error) {
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
        for (const picture of pictures) {
            const user = await getPhotoUserData(picture.id);
            const username = user?.username ?? 'Username';
            container.innerHTML += `
                <div class="feed">
                    <section class="username">
                        <div class="id">
                            <p>${username}</p>
                        </div>
                    </section>
                    <section class="post" style="cursor: pointer;" onclick="window.location.href='/photo/${picture.id}'">
                        <img src="${picture.path}" alt="${picture.path}">
                    </section>
                    <section class="btn-group">
                        <button type="button" data-id="${picture.id}" id="likeBtn" class="btn-love"><i class="far fa-heart fa-lg"></i></button>
                        <button type="button" class="btn-comment"><i class="far fa-comment fa-lg"></i></button>
                    </section>
                    <section class="caption">
                        <p class="like" data-id="${picture.id}">${picture.likes ?? 0} likes</p>
                        <p><b class="id">${username}</b><span> ${picture.description ?? ''}</span></p>
                        <p class="time">${picture.time ?? ''}</p>
                    </section>
                </div>
            `;
		}
    } else {
        const message = document.createElement('p');
        message.textContent = 'No image available';
        container.appendChild(message);
    }
	pictures.forEach(p => getLikes(p.id));
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
