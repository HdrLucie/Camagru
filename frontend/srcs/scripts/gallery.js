import { check_token } from './check-token.js';

document.addEventListener('DOMContentLoaded', async () => {
	const r = await check_token();
	if (r == true) {
		displayGallery();
	} else {
		displaySimpleGallery();
	}
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
		window.location.href = "/";
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

async function getComments(pId) {
    const token = localStorage.getItem('token');
    try {
        const response = await fetch(`/getComments/${pId}`, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json",
            },
        });
        return await response.json();
    } catch (error) {
        console.error("Erreur:", error);
        return [];
    }
}

document.getElementById('galleryContainer').addEventListener('click', async function(e) {
	const btn = e.target.closest('.btn-delete');
	if (!btn)
		return;
	const pId = btn.dataset.id;
	const user = await getUser();
	const token = localStorage.getItem('token');
	const response = await fetch("/deleteImg", {
		method: "POST",
		headers: {
			"Authorization": `Bearer ${token}`,
			"Content-type": "application/json"
		},
		body: JSON.stringify({
			"Username": user.username,
			"uId": user.id,
			"pId": Number(pId),
		})
	});
	if (response.ok) {
		window.location.href = "/gallery/1";
	} else if (response.status === 403) {
		const data = await response.json();
		alert(data.error);
	}
});

document.getElementById('galleryContainer').addEventListener('submit', async function(e) {
    e.preventDefault();
    const form = e.target.closest('.com-form');
    if (!form) return;

    const pId = form.dataset.id;
    const input = form.querySelector('.comment-input-field');
    const u = await getUser();
    const token = localStorage.getItem('token');

    try {
        const response = await fetch("/sendComments", {
            method: "POST",
            headers: {
                "authorization": `Bearer ${token}`,
                "Content-type": "application/json"
            },
            body: JSON.stringify({
                "Username": u.username,
                "Id": u.id,
                "Photo": Number(pId),
                "Comment": input.value,
            })
        });
        const data = await response.json();
		const commentList = document.querySelector(`.commentList[data-id="${pId}"]`);
		commentList.innerHTML += `
			<p><b class="id">${data.Username}</b><span> ${data.Comment}</span></p>
`;
          form.reset();

    } catch (error) {
        console.error("Error: ", error);
    }
});

async function displaySimpleGallery() {
	const data = await getPictures();
    const container = document.getElementById('galleryContainer');
    container.innerHTML = '';
    const pictures = data?.pictures;
    const isLast = data?.last;
	if (pictures && pictures.length > 0) {
		for (const picture of pictures) {
			container.innerHTML += `
				<div class="feed">
				<section class="username">
				<div class="id">
				</div>
				</section>
				<section class="post">
				<img src="${picture.path}" alt="${picture.path}">
				</section>
				<section class="btn-group">
			
				</section>
				<section class="caption">
				<p class="like" data-id="${picture.id}">${picture.likes ?? 0} likes</p>
				<form class="com-form" data-id="${picture.id}">
				<div class="commentList" data-id="${picture.id}"></div>
				</form>
				
				</section>
				</div>
				
				`;
		}
    } else {
        const message = document.createElement('p');
        message.textContent = 'No image available';
        container.appendChild(message);
    }
}

async function displayGallery() {
	const data = await getPictures();
    const container = document.getElementById('galleryContainer');
    container.innerHTML = '';
	const currentUser = await getUser();
    const pictures = data?.pictures;
    const isLast = data?.last;
	if (pictures && pictures.length > 0) {
		for (const picture of pictures) {
			const [user, comments] = await Promise.all([
				getPhotoUserData(picture.id),
				getComments(picture.id)
			]);
			const commentsHTML = (comments ?? []).map(c => `
				<p><b class="id">${c.Username}</b><span> ${c.Comment}</span></p>`).join('');
			const username = user?.username ?? 'Username';
			container.innerHTML += `
				<div class="feed">
				<section class="username">
				<div class="id">
				<p><b>${username}</b></p>
				</div>
				</section>
				<section class="post">
				<img src="${picture.path}" alt="${picture.path}">
				</section>
				<section class="btn-group">
				<button type="button" data-id="${picture.id}" id="likeBtn" class="btn-love"><i class="far fa-heart fa-lg"></i></button>
				<button type="button" data-id="${picture.id}" style="display:none;" class="btn-delete"><i class="fa fa-trash fa-lg"></i></button>
			
				</section>
				<section class="caption">
				<p class="like" data-id="${picture.id}">${picture.likes ?? 0} likes</p>
				<p><span></span></p>${commentsHTML ?? ''}
				<form class="com-form" data-id="${picture.id}">
				<div class="commentList" data-id="${picture.id}"></div>
				<div class="comment-input">
				<input type="text" class="comment-input-field" placeholder="Say something..." required />
				<button type="submit">Envoyer</button>
				</div>
				</form>
				
				</section>
				</div>
				
				`;
			    if (currentUser.id === picture.userId) {
					const btn = document.querySelector(`.btn-delete[data-id="${picture.id}"]`);
					btn.style.display = 'inline';
				}
		}
    } else {
        const message = document.createElement('p');
        message.textContent = 'No image available';
        container.appendChild(message);
    }
	if (pictures && pictures.length > 0) 
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
