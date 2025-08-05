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

	const width = 680; // We will scale the photo width to this
	let height = 550; // This will be computed based on the input stream
	let streaming = false;
	let video = null;
	let canvas = null;
	let photo = null;
	let startbutton = null;
	let sendButton = null;

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
		const user = await getUser();
		const imgBlob = await new Promise(resolve => {
            canvas.toBlob(resolve, 'image/jpeg', 0.8); // Compression JPEG à 80%
        });
		const blobUrl = URL.createObjectURL(imgBlob);
		const formData = new FormData();
		console.log(sticker[0]);
        formData.append('image', imgBlob, 'photo.jpg');
		formData.append('id', user.id);
		formData.append('imageId', sticker[0].id);
		const relativeX = sticker[0].dataset.relativeX;
		const relativeY = sticker[0].dataset.relativeY;
		formData.append('stickerPath', sticker[0].src);
		formData.append('posX', JSON.stringify(Math.floor(relativeX)));
		formData.append('posY', JSON.stringify(Math.floor(relativeY)));
        formData.append('timestamp', new Date().toISOString());
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
