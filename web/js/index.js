import {initMenuNavigation, pageBody, request} from "./utils.js";
import {displayJobs} from "./jobs.js";
import {displayLogin} from "./login.js";
import {displayUsers} from "./users.js";
import {displayImportExport} from "./importExport.js";

function initGlobalActions() {
    document.querySelector('[data-action=logout]').addEventListener('click', async () => {
        await request('/api/login/', 'DELETE');
        window.location.reload();
    });
}

document.addEventListener('DOMContentLoaded', async () => {
    initGlobalActions();
    initMenuNavigation();

    try {
        const response = await request('/api/login/', 'HEAD');
        if (response.status === 204) {
            pageBody.classList.remove('jinya-page-content--login');
            if (window.location.hash === '#users') {
                await displayUsers();
            } else if (window.location.hash === '#import-export') {
                await displayImportExport();
            } else {
                await displayJobs();
            }
        } else {
            displayLogin();
        }
    } catch (e) {
        displayLogin();
    }
    document.querySelector('.cosmo-page-layout').classList.remove('jinya-hidden');
});