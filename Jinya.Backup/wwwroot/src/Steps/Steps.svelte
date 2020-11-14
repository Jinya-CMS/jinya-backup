<script>
    import {getKey} from "../authenticationStorage";

    export let params;

    let steps = [];
    let job = {};
    let loadPromise = load();

    async function load() {
        const stepsResponse = await fetch(`/api/job/${params.jobId}/step`, {
            headers: {
                AuthKey: getKey(),
            },
        });
        steps = await stepsResponse.json();

        const jobResponse = await fetch(`/api/job/${params.jobId}`, {
            headers: {
                AuthKey: getKey(),
            },
        });
        job = await jobResponse.json();

        return {job, steps};
    }

    async function deleteStep(stepId) {
        const confirmDelete = confirm('Delete step?');
        if (confirmDelete) {
            const response = await fetch(`/api/job/${params.jobId}/step/${stepId}`, {
                method: 'DELETE',
                headers: {
                    AuthKey: getKey(),
                },
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
            headers: {
                AuthKey: getKey(),
            },
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
    <p class="info">Loading job and steps...</p>
{:then data}
    <h2>{data.job.name} – Steps</h2>
    <div>
        <a class="button" href="/job/{params.jobId}/steps/add/0">Add step</a>
    </div>
{:catch error}
    <p class="error">{error.message}</p>
{/await}
{#each steps as step, i}
    <hr>
    <dl>
        <dt><label for="position">Position</label></dt>
        <dd>
            <input type="number" id="position" bind:value={step.position}>
            <button class="button" on:click={() => savePosition(step)}>Move Step</button>
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
        <a class="button" href="/job/{params.jobId}/steps/{step.id}/edit">Edit</a>
        <button class="button" on:click={() => deleteStep(step.id)}>Delete step</button>
    </div>
    <hr>
    <a class="button" href="/job/{params.jobId}/steps/add/{i + 1}">Add step</a>
{/each}

<style>
    div {
        margin-bottom: 1rem;
    }

    input {
        padding: 0.5rem;
        width: 20rem;
        background: #fff;
        border: 1px solid #afafaf;
        border-radius: 4px;
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

    fieldset {
        border: none;
        margin: 0;
        padding: 1rem 0 0;
    }

    legend {
        font-size: 1.25rem;
    }

    label {
        font-weight: bold;
    }

    dd {
        margin-left: 0;
        padding-left: 0;
        margin-bottom: 1rem;
    }

    dt {
        font-weight: bold;
    }

    hr {
        border-bottom: unset;
        border-top: 1px solid #333;
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
</style>