import {confirm, displayModal, pageBody, request, resetContent, resetListClass, resetTabClass,} from "./utils.js";

let jobs;
const jobDetailTemplate = Handlebars.compile(document.getElementById('jobItemDetailsTemplate').innerHTML);
let jobDetailContainer;

async function displaySpecificJob(id) {
    resetListClass();
    document.querySelector(`[data-job-id="${id}"]`)?.classList.add('cosmo-list__item--active');
    const [job] = jobs.filter(f => f.id === id);
    if (job) {
        async function renderDetails() {
            const backups = [...await (await request(`/api/backup-job/${id}/backup`, 'GET')).json()];
            const templateData = {
                job,
                backups,
            };
            Handlebars.registerHelper('formatDate', (data) => {
                const date = new Date(Date.parse(data));
                return date.toLocaleString();
            })
            jobDetailContainer.innerHTML = jobDetailTemplate(templateData);
            const detailsTab = document.querySelector('[data-content-type=details]');
            const backupsTab = document.querySelector('[data-content-type=backups]');
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
            document.querySelectorAll('[data-action="delete-backup"]').forEach((item) => {
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
            document.querySelector('[data-action="edit-job"]').addEventListener('click', (e) => {
                e.preventDefault();
                const template = Handlebars.compile(document.getElementById('editJobTemplate').innerHTML);
                const closeModal = displayModal(template(job));
                document.querySelector('[data-role=edit-job-modal]').addEventListener('submit', async (e) => {
                    e.preventDefault();
                    const data = {
                        name: document.getElementById('name').value,
                        host: document.getElementById('host').value,
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
                        alert('Error creating job');
                    }
                });
                document.querySelector('[data-role=cancel-job-modal]').addEventListener('click', closeModal);
            });
            document.querySelector('[data-action="delete-job"]').addEventListener('click', async (e) => {
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
        const jobsTemplate = Handlebars.compile(document.getElementById('jobListTemplate').innerHTML);
        jobs = await response.json();
        pageBody.innerHTML = jobsTemplate({jobs});
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
