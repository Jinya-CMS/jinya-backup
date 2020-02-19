<script>
    import {onMount} from 'svelte';

    let jobs = [];
    let loadPromise = load();

    async function load() {
        const response = await fetch('/api/job');
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
    <a href="/job/new">New Job</a>
</div>
{#await loadPromise}
    <progress></progress>
    <p>Loading jobs...</p>
{:then jobs}
    {#if jobs.length === 0}
        <div>No Jobs found</div>
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
                        <a href="/job/{job.id}">Edit</a>
                        <a href="/job/{job.id}/steps">Steps</a>
                        <button on:click={() => enqueue(job.id)}>Enqueue</button>
                    </td>
                </tr>
            {/each}
            </tbody>
        </table>
    {/if}
{:catch error}
    <p>{error.message}</p>
{/await}