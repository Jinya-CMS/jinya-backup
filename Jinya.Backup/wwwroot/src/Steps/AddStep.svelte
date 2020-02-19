<script>
    import router from 'page';

    export let params;

    let step = {};

    async function save() {
        step.position = parseInt(params.position);
        step.stepType = parseInt(step.stepType);
        const response = await fetch(`/api/job/${params.jobId}/step`, {
            method: 'POST',
            body: JSON.stringify(step),
            headers: {
                'Content-Type': 'application/json',
            },
        });
        if (response.ok) {
            router(`/job/${params.jobId}/steps`);
        } else {
            alert('Error saving step');
        }
    }
</script>

<fieldset>
    <legend>Add step</legend>
    <div>
        <label for="stepType">Step Type</label>
        <select name="stepType" id="stepType" bind:value={step.stepType}>
            <option value="1">SSH</option>
            <option value="2">File transfer</option>
            <option value="3">Cleanup</option>
        </select>
    </div>
    {#if step.stepType != 3}
        <div>
            <label for="server">Server</label>
            <input type="text" id="server" bind:value={step.server}>
        </div>
        <div>
            <label for="username">Username</label>
            <input type="text" id="username" bind:value={step.username}>
        </div>
        <div>
            <label for="password">Password</label>
            <input type="password" id="password" bind:value={step.password}>
        </div>
        {#if step.stepType == 1}
            <div>
                <label for="command">Command</label>
                <textarea id="command" bind:value={step.command}></textarea>
            </div>
        {/if}
        {#if step.stepType == 2}
            <div>
                <label for="source">Source</label>
                <input type="text" id="source" bind:value={step.source}>
            </div>
        {/if}
    {/if}
    {#if step.stepType == 2 || step.stepType == 3}
        <div>
            <label for="target">Target</label>
            <input type="text" id="target" bind:value={step.target}>
        </div>
    {/if}
    <div>
        <button on:click={save}>Save</button>
    </div>
</fieldset>