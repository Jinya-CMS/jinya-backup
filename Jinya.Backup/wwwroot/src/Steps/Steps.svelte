<script>
    export let params;

    let steps = [];
    let job = {};
    let loadPromise = load();

    async function load() {
        const stepsResponse = await fetch(`/api/job/${params.jobId}/step`);
        steps = await stepsResponse.json();

        const jobResponse = await fetch(`/api/job/${params.jobId}`);
        job = await jobResponse.json();

        return {job, steps};
    }

    async function deleteStep(stepId) {
        const confirmDelete = confirm('Delete step?');
        if (confirmDelete) {
            const response = await fetch(`/api/job/${params.jobId}/step/${stepId}`, {
                method: 'DELETE',
            });
            if (response.ok) {
                loadPromise = load();
            } else {
                alert('Failed to delete step');
            }
        }
    }
    
    async function savePosition(step) {
        const response = await fetch(`/api/job/${params.jobId}/step/${step.id}/move/${step.position}`, {
            method: 'PUT',
        });
        if (response.ok) {
            loadPromise = load();
        } else {
            alert('Failed to move step');
        }
    }
</script>

{#await loadPromise}
    <progress></progress>
    <p>Loading job and steps...</p>
{:then data}
    <h1>{data.job.name} – Steps</h1>
    <div>
        <a href="/job/{params.jobId}/steps/add/0">Add step</a>
    </div>
{:catch error}
    <p>{error.message}</p>
{/await}
{#each steps as step, i}
    <hr>
    <dl>
        <dt><label for="position">Position</label></dt>
        <dd>
            <input type="number" id="position" bind:value={step.position}><button on:click={() => savePosition(step)}>Move Step</button>
        </dd>
        <dt>Type</dt>
        <dd>
            {#if step.stepType === 1}
                SSH
            {:else if step.stepType === 2}
                File transfer
            {:else if step.stepType === 3}
                Cleanup
            {/if}
        </dd>
        {#if step.stepType === 1 || step.stepType === 2}
            <dt>Server</dt>
            <dd>{step.server}</dd>
        {/if}
        {#if step.stepType === 1 || step.stepType === 2}
            <dt>Username</dt>
            <dd><code>{step.username}</code></dd>
        {/if}
        {#if step.stepType === 1 || step.stepType === 2}
            <dt>Password</dt>
            <dd><code>{step.password}</code></dd>
        {/if}
        {#if step.stepType === 2}
            <dt>Source</dt>
            <dd><code>{step.source}</code></dd>
        {/if}
        {#if step.stepType === 2 || step.stepType === 3}
            <dt>Target</dt>
            <dd><code>{step.target}</code></dd>
        {/if}
        {#if step.stepType === 1}
            <dt>Command</dt>
            <dd>
                <pre>{step.command}</pre>
            </dd>
        {/if}
    </dl>
    <div>
        <a href="/job/{params.jobId}/steps/{step.id}/edit">Edit</a>
        <button on:click={() => deleteStep(step.id)}>Delete step</button>
    </div>
    <hr>
    <a href="/job/{params.jobId}/steps/add/{i + 1}">Add step</a>
{/each}