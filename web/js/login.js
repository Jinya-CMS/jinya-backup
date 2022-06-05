import {pageBody, request, resetContent} from "./utils.js";
import {displayJobs} from "./jobs.js";

export function displayLogin() {
    resetContent(true);
    document.querySelector('[data-menu=login]').classList.remove('jinya-hidden');
    document.querySelector('[data-menu=app]').classList.add('jinya-hidden');
    const loginTemplate = Handlebars.compile(document.getElementById('loginTemplate').innerHTML);
    pageBody.innerHTML = loginTemplate();
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