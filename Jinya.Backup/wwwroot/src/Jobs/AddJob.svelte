<script>
    import router from 'page';

    export let params;

    let name;

    async function save() {
        const response = await fetch(`/api/job`, {
            method: 'POST',
            body: JSON.stringify({name}),
            headers: {
                'Content-Type': 'application/json',
            },
        });
        if (response.ok) {
            router('/');
        } else {
            alert('Error saving job');
        }
    }
</script>

<div>
    <label for="name">Name</label><br/>
    <input type="text" bind:value={name} id="name">
</div>
<div>
    <button on:click={save}>Save</button>
</div>