import { check_token } from './check-token.js';
import { getUser } from './get-user.js';

document.addEventListener('DOMContentLoaded', async () => {
	const r = await check_token();
	if (r === false) {
		window.location.href = '/';
	}
});

(() => {
	let width = 680;
	let height = 550;
	let streaming = false;
	let video = null;
	let canvas = null;
	let startButton = null;
	let sendButton = null;
	let clearButton = null;

	function cameraSize() {
		if (window.innerWidth <= 760) {
			width = 430;
			height = 350;
		}
	}

	function startup() {
		cameraSize();

		video = document.getElementById("video");
		canvas = document.getElementById("canvas");
		startButton = document.getElementById("startbutton");
		sendButton = document.getElementById("sendButton");
		clearButton = document.getElementById("clearButton");

		navigator.mediaDevices
			.getUserMedia({ video: true, audio: false })
			.then((stream) => {
				video.srcObject = stream;
				video.play();
			})
			.catch((err) => {
				console.error(`An error occurred: ${err}`);
			});

		video.addEventListener("canplay", () => {
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
		}, false);

		startButton.addEventListener("click", (ev) => {
			takePicture();
			ev.preventDefault();
		}, false);

		sendButton.addEventListener("click", (ev) => {
			sendPicture();
			ev.preventDefault();
		}, false);

		clearButton.addEventListener("click", (ev) => {
			clearPhoto();
			ev.preventDefault();
		}, false);
	}

	function clearPhoto() {
		const context = canvas.getContext("2d");
		context.fillStyle = "#aaaaaa";
		context.fillRect(0, 0, canvas.width, canvas.height);

		document.getElementById("photo").style.display = "none";
		document.getElementById("output").style.display = "none";
		document.getElementById("camera").style.display = "inline-block";
		document.querySelectorAll(".placed-sticker").forEach((sticker) => sticker.remove());
	}

	function takePicture() {
		if (!width || !height) return;

		const context = canvas.getContext("2d");
		canvas.width = width;
		canvas.height = height;
		context.drawImage(video, 0, 0, width, height);

		const photo = document.getElementById("photo");
		photo.src = canvas.toDataURL("image/png"); // affichage seulement
		photo.style.display = "inline-block";
		document.getElementById("output").style.display = "inline-block";
		document.getElementById("camera").style.display = "none";
	}

	async function sendPicture() {
		const token = localStorage.getItem('token');
		const elements = document.querySelectorAll('.placed-sticker');
		if (elements.length === 0) {
			alert('Veuillez placer un sticker sur votre photo avant d\'envoyer !');
			return;
		}

		sendButton.disabled = true;

		try {
			const user = await getUser();

			const imgBlob = await new Promise((resolve) => {
				canvas.toBlob(resolve, 'image/jpeg', 0.8);
			});

			const stickers = [...elements].map((sticker) => ({
				path: sticker.src,
				posX: Math.floor(sticker.dataset.relativeX),
				posY: Math.floor(sticker.dataset.relativeY),
				id: parseInt(sticker.id),
			}));

			const formData = new FormData();
			formData.append('image', imgBlob, 'photo.jpg');
			formData.append('id', user.id); // TODO : à supprimer, le back doit lire l'id dans le JWT
			formData.append('stickers', JSON.stringify(stickers));
			formData.append('timestamp', new Date().toISOString());

			const response = await fetch("/sendImage", {
				method: "POST",
				headers: {
					"Authorization": `Bearer ${token}`,
				},
				body: formData,
			});

			if (!response.ok) {
				throw new Error(`Response status: ${response.status}`);
			}
			window.location.href = "/gallery/1";
		} catch (error) {
			console.error("Erreur:", error);
			alert("L'envoi a échoué, réessaie.");
			sendButton.disabled = false;
		}
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

