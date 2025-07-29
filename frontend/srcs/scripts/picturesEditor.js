document.addEventListener('DOMContentLoaded', () => {
	checkToken();
	//loadUserData();
});

async function checkToken() {
	const token = localStorage.getItem('token');
	console.log("Function check token")
	console.log(token)
	if (!token) {
		// alert('No token found. Please login.');
		window.location.href = '/';
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

(() => {
	// The width and height of the captured photo. We will set the
	// width to the value defined here, but the height will be
	// calculated based on the aspect ratio of the input stream.

	const width = 100; // We will scale the photo width to this
	let height = 100; // This will be computed based on the input stream

	// |streaming| indicates whether or not we're currently streaming
	// video from the camera. Obviously, we start at false.

	let streaming = false;

	// The various HTML elements we need to configure or control. These
	// will be set by the startup() function.

	let video = null;
	let canvas = null;
	let photo = null;
	let startbutton = null;
	let sendButton = null;

	function showViewLiveResultButton() {
		if (window.self !== window.top) {
			// Ensure that if our document is in a frame, we get the user
			// to first open it in its own tab or window. Otherwise, it
			// won't be able to request permission for camera access.
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
		photo = document.getElementById("photo");
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

					// Firefox currently has a bug where the height can't be read from
					// the video, so we will make assumptions if this happens.

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
		clearphoto();
	}

		function clearphoto() {
			const context = canvas.getContext("2d");
			context.fillStyle = "#AAA";
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
				photo.setAttribute("src", data);
			} else {
				clearphoto();
			}
		}

	async function sendPictures() {

		console.log("\n\nsendButton function.\n\n");
		var path;
		const token = localStorage.getItem('token');
		const sticker = document.getElementsByClassName('placed-sticker');
		const str = sticker[0];
		console.log(sticker[0], sticker);
		const user = await getUser();
		const imgBlob = await new Promise(resolve => {
            canvas.toBlob(resolve, 'image/jpeg', 0.8); // Compression JPEG à 80%
        });
		const blobUrl = URL.createObjectURL(imgBlob);
		console.log("Votre photo capturée:", blobUrl);
		const formData = new FormData();
        formData.append('image', imgBlob, 'photo.jpg');
		formData.append('id', user.id);
		formData.append('imageId', sticker[0].id);
		const relativeX = sticker[0].relativeX;
		const relativeY = sticker[0].relativeY;

		formData.append('posX', JSON.stringify(relativeX));
		formData.append('posY', JSON.stringify(relativeY));
        formData.append('timestamp', new Date().toISOString());
		console.log(formData);
		const response = await fetch("/sendImage", {
			method: "POST",
			headers: {
				"Authorization": `Bearer ${token}`
			},
			body: formData
		});
		console.log(await response.json());

	}
		window.addEventListener("load", startup, false);
})();


