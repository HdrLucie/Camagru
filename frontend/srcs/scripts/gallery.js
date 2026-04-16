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

// async function getComments(pId) {
// 	const pId = window.location.pathname.split("/").pop();
// 	const token = localStorage.getItem('token');
//     try {
//         const response = await fetch(`/getComments/${pId}`, {
//             method: "GET",
//             headers: {
//                 "Authorization": `Bearer ${token}`,
//                 "Content-Type": "application/json",
//             },
//         });
//         const comments = await response.json();
// 		return comments;
//     } catch (error) {
//         return null;
//     }
//
// }
//
// async function displayComments() {
// 	comments = await getComments();
//     if (!comments) return;
// 	listComments = document.getElementById("commentList");
//
// 	comments.forEach(comment=>{
//
// 		const li = document.createElement("li");
// 		li.classList.add("comment-item");
//
// 		const username = document.createElement("span");
// 		username.classList.add("username");
// 		username.textContent = comment.Username;
// 		username.className="user-data";
//
// 		const content = document.createElement("p");
// 		content.textContent = comment.Comment;
//
// 		li.appendChild(username);
// 		li.appendChild(content);
// 		listComments.appendChild(li);
// 	});
// }



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
                        <button type="button" class="btn-love"><i class="far fa-heart fa-lg"></i></button>
                        <button type="button" class="btn-comment"><i class="far fa-comment fa-lg"></i></button>
                    </section>
                    <section class="caption">
                        <p class="like">${picture.likes ?? 0} likes</p>
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
