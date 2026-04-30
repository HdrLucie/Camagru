import { check_token } from './check-token.js';
import { getUser } from './get-user.js';

document.addEventListener('DOMContentLoaded', async () => {
	const r = check_token();
	if (r == false) {
		window.location.href = '/';
	}

	await displayComments();
	await getLikes();
});

async function checkToken() {
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/';
    }
}

async function getLikes() {
    const pId = window.location.pathname.split("/").pop();
    const token = localStorage.getItem('token');
    try {
        const response = await fetch(`/getLikes/${pId}`, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json",
            },
        });
		if (!response.ok) {
			throw new Error(`Response status: ${response.status}`);
		}
        const likes = await response.json();
        const user = await getUser();
        const hasLiked = likes.some(like => like.uId === user.id);
        const heart = document.getElementById('sendLikes');
        if (hasLiked) {
            heart.classList.add('is-toggled');
        } else {
            heart.classList.remove('is-toggled');
        }
        return likes;
    } catch (error) {
        return null;
    }
}

async function getComments() {
	const pId = window.location.pathname.split("/").pop();
	const token = localStorage.getItem('token');
    try {
        const response = await fetch(`/getComments/${pId}`, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json",
            },
        });
		if (!response.ok) {
			throw new Error(`Response status: ${response.status}`);
		}
        const comments = await response.json();
		return comments;
    } catch (error) {
        return null;
    }

}

async function displayComments() {
	comments = await getComments();
    if (!comments) return;
	listComments = document.getElementById("commentList");

	comments.forEach(comment=>{
		const li = document.createElement("li");
		li.classList.add("comment-item");
		const username = document.createElement("span");
		username.classList.add("username");
		username.textContent = comment.Username;
		username.className="user-data";
		const content = document.createElement("p");
		content.textContent = comment.Comment;
		li.appendChild(username);
		li.appendChild(content);
		listComments.appendChild(li);
	});
}

document.getElementById("sendLikes").onclick = async function () {
	const photoId = window.location.pathname.split("/").pop();
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

document.getElementById("deleteButton").onclick = async function () {
	const pId = window.location.pathname.split("/").pop();
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
}

const btn = document.querySelector('.js-button')
btn.addEventListener('click', (e) => {
	e.preventDefault()
	btn.classList.add('is-active')
	btn.classList.toggle('is-toggled')
	btn.blur()
	setTimeout(() => btn.classList.remove('is-active'), 400)
})

document.getElementById('sendLikes').addEventListener('click', function() {
	this.classList.toggle('liked');
});

const form = document.getElementById('com-form');

form.addEventListener('submit', async function (e) {
  e.preventDefault();
	var c = document.getElementById('comment');
	var u = await getUser();
	const token = localStorage.getItem('token');
	const pId = window.location.pathname.split("/").pop();
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
				"Comment": c.value,
			})
		});
		const data = await response.json();
		listComments = document.getElementById("commentList");
		const li = document.createElement("li");
		li.classList.add("comment-item");

		const username = document.createElement("span");
		username.classList.add("username");
		username.textContent = data.Username;

		const content = document.createElement("p");
		content.textContent = data.Comment;

		li.appendChild(username);
		li.appendChild(content);
		listComments.appendChild(li);
		c = "";
		form.reset();
	
	} catch (error) {
		console.error("Error: ", error);
	}
});

