import {displayJobs} from "./jobs.js";
import {displayUsers} from "./users.js";

export const pageBody = document.querySelector('[data-role=page-body]');

export function resetContent(login = false) {
    pageBody.innerHTML = "";
    if (login) {
        document.querySelector('[data-menu=login]').classList.remove('jinya-hidden');
        document.querySelector('[data-menu=app]').classList.add('jinya-hidden');
    } else {
        document.querySelector('[data-menu=login]').classList.add('jinya-hidden');
        document.querySelector('[data-menu=app]').classList.remove('jinya-hidden');
    }
}

export async function request(url, method, body = {}) {
    const data = {
        method,
    };
    if (method !== 'GET' && method !== 'HEAD') {
        data.body = JSON.stringify(body);
    }

    return await fetch(url, data);
}

export function displayModal(modalHtml) {
    const container = document.createElement('div');
    container.innerHTML = modalHtml;
    container.classList.add('cosmo-modal__container');

    const backdrop = document.createElement('div');
    backdrop.classList.add('cosmo-modal__backdrop');

    document.body.appendChild(backdrop);
    document.body.appendChild(container);

    function closeModal() {
        backdrop.remove();
        container.remove();
    }

    return closeModal;
}

export function confirm(title, message, okButton = 'Ok', cancelButton = 'Cancel') {
    const modalHtml = `
<div class="cosmo-modal">
    <h1 class="cosmo-modal__title">${title}</h1>
    <p class="cosmo-modal__content">${message}</p>
    <div class="cosmo-modal__button-bar">
        <button class="cosmo-button" data-action="modal-ok">${okButton}</button>
        <button class="cosmo-button" data-action="modal-cancel">${cancelButton}</button>
    </div>
</div>`;

    const closeModal = displayModal(modalHtml);

    return new Promise((resolve) => {
        document.querySelector('[data-action="modal-ok"]').addEventListener('click', () => {
            resolve(true);
            closeModal();
        });
        document.querySelector('[data-action="modal-cancel"]').addEventListener('click', () => {
            resolve(false);
            closeModal();
        });
    });
}

export function resetListClass() {
    document.querySelectorAll('.cosmo-list__item--active').forEach(item => item.classList.remove('cosmo-list__item--active'));
}

export function resetTabClass() {
    document.querySelectorAll('.cosmo-tab-control__tab-link--active').forEach(item => item.classList.remove('cosmo-tab-control__tab-link--active'));
}

export function resetMenuClass() {
    document.querySelectorAll('.cosmo-menu-bar__main-item--active').forEach(item => item.classList.remove('cosmo-menu-bar__main-item--active'));
}

export function initMenuNavigation() {
    document.querySelector('[data-menu-item=backups]').addEventListener('click', async (e) => {
        e.preventDefault();
        resetMenuClass();
        e.target.classList.add('cosmo-menu-bar__main-item--active');
        await displayJobs();
        location.hash = '#backups';
    });
    document.querySelector('[data-menu-item=users]').addEventListener('click', async (e) => {
        e.preventDefault();
        resetMenuClass();
        e.target.classList.add('cosmo-menu-bar__main-item--active');
        await displayUsers();
        location.hash = '#users';
    });
    if (window.location.hash === '#users') {
        resetMenuClass();
        document.querySelector('[data-menu-item=users]').classList.add('cosmo-menu-bar__main-item--active');
    }
}