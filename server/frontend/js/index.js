import {pageBody, request} from "./utils.js";
import {displayJobs} from "./jobs.js";
import {displayLogin} from "./login.js";

function initGlobalActions() {
    document.querySelector('[data-action=logout]').addEventListener('click', async () => {
        await request('/api/login/', 'DELETE');
        window.location.reload();
    });
}

document.addEventListener('DOMContentLoaded', async () => {
    initGlobalActions();

    try {
        const response = await request('/api/login/', 'HEAD');
        if (response.status === 204) {
            pageBody.classList.remove('jinya-page-content--login');
            await displayJobs();
        } else {
            displayLogin();
        }
    } catch (e) {
        displayLogin();
    }
});