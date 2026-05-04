import { check_token } from './check-token.js';
import { getUser } from './get-user.js';

document.addEventListener('DOMContentLoaded', async () => {
	const r = await check_token();
	if (r === false) {
		window.location.href = '/';
	}
});

async function getPictures() {
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
			const img = document.createElement('img');
			img.src = "/pictures/" + sticker.Path;
			img.alt = sticker.Path;
			img.id = sticker.id;
			container.appendChild(img);
		});
	} else {
		const message = document.createElement('p');
		message.textContent = 'No sticker available';
		container.appendChild(message);
	}
}

(() => {

	let width = 680;
	let height = 550;
	width, height = cameraSize();
	let streaming = false;
	let video = null;
	let canvas = null;
	let startbutton = null;
	let sendButton = null;
	
	function cameraSize() {
		console.log(window.innerWidth);
		if (window.innerWidth <= 760) {
			width = 430;
			height = 350;
			return width, height
		}
		return width, height
	}

	function showViewLiveResultButton() {
		if (window.self !== window.top) {
			document.querySelector(".contentarea").remove();
			const button = document.createElement("button");
			button.textContent = "View live result of the example code above";
			document.body.append(button);
			button.addEventListener("click", () => window.open(location.href));
			return true;
		}
		return false;
	}

	function startup() {
		if (showViewLiveResultButton()) {
			return;
		}
		video = document.getElementById("video");
		canvas = document.getElementById("canvas");
		startbutton = document.getElementById("startbutton");
		sendButton = document.getElementById("sendButton");

		navigator.mediaDevices
			.getUserMedia({ video: true, audio: false })
			.then((stream) => {
				video.srcObject = stream;
				video.play();
			})
			.catch((err) => {
				console.error(`An error occurred: ${err}`);
			});

		video.addEventListener(
			"canplay",
			(ev) => {
				if (!streaming) {
					height = video.videoHeight / (video.videoWidth / width);
					if (isNaN(height)) {
						height = width / (4 / 3);
					}
					video.setAttribute("width", width);
					video.setAttribute("height", height);
					canvas.setAttribute("width", width);
					canvas.setAttribute("height", height);
					streaming = true;
				}
			},
			false,
		);

		startbutton.addEventListener(
			"click",
			(ev) => {
				takepicture();
				ev.preventDefault();
			},
			false,
		);

		sendButton.addEventListener(
			"click", 
			(ev) => {
				sendPictures();
				ev.preventDefault();
			},
			false
		);
	}

	function clearPhoto() {
		const context = canvas.getContext("2d");
		context.fillStyle = "#aaaaaa";
		context.fillRect(0, 0, canvas.width, canvas.height);

		const data = canvas.toDataURL("image/png");
		photo.setAttribute("src", data);
	}

	function takepicture() {
		const context = canvas.getContext("2d");
		if (width && height) {
			canvas.width = width;
			canvas.height = height;
			context.drawImage(video, 0, 0, width, height);

			const data = canvas.toDataURL("image/png");
			const photo = document.getElementById("photo");
			photo.style.display = "inline-block"
			const output = document.getElementById('output');
			output.style.display = 'inline-block';
			const cam = document.getElementById('camera');
			cam.style.display='none';
			photo.src = data;
		}
	}

	async function sendPictures() {
		var path;
		const token = localStorage.getItem('token');
		const stickers = [];
		const elements = document.querySelectorAll('.placed-sticker');
		if (!elements || elements.length === 0) {
			alert('Veuillez placer un sticker sur votre photo avant d\'envoyer !');
			return;
		}
		const user = await getUser();
		const imgBlob = await new Promise(resolve => {
			canvas.toBlob(resolve, 'image/jpeg', 0.8);
		});
		const blobUrl = URL.createObjectURL(imgBlob);
		const formData = new FormData();
		formData.append('image', imgBlob, 'photo.jpg');
		formData.append('id', user.id);
		elements.forEach(sticker => {
			stickers.push({
				path: sticker.src,
				posX: Math.floor(sticker.dataset.relativeX),
				posY: Math.floor(sticker.dataset.relativeY),
				id: parseInt(sticker.id) 
			});
		});
		formData.append('stickers', JSON.stringify(stickers));
		formData.append('timestamp', new Date().toISOString());
		const response = await fetch("/sendImage", {
			method: "POST",
			headers: {
				"Authorization": `Bearer ${token}`
			},
			body: formData
		});
		window.location.href = `/gallery/1`;

	}

	window.addEventListener("load", startup, false);
})();

	function clearPhoto() {
		const context = canvas.getContext("2d");
		context.fillStyle = "#aaaaaa";
		context.fillRect(0, 0, canvas.width, canvas.height);

		const photo = document.getElementById("photo");
		photo.style.display = "none"
		const output = document.getElementById('output');
		output.style.display = 'none';
		const cam = document.getElementById('camera');
		cam.style.display='inline-block';

	}

