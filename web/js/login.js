import {pageBody, request, resetContent} from "./utils.js";
import {displayJobs} from "./jobs.js";
import {html} from "../lib/js/jinya-html.js";

export function displayLogin() {
    resetContent(true);
    document.querySelector('[data-menu=login]').classList.remove('jinya-hidden');
    document.querySelector('[data-menu=app]').classList.add('jinya-hidden');
    pageBody.innerHTML = html`
        <form id="loginForm">
            <div class="cosmo-input__group">
                <label class="cosmo-label" for="username">Username</label>
                <input class="cosmo-input" id="username" required type="text">
                <label class="cosmo-label" for="password">Password</label>
                <input class="cosmo-input" id="password" required type="password">
            </div>
            <div class="cosmo-button__container">
                <button class="cosmo-button" data-action="login" type="submit">Login</button>
            </div>
        </form>`;
    pageBody.classList.add('jinya-page-content--login');
    pageBody.querySelector('#loginForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        try {
            const response = await request('/api/login/', 'POST', {
                username: document.getElementById('username').value,
                password: document.getElementById('password').value,
            });
            if (response.status !== 200) {
                alert('Login failed')
            } else {
                pageBody.classList.remove('jinya-page-content--login');
                await displayJobs();
            }
        } catch (e) {
            alert(e);
        }
    })
}