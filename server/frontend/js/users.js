import {confirm, displayModal, pageBody, request, resetContent} from "./utils.js";

export async function displayUsers() {
    resetContent();
    const users = await (await request('/api/user/', 'GET')).json()
    const template = Handlebars.compile(document.getElementById('userListTemplate').innerHTML);
    pageBody.innerHTML = template({users});
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
            const user = users.filter(f => f.id === item.getAttribute('data-id'));
            const template = Handlebars.compile(document.getElementById('editUserTemplate').innerHTML);
            const closeModal = await displayModal(template(user));
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
            const template = Handlebars.compile(document.getElementById('addUserTemplate').innerHTML);
            const closeModal = await displayModal(template());
            document.querySelector('[data-role=add-user-modal]').addEventListener('submit', async (e) => {
                e.preventDefault();
                const data = {
                    username: document.getElementById('username').value,
                    password: document.getElementById('password').value,
                };
                if ((await request(`/api/user/`, 'POST', data)).status === 201) {
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