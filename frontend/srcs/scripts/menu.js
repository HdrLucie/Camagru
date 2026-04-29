import { check_token } from './check-token.js';

// document.addEventListener('DOMContentLoaded', async () => {
//     const isConnected = await check_token();
//
//     const menuConnected = document.querySelector('.menu-connected');
//     const menuGuest = document.querySelector('.menu-guest');
//
//     if (isConnected) {
//         menuConnected.style.display = 'flex';
//         menuGuest.style.display = 'none';
//     } else {
//         menuConnected.style.display = 'none';
//         menuGuest.style.display = 'flex';
//     }
// });
//
// document.getElementById("burger").onclick = function () {
//     let burger = document.querySelector(".js-burger");
//     let nav = document.querySelector(".js-nav");
//
//     nav.classList.toggle("_active");
//     burger.classList.toggle("_active");
// }
//
// function redirectionPage(path) {
//     window.location.href = path;
// }

document.addEventListener('DOMContentLoaded', async () => {

    document.querySelectorAll('[data-path]').forEach(btn => {
        btn.addEventListener('click', () => {
            window.location.href = btn.dataset.path;
        });
    });

    document.getElementById("burger").addEventListener('click', () => {
        document.querySelector(".js-burger").classList.toggle("_active");
        document.querySelector(".js-nav").classList.toggle("_active");
    });

    const isConnected = await check_token();
    const menuConnected = document.querySelector('.menu-connected');
    const menuGuest = document.querySelector('.menu-guest');

    if (isConnected) {
        menuConnected.style.display = 'flex';
        menuGuest.style.display = 'none';
    } else {
        menuConnected.style.display = 'none';
        menuGuest.style.display = 'flex';
    }
});
