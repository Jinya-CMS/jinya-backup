<script>
    import {onMount} from 'svelte';
    import router from 'page';

    export let params;

    let name;

    onMount(async () => {
        const response = await fetch(`/api/job/${params.id}`);
        if (response.ok) {
            const json = await response.json();
            name = json.name;
        } else {
            alert('Error loading job');
        }
    });

    async function save() {
        const response = await fetch(`/api/job/${params.id}`, {
            method: 'PUT',
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