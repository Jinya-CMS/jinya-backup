<script>
    import router from 'page';
    import {getKey} from "../authenticationStorage";

    export let params;

    let step = {};
    let loadPromise = loadStep();

    async function loadStep() {
        const response = await fetch(`/api/job/${params.jobId}/step/${params.stepId}`, {
            headers: {
                AuthKey: getKey(),
            },
        });
        step = await response.json();

        return step;
    }

    async function save() {
        step.stepType = parseInt(step.stepType);
        const response = await fetch(`/api/job/${params.jobId}/step/${params.stepId}`, {
            method: 'PUT',
            body: JSON.stringify(step),
            headers: {
                'Content-Type': 'application/json',
                AuthKey: getKey(),
            },
        });
        if (response.ok) {
            router(`/job/${params.jobId}/steps`);
        } else {
            alert('Error saving step');
        }
    }
</script>

{#await loadPromise}
    <progress></progress>
    <p class="info">Step is loading...</p>
{:then step}
{:catch error}
    <p class="error">{error.message}</p>
{/await}
<fieldset>
    <legend>Edit step</legend>
    <div>
        <label for="stepType">Step Type</label><br/>
        <select name="stepType" id="stepType" bind:value={step.stepType}>
            <option value="1">SSH</option>
            <option value="2">File transfer</option>
            <option value="3">Cleanup</option>
        </select>
    </div>
    {#if step.stepType != 3}
        <div>
            <label for="server">Server</label><br/>
            <input type="text" id="server" bind:value={step.server}>
        </div>
        <div>
            <label for="username">Username</label><br/>
            <input type="text" id="username" bind:value={step.username}>
        </div>
        <div>
            <label for="password">Password</label><br/>
            <input type="password" id="password" bind:value={step.password}>
        </div>
        {#if step.stepType == 1}
            <div>
                <label for="command">Command</label><br/>
                <textarea id="command" bind:value={step.command}></textarea>
            </div>
        {/if}
        {#if step.stepType == 2}
            <div>
                <label for="source">Source</label><br/>
                <input type="text" id="source" bind:value={step.source}>
            </div>
        {/if}
    {/if}
    {#if step.stepType == 2 || step.stepType == 3}
        <div>
            <label for="target">Target</label><br/>
            <input type="text" id="target" bind:value={step.target}>
        </div>
    {/if}
    <div>
        <button on:click={save}>Save</button>
    </div>
</fieldset>

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

    button {
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

    button:hover {
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