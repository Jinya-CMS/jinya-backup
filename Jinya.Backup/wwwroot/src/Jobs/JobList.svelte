<script>
    import {getKey} from "../authenticationStorage";

    let jobs = [];
    let loadPromise = load();

    async function load() {
        const response = await fetch('/api/job', {
            headers: {
                AuthKey: getKey(),
            },
        });
        if (response.ok) {
            return await response.json();
        } else {
            throw new Error('Error loading jobs');
        }
    }

    async function enqueue(jobId) {
        const cron = prompt('Please enter a cron schedule');
        if (cron) {
            const response = await fetch(`/api/job/${jobId}/enqueue?cron=${cron}`, {
                method: 'POST',
                headers: {
                    AuthKey: getKey(),
                },
            });
            if (response.ok) {
                alert('Enqueued successfully');
            } else {
                alert('Enqueue failed');
            }
        }
    }
</script>

<div>
    <a class="button" href="/job/new">New Job</a>
</div>
{#await loadPromise}
    <progress></progress>
    <p class="info">Loading jobs...</p>
{:then jobs}
    {#if jobs.length === 0}
        <p class="info">No Jobs found</p>
    {:else}
        <table>
            <thead>
            <tr>
                <th>#</th>
                <th>Name</th>
                <th>Actions</th>
            </tr>
            </thead>
            <tbody>
            {#each jobs as job}
                <tr>
                    <td>{job.id}</td>
                    <td>{job.name}</td>
                    <td>
                        <a class="button" href="/job/{job.id}">Edit</a>
                        <a class="button" href="/job/{job.id}/steps">Steps</a>
                        <button class="button" on:click={() => enqueue(job.id)}>Enqueue</button>
                    </td>
                </tr>
            {/each}
            </tbody>
        </table>
    {/if}
{:catch error}
    <p class="error">{error.message}</p>
{/await}

<style>
    table {
        width: 100%;
        border-spacing: 0;
        border-collapse: collapse;
        margin-top: 1rem;
    }

    th {
        text-align: left;
        border-bottom: 2px solid #333;
        padding: 0.5rem;
    }

    td {
        border-bottom: 1px solid #333;
        padding: 0.5rem;
    }

    tbody tr:nth-child(2n - 1) {
        background: #f3f3f3;
    }

    progress {
        margin: 1rem 0;
        width: 100%;
    }

    .error {
        padding: 0.5rem;
        background: #c4363c;
        border-radius: 5px;
        color: #eee;
    }

    .info {
        padding: 0.5rem;
        background: #61b4c4;
        border-radius: 5px;
        color: #eee;
    }

    .button {
        border-radius: 5px;
        padding: 0.25rem 0.5rem;
        background: #eee;
        border: 1px solid #ccc;
        text-decoration: none;
        color: #333;
        transition: all 0.3s;
        line-height: normal;
        cursor: pointer;
        box-sizing: border-box;
    }

    .button:hover {
        background: #fff;
    }
</style>