document.addEventListener('DOMContentLoaded', () => {
    displayStickers();
});

const draggableItem = document.getElementById("draggableItem");

async function getStickers() {
    const token = localStorage.getItem('token');
    console.log("getStickers function")    
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
    if (existingSticker) {
        existingSticker.remove();
    }
    if (document.readyState === 'loading') {
        await new Promise(resolve => {
            document.addEventListener('DOMContentLoaded', resolve);
        });
    }

	const name = await getStickerNameById(stickerId);
	console.log(name);
    const stickerElement = document.createElement('img');
    stickerElement.src = "/stickers/" + name;
    stickerElement.className = 'placed-sticker';
    stickerElement.style.position = 'absolute';
    stickerElement.style.left = x + 'px';
    stickerElement.style.top = y + 'px';
    stickerElement.style.width = 'auto';
    stickerElement.style.height = 'auto';
	stickerElement.style.zIndex = 1000;
    stickerElement.style.pointerEvents = 'none';
    const dropZone = document.getElementById('camera');
	console.log(dropZone);
    dropZone.appendChild(stickerElement);
}

function handleDrop(event) {
	event.preventDefault();
	console.log("handle drop");
	const id = event
	.dataTransfer
	.getData('text/plain');
	console.log(id);
	const dropX = event.offsetX;
    const dropY = event.offsetY;
    console.log("Drop position:", dropX, dropY);
    createStickerOnImage(id, dropX, dropY);
}

function handleDragOver(event) {
  event.preventDefault();
}

function handleDragStart(event) {
	console.log("Handle drag start");
	event
    .dataTransfer
    .setData('text/plain', event.target.id);
	console.log(event.target.id);
}

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
			img.addEventListener('dragstart', handleDragStart);
            container.appendChild(img);
        });
    } else {
        const message = document.createElement('p');
        message.textContent = 'No sticker available';
        container.appendChild(message);
    }
}
