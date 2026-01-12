document.addEventListener('DOMContentLoaded', () => {
	checkToken();
	displayComments();
});

async function checkToken() {
    const token = localStorage.getItem('token');
    console.log("Function check token")
    console.log(token)
    if (!token) {
        window.location.href = '/';
    }
}

async function getComments() {
	const token = localStorage.getItem('token');
    try {
        const response = await fetch("/getComments", {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json",
            },
        });
        const comments = await response.json();
		return comments;
    } catch (error) {
        console.error("Erreur:", error);
        return null;
    }

}

async function displayComments() {
	comments = await getComments();
	listComments = document.getElementById("commentList");

	comments.forEach(comment=>{

		const li = document.createElement("li");
		li.classList.add("comment-item");

		const username = document.createElement("span");
		username.classList.add("username");
		username.textContent = comment.Username;

		const content = document.createElement("p");
		content.textContent = comment.Comment;

		li.appendChild(username);
		li.appendChild(content);
		listComments.appendChild(li);
	});
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
		console.log(userData);
		return userData;
    } catch (error) {
        console.error("Erreur:", error);
        return null;
    }
}

document.getElementById("sendLikes").onclick = async function () {
	console.log("Likes button");
	const photoId = window.location.pathname.split("/").pop();
	console.log(photoId);
	const user = await getUser();
	console.log(user);
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

document.getElementById('sendLikes').addEventListener('click', function() {
    this.classList.toggle('fa-regular');
    this.classList.toggle('fa-solid');
});

const form = document.getElementById('com-form');

form.addEventListener('submit', async function (e) {
  e.preventDefault();
	console.log('\n\nSend Comments\n\n');
	var c = document.getElementById('comment');
	var u = await getUser();
	const token = localStorage.getItem('token');
	const pId = window.location.pathname.split("/").pop();

	try {
		const response = await fetch("/sendComments", {
			method: "POST",
			headers: {
				"Authorization": `Bearer ${token}`,
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

