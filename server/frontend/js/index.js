let authenticationToken = '';

async function request(url, method, body = {}) {
  const data = {
    method,
    // headers: new Headers({ 'Jinya-Auth-Key': authenticationToken })
  };
  if (method !== 'GET' && method !== 'HEAD') {
    data.body = JSON.stringify(body);
  }

  return await fetch(url, data);
}


$(() => {
  const $content = $('[data-role=content]');
  const $loginModal = $('[data-role=login-modal]');
  const $editUserModal = $('[data-role=edit-user-modal]');
  const $deleteUserModal = $('[data-role=delete-user-modal]');
  const $addUserModal = $('[data-role=add-user-modal]');

  const $addJobModal = $('[data-role=add-job-modal]');
  const $editJobModal = $('[data-role=edit-job-modal]');
  const $deleteJobModal = $('[data-role=delete-job-modal]');
  const $deleteBackupModal = $('[data-role=delete-backup-modal]');
  let currentJobId;

  async function renderBackups(jobId) {
    currentJobId = jobId;
    const response = await request(`/api/backup-job/${jobId}/backup`, 'GET');
    const backups = await response.json();
    const tmpl = Handlebars.compile($('#backups').html());
    const result = tmpl({ backups });
    $content.html(result).after(() => {
      $content.find('[data-action=delete]').on('click', (e) => {
        const id = $(e.currentTarget).data('id');
        $deleteBackupModal.data('id', id);
        $deleteBackupModal.modal('show');
      });
    });
  }

  $loginModal.modal({
    closable: false,
    onHide() {
      return !!authenticationToken;
    }
  });
  $loginModal.modal('show');
  $loginModal.on('submit', async (e) => {
    e.preventDefault();
    const username = $('#username').val();
    const password = $('#password').val();
    const response = await request('/api/login/', 'POST', { username, password });
    if (response.status === 200) {
      const body = await response.json();
      authenticationToken = body.token;
      $loginModal.modal('hide');
    } else {
      $('[data-role=login-form]').addClass('error');
    }
  });

  $addUserModal.on('submit', async (e) => {
    e.preventDefault();
    const username = $addUserModal.find('[name=username]').val();
    const password = $addUserModal.find('[name=password]').val();
    const data = {
      username,
      password,
    };

    await request(`/api/user/`, 'POST', data);
    $('[data-action=users]').click();
  });
  $editUserModal.on('submit', async (e) => {
    e.preventDefault();
    const id = $editUserModal.data('id');
    const username = $editUserModal.find('[name=username]').val();
    const password = $editUserModal.find('[name=password]').val();
    const data = {
      username,
    };
    if (password) {
      data.password = password;
    }

    await request(`/api/user/${id}`, 'PUT', data);
    $('[data-action=users]').click();
  });
  $deleteUserModal.find('[data-action=delete]').on('click', async (e) => {
    e.preventDefault();
    const id = $deleteUserModal.data('id');
    await request(`/api/user/${id}`, 'DELETE');
    $('[data-action=users]').click();
  });

  $addJobModal.on('submit', async (e) => {
    e.preventDefault();
    const name = $addJobModal.find('[name=name]').val();
    const host = $addJobModal.find('[name=host]').val();
    const port = $addJobModal.find('[name=port]').val();
    const username = $addJobModal.find('[name=username]').val();
    const password = $addJobModal.find('[name=password]').val();
    const remotePath = $addJobModal.find('[name=remote_path]').val();
    const localPath = $addJobModal.find('[name=local_path]').val();
    const data = {
      name,
      host,
      port: parseInt(port),
      username,
      password,
      remotePath,
      localPath,
    };

    await request('/api/backup-job/', 'POST', data);
    $('[data-action=jobs]').click();
  });
  $editJobModal.on('submit', async (e) => {
    e.preventDefault();
    const id = $editJobModal.data('id');
    const name = $editJobModal.find('[name=name]').val();
    const host = $editJobModal.find('[name=host]').val();
    const port = $editJobModal.find('[name=port]').val();
    const username = $editJobModal.find('[name=username]').val();
    const password = $editJobModal.find('[name=password]').val();
    const remotePath = $editJobModal.find('[name=remote_path]').val();
    const localPath = $editJobModal.find('[name=local_path]').val();
    const data = {
      name,
      host,
      port: parseInt(port),
      username,
      remotePath,
      localPath,
    };
    if (password) {
      data.password = password;
    }

    await request(`/api/backup-job/${id}`, 'PUT', data);
    $('[data-action=jobs]').click();
  });
  $deleteJobModal.find('[data-action=delete]').on('click', async (e) => {
    e.preventDefault();
    const id = $deleteJobModal.data('id');
    await request(`/api/backup-job/${id}`, 'DELETE');
    $('[data-action=jobs]').click();
  });
  $deleteBackupModal.find('[data-action=delete]').on('click', async (e) => {
    e.preventDefault();
    const id = $deleteBackupModal.data('id');
    await request(`/api/backup-job/asd/backup/${id}`, 'DELETE');
    await renderBackups(currentJobId);
  });

  $('[data-action=users]').on('click', async () => {
    const response = await request('/api/user/', 'GET');
    if (response.status === 200) {
      const users = await response.json();
      const tmpl = Handlebars.compile($('#users').html());
      const result = tmpl({ users });
      $content.html(result).after(() => {
        $content.find('[data-action=edit]').on('click', async (e) => {
          const id = $(e.currentTarget).data('id');
          const user = await request(`/api/user/${id}`, 'GET');
          $editUserModal.find('[name=username]').val(user.name);
          $editUserModal.data('id', id);
          $editUserModal.modal('show');
        });
        $content.find('[data-action=delete]').on('click', (e) => {
          const id = $(e.currentTarget).data('id');
          $deleteUserModal.data('id', id);
          $deleteUserModal.modal('show');
        });
        $content.find('[data-action=add]').on('click', () => {
          $addUserModal.find('input').val('');
          $addUserModal.modal('show');
        });
      });
    }
  });
  $('[data-action=jobs]').on('click', async () => {
    const response = await request('/api/backup-job/', 'GET');
    if (response.status === 200) {
      const jobs = await response.json();
      const tmpl = Handlebars.compile($('#jobs').html());
      const result = tmpl({ jobs });
      $content.html(result).after(() => {
        $content.find('[data-action=edit]').on('click', async (e) => {
          const id = $(e.currentTarget).data('id');
          const response = await request(`/api/backup-job/${id}`, 'GET');
          const job = await response.json();
          $editJobModal.find('[name=name]').val(job.name);
          $editJobModal.find('[name=host]').val(job.host);
          $editJobModal.find('[name=port]').val(job.port);
          $editJobModal.find('[name=username]').val(job.username);
          $editJobModal.find('[name=remote_path]').val(job.remotePath);
          $editJobModal.find('[name=local_path]').val(job.localPath);
          $editJobModal.data('id', id);
          $editJobModal.modal('show');
        });
        $content.find('[data-action=delete]').on('click', (e) => {
          const id = $(e.currentTarget).data('id');
          $deleteJobModal.data('id', id);
          $deleteJobModal.modal('show');
        })
        $content.find('[data-action=add]').on('click', () => {
          $addJobModal.find('input').val('');
          $addJobModal.modal('show');
        });
        $content.find('[data-action=backups]').on('click', async (e) => {
          await renderBackups($(e.currentTarget).data('id'));
        });
      });
    }
  });
});