import { check_token } from './check-token.js';

document.addEventListener('DOMContentLoaded', async () => {
	const r = await check_token();
	if (r == false) {
		window.location.href = '/';
	}
    displayStickers();
});

let selectedStickerId = null;
let floatingSticker = null;
let isDragging = false;

async function getStickers() {
    const token = localStorage.getItem('token');
    try {
        const response = await fetch("/getStickers", {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json",
            },
        });
        const stickers = await response.json();
        return stickers;
    } catch (error) {
        console.error("Erreur:", error);
        return null;
    }
}

async function getStickerNameById(stickerId) {
    const token = localStorage.getItem('token');
    try {
        const response = await fetch(`/getSticker/${stickerId}`, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Content-Type": "application/json",
            }
        });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const sticker = await response.json();
        return sticker.name;
    } catch (error) {
        console.error("Error while getting sticker:", error);
        return null;
    }
}

async function createStickerOnImage(stickerId, x, y) {
    const existingSticker = document.querySelector('.placed-sticker');
    const name = await getStickerNameById(stickerId);
    const stickerElement = document.createElement('img');
    stickerElement.src = "/stickers/" + name;
    stickerElement.className = 'placed-sticker';

    const tmp = document.getElementById('camera-section') || document.getElementById('file-upload-content');
    stickerElement.id = stickerId;
    stickerElement.style.left = x + 'px';
    stickerElement.style.top = y + 'px';
    stickerElement.dataset.relativeX = x;
    stickerElement.dataset.relativeY = y;
    tmp.appendChild(stickerElement);
}

function handleDrop(event) {
	event.preventDefault();
}

function handleDragOver(event) {
  event.preventDefault();
}

function handleDragStart(event) {
	event.preventDefault();
	isDragging = true;
	selectedStickerId = event.target.id;
	createFloatingSticker(event.target.src, event.clientX, event.clientY);
	document.body.style.cursor = 'grabbing';
}

function removeFloatingSticker() {
	if (floatingSticker) {
		floatingSticker.remove();
		floatingSticker = null;
	}
}

function createFloatingSticker(src, x, y) {
	removeFloatingSticker();

	floatingSticker = document.createElement('img');
	floatingSticker.src = src;
	floatingSticker.className = "floating-sticker";
	floatingSticker.style.cssText = `
    position: fixed;
    pointer-events: none;
    z-index: 9999;
    opacity: 0.8;
	max-width: 128px;
	max-height: 128px;
    // transform: translate(-50%, -50%)`;
	floatingSticker.style.left = x + 'px';
	floatingSticker.style.top = y + 'px';
	document.body.appendChild(floatingSticker);
}

document.addEventListener('mousemove', (event) => {
	if (floatingSticker) {
			floatingSticker.style.left = event.clientX + 'px';
			floatingSticker.style.top = event.clientY + 'px';
	}
});

function cancelSelection() {
	selectedStickerId = null;
	isDragging = false;
	removeFloatingSticker();
	document.body.style = '';
}


document.addEventListener('mouseup', (event) => {
	if (!isDragging || !selectedStickerId) return;

	const dropZone = document.getElementById('camera-section') || document.getElementById('file-upload-content');
	if (dropZone && dropZone.contains(event.target)) {
		const rect = dropZone.getBoundingClientRect();
		const x = event.clientX - rect.left;
		const y = event.clientY - rect.top;
		const id = selectedStickerId;
		createStickerOnImage(id, x, y);
	}
	cancelSelection();
});

document.addEventListener('click', (event) => {
	if (selectedStickerId === null) return;

	const dropZone = document.getElementById('camera-section') || document.getElementById('file-upload-content');
	if (dropZone && dropZone.contains(event.target)) {
		const rect = dropZone.getBoundingClientRect();
		const x = event.clientX - rect.left;
		const y = event.clientY - rect.top;
		const id = selectedStickerId;
		createStickerOnImage(id, x, y);
	}
	cancelSelection();
});

async function displayStickers() {
    const stickers = await getStickers();
    const container = document.getElementById('stickersContainer');
    container.innerHTML = '';
    if (stickers && stickers.length > 0) {
        stickers.forEach(sticker => {
            const img = document.createElement('img');
            img.src = "/stickers/" + sticker.name;
			img.alt = sticker.name;
			img.id = sticker.id;
			img.draggable = true;
			img.className = 'sticker-item';
			img.style.cursor = 'grab';
			img.addEventListener('dragstart', handleDragStart);
			img.addEventListener('click', (event) => {
				event.stopPropagation();
				if (selectedStickerId === sticker.id) return;
				selectedStickerId = sticker.id;
				createFloatingSticker(img.src, event.clientX, event.clientY);
				document.body.style.cursor = 'grabbing';
			});
            container.appendChild(img);
        });
    } else {
        const message = document.createElement('p');
        message.textContent = 'No sticker available';
        container.appendChild(message);
    }
}
