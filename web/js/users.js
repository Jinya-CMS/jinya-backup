import {confirm, displayModal, pageBody, request, resetContent} from "./utils.js";
import {html} from "../lib/js/jinya-html.js";

export async function displayUsers() {
    resetContent();
    const users = await (await request('/api/user', 'GET')).json()
    pageBody.innerHTML = html`
        <div class="cosmo-page-body__content">
            <div class="cosmo-toolbar">
                <div class="cosmo-toolbar__group">
                    <button data-action="create-user" class="cosmo-button" type="button">Create user</button>
                </div>
            </div>
            <table class="cosmo-table">
                <thead>
                <tr>
                    <th>#</th>
                    <th>Name</th>
                    <th>Actions</th>
                </tr>
                </thead>
                <tbody>
                ${users.map(item => `<tr>
                    <td>${item.id}</td>
                    <td>${item.name}</td>
                    <td>
                        <div class="cosmo-toolbar__group">
                            <button type="button" data-id="${item.id}" data-action="edit-user" class="cosmo-button">Edit
                            </button>
                            <button type="button" data-id="${item.id}" class="cosmo-button" data-action="delete-user">Delete
                            </button>
                        </div>
                    </td>
                </tr>`)}
                </tbody>
            </table>
        </div>`;
    document.querySelectorAll('[data-action=delete-user]').forEach((item) => {
        item.addEventListener('click', async (e) => {
            e.preventDefault();
            const id = item.getAttribute('data-id');
            const result = await confirm('Delete user', 'Do you really want to delete the user?', 'Delete', 'Keep');
            if (result) {
                const response = await request(`/api/user/${id}`, 'DELETE');
                if (response.status !== 204) {
                    alert('Failed to delete');
                } else {
                    displayUsers();
                }
            }
        });
    });
    document.querySelectorAll('[data-action=edit-user]').forEach((item) => {
        item.addEventListener('click', async (e) => {
            const user = users.filter(f => f.id === item.getAttribute('data-id'))[0];
            const closeModal = await displayModal(html`
                <form class="cosmo-modal" data-role="edit-user-modal">
                    <div class="cosmo-modal__title">Edit user</div>
                    <div class="cosmo-modal__content">
                        <div class="cosmo-input__group">
                            <label class="cosmo-label" for="username">Username</label>
                            <input class="cosmo-input" autocomplete="false" id="username" name="username"
                                   placeholder="Username"
                                   required type="text"
                                   value="${user.name}">
                            <label class="cosmo-label" for="password">Password</label>
                            <input class="cosmo-input" autocomplete="false" id="password" name="password"
                                   placeholder="Password"
                                   required
                                   type="password">
                        </div>
                    </div>
                    <div class="cosmo-modal__button-bar">
                        <button class="cosmo-button" type="submit">Save user</button>
                        <button class="cosmo-button" data-action="cancel-modal" type="button">Cancel</button>
                    </div>
                </form>`);
            document.querySelector('[data-role=edit-user-modal]').addEventListener('submit', async (e) => {
                e.preventDefault();
                const data = {
                    username: document.getElementById('username').value,
                };
                if (document.getElementById('password').value) {
                    data.password = document.getElementById('password').value;
                }
                if ((await request(`/api/user/${item.getAttribute('data-id')}`, 'PUT', data)).status === 204) {
                    closeModal();
                    await displayUsers();
                } else {
                    alert('Failed to update user');
                }
            });
            document.querySelector('[data-action=cancel-modal]').addEventListener('click', closeModal);
        });
    });
    document.querySelectorAll('[data-action=create-user]').forEach((item) => {
        item.addEventListener('click', async (e) => {
            const user = users.filter(f => f.id === item.getAttribute('data-id'));
            const closeModal = await displayModal(html`
                <form class="cosmo-modal" data-role="add-user-modal">
                    <div class="cosmo-modal__title">Add user</div>
                    <div class="cosmo-modal__content">
                        <div class="cosmo-input__group" data-role="login-form">
                            <label class="cosmo-label" for="username">Username</label>
                            <input class="cosmo-input" id="username" name="username" placeholder="Username" required
                                   type="text">
                            <label class="cosmo-label" for="password">Password</label>
                            <input class="cosmo-input" id="password" name="password" placeholder="Password" required
                                   type="password">
                        </div>
                    </div>
                    <div class="cosmo-modal__button-bar">
                        <button class="cosmo-button" type="submit">Save user</button>
                        <button class="cosmo-button" data-action="cancel-modal" type="button">Cancel</button>
                    </div>
                </form>`);
            document.querySelector('[data-role=add-user-modal]').addEventListener('submit', async (e) => {
                e.preventDefault();
                const data = {
                    username: document.getElementById('username').value,
                    password: document.getElementById('password').value,
                };
                if ((await request(`/api/user`, 'POST', data)).status === 201) {
                    closeModal();
                    await displayUsers();
                } else {
                    alert('Failed to create user');
                }
            });
            document.querySelector('[data-action=cancel-modal]').addEventListener('click', closeModal);
        });
    });
}