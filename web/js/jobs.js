import {confirm, displayModal, pageBody, request, resetContent, resetListClass, resetTabClass,} from "./utils.js";
import {html} from "../lib/js/jinya-html.js";

let jobs;
let jobDetailContainer;

async function displaySpecificJob(id) {
    resetListClass();
    document.querySelector(`[data-job-id="${id}"]`)?.classList.add('cosmo-list__item--active');
    const [job] = jobs.filter(f => f.id === id);
    if (job) {
        async function renderDetails() {
            const backups = [...await (await request(`/api/backup-job/${id}/backup`, 'GET')).json()];
            jobDetailContainer.innerHTML = html`
                <div class="cosmo-tab-control">
                    <div class="cosmo-tab-control__tabs">
                        <a data-action="details"
                           class="cosmo-tab-control__tab-link cosmo-tab-control__tab-link--active">Details</a>
                        <a data-action="backups" class="cosmo-tab-control__tab-link">Backups</a>
                    </div>
                    <div class="cosmo-tab-control__content" data-content-type="details">
                        <div class="cosmo-toolbar">
                            <div class="cosmo-toolbar__group">
                                <button type="button" data-id="${id}" class="cosmo-button" data-action="edit-job">
                                    Edit
                                </button>
                            </div>
                            <div class="cosmo-toolbar__group">
                                <button type="button" data-id="${id}" class="cosmo-button" data-action="delete-job">
                                    Delete
                                </button>
                            </div>
                        </div>
                        <span class="cosmo-title">${job.name}</span>
                        <dl class="cosmo-key-value-list">
                            <dt class="cosmo-key-value-list__key">ID</dt>
                            <dd class="cosmo-key-value-list__value">${job.id}</dd>
                            <dt class="cosmo-key-value-list__key">Type</dt>
                            <dd class="cosmo-key-value-list__value">${job.type}</dd>
                            <dt class="cosmo-key-value-list__key">Host</dt>
                            <dd class="cosmo-key-value-list__value">${job.host}</dd>
                            <dt class="cosmo-key-value-list__key">Port</dt>
                            <dd class="cosmo-key-value-list__value">${job.port}</dd>
                            <dt class="cosmo-key-value-list__key">Remote path</dt>
                            <dd class="cosmo-key-value-list__value">${job.remotePath}</dd>
                            <dt class="cosmo-key-value-list__key">Local path</dt>
                            <dd class="cosmo-key-value-list__value">${job.localPath}</dd>
                        </dl>
                    </div>
                    <div class="cosmo-tab-control__content jinya-hidden" data-content-type="backups">
                        <table class="cosmo-table">
                            <thead>
                            <tr>
                                <th>#</th>
                                <th>Name</th>
                                <th>Backup Date</th>
                                <th>Actions</th>
                            </tr>
                            </thead>
                            <tbody>
                            ${backups.map(item => html`
                                <tr>
                                    <td>${item.id}</td>
                                    <td>${item.name}</td>
                                    <td>${new Date(Date.parse(item.backupDate)).toLocaleString()}</td>
                                    <td>
                                        <div class="cosmo-toolbar__group">
                                            <a target="_blank" href="/api/backup-job/${job.id}/backup/${item.id}"
                                               class="cosmo-button">Download</a>
                                            <button type="button" data-job-id="${job.id}" data-id="${item.id}"
                                                    class="cosmo-button" data-action="delete-backup">Delete
                                            </button>
                                        </div>
                                    </td>
                                </tr>`)}
                            </tbody>
                        </table>
                    </div>
                </div>`;
            const detailsTab = document.querySelector('[data-content-type=details]');
            const backupsTab = document.querySelector('[data-content-type=backups]');
            document.querySelectorAll('[data-action=delete-backup]').forEach((item) => {
                item.addEventListener('click', async (e) => {
                    e.preventDefault();
                    const result = await confirm('Delete backup', 'Do you want to delete the backup?', 'Delete', 'Keep');
                    if (result) {
                        await request(`/api/backup-job/${e.target.getAttribute('data-job-id')}/backup/${e.target.getAttribute('data-id')}`, 'DELETE');
                        await renderDetails();
                        document.querySelector('[data-action=backups]').click();
                    }
                })
            });
            document.querySelector('[data-action=backups]').addEventListener('click', (e) => {
                e.preventDefault();
                detailsTab.classList.add('jinya-hidden');
                backupsTab.classList.remove('jinya-hidden');
                resetTabClass();
                e.target.classList.add('cosmo-tab-control__tab-link--active');
            });
            document.querySelector('[data-action=details]').addEventListener('click', (e) => {
                e.preventDefault();
                detailsTab.classList.remove('jinya-hidden');
                backupsTab.classList.add('jinya-hidden');
                resetTabClass();
                e.target.classList.add('cosmo-tab-control__tab-link--active');
            });
            document.querySelector('[data-action=edit-job]').addEventListener('click', (e) => {
                e.preventDefault();
                const closeModal = displayModal(html`
                    <form class="cosmo-modal" data-role="edit-job-modal">
                        <h1 class="cosmo-modal__title">Edit FTP job</h1>
                        <div class="cosmo-modal__content">
                            <div class="cosmo-input__group">
                                <label class="cosmo-label" for="name">Name</label>
                                <input class="cosmo-input" id="name" placeholder="Name" required type="text"
                                       value="${job.name}">
                                <label class="cosmo-label" for="type">Server type</label>
                                <select class="cosmo-select" id="type" required>
                                    <option ${job.host.type === 'ftp' ? 'selected' : ''} value="ftp">FTP</option>
                                    <option ${job.host.type === 'sftp' ? 'selected' : ''} value="ftp">SFTP</option>
                                </select>
                                <label class="cosmo-label" for="host">Host</label>
                                <input class="cosmo-input" id="host" placeholder="Host" required type="text"
                                       value="${job.host}">
                                <label class="cosmo-label" for="port">Port</label>
                                <input class="cosmo-input" id="port" placeholder="Port" required type="number"
                                       value="${job.port}">
                                <label class="cosmo-label" for="username">Username</label>
                                <input class="cosmo-input" id="username" placeholder="Username" required type="text"
                                       value="${job.username}">
                                <label class="cosmo-label" for="password">Password</label>
                                <input class="cosmo-input" id="password" placeholder="Password" type="password"
                                       value="${job.password}">
                                <label class="cosmo-label" for="remotePath">Remote path</label>
                                <input class="cosmo-input" id="remotePath" placeholder="Remote path" required
                                       type="text" value="${job.remotePath}">
                                <label class="cosmo-label" for="localPath">Local path</label>
                                <input class="cosmo-input" id="localPath" placeholder="Local path" required type="text"
                                       value="${job.localPath}">
                            </div>
                        </div>
                        <div class="cosmo-modal__button-bar">
                            <button class="cosmo-button" type="submit">Save job</button>
                            <button class="cosmo-button" data-role="cancel-job-modal" type="button">Cancel</button>
                        </div>
                    </form>`);
                document.querySelector('[data-role=edit-job-modal]').addEventListener('submit', async (e) => {
                    e.preventDefault();
                    const data = {
                        name: document.getElementById('name').value,
                        host: document.getElementById('host').value,
                        type: document.getElementById('type').value,
                        port: parseInt(document.getElementById('port').value),
                        localPath: document.getElementById('localPath').value,
                        remotePath: document.getElementById('remotePath').value,
                        username: document.getElementById('username').value,
                    };
                    if (document.getElementById('password').value) {
                        data.password = document.getElementById('password').value;
                    }
                    const response = await request(`/api/backup-job/${job.id}`, 'PUT', data);
                    closeModal();
                    if (response.status === 204) {
                        await displayJobs();
                        await displaySpecificJob(job.id);
                    } else {
                        alert('Error updating job');
                    }
                });
                document.querySelector('[data-role=cancel-job-modal]').addEventListener('click', closeModal);
            });
            document.querySelector('[data-action=delete-job]').addEventListener('click', async (e) => {
                e.preventDefault();
                const result = await confirm('Delete job', 'Do you want to delete the job?', 'Delete', 'Keep');
                if (result) {
                    await request(`/api/backup-job/${job.id}`, 'DELETE');
                    await displayJobs();
                }
            });
        }

        await renderDetails();
    }
}

export async function displayJobs() {
    resetContent();
    document.querySelector('[data-menu=login]').classList.add('jinya-hidden');
    document.querySelector('[data-menu=app]').classList.remove('jinya-hidden');
    const response = await request('/api/backup-job', 'GET');
    if (response.status !== 200) {
        alert('Loading failed')
    } else {
        jobs = await response.json();
        pageBody.innerHTML = html`
            <div class="cosmo-page-body__content">
                <div class="cosmo-list">
                    <nav class="cosmo-list__items">
                        ${jobs.map(item => `<a class="cosmo-list__item" data-action="change-job" data-job-id="${item.id}">${item.name}</a>
`)}
                        <button class="cosmo-button cosmo-button--full-width" data-action="add-job">Create Job</button>
                    </nav>
                    <div class="cosmo-list__content">
                    </div>
                </div>
            </div>`;
        jobDetailContainer = document.querySelector('.cosmo-list__content');
        if (jobs.length > 0) {
            await displaySpecificJob(jobs[0].id);
        }
        document.querySelectorAll('[data-action=change-job]').forEach(job => {
            job.addEventListener('click', async (e) => {
                e.preventDefault();
                await displaySpecificJob(e.target.getAttribute('data-job-id'));
            });
        });
        document.querySelector('[data-action=add-job]').addEventListener('click', (e) => {
            e.preventDefault();
            const template = Handlebars.compile(document.getElementById('createJobTemplate').innerHTML);
            const closeModal = displayModal(template());
            document.querySelector('[data-role=add-job-modal]').addEventListener('submit', async (e) => {
                e.preventDefault();
                const data = {
                    name: document.getElementById('name').value,
                    host: document.getElementById('host').value,
                    type: document.getElementById('type').value,
                    port: parseInt(document.getElementById('port').value),
                    localPath: document.getElementById('localPath').value,
                    remotePath: document.getElementById('remotePath').value,
                    password: document.getElementById('password').value,
                    username: document.getElementById('username').value,
                };
                const response = await request('/api/backup-job', 'POST', data);
                closeModal();
                if (response.status === 201) {
                    await displayJobs();
                } else {
                    alert('Error creating job');
                }
            });
            document.querySelector('[data-role=cancel-job-modal]').addEventListener('click', closeModal);
        });
    }
}
