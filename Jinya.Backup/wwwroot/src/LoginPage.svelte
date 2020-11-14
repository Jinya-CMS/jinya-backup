<script>
    import {saveKey} from "./authenticationStorage";
    import router from 'page';

    let authKey;

    function login() {
        saveKey(authKey);
        router('/');
    }
</script>

<fieldset>
    <legend>Login</legend>
    <div>
        <label for="password">Authentication Key</label><br/>
        <input type="password" id="password" bind:value={authKey}>
    </div>
    <div>
        <button on:click={login}>Login</button>
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