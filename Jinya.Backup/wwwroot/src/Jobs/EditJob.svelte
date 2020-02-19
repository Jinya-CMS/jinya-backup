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
</style>