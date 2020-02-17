<script>
    let email;
    let password;
    let apiKey;
    
    async function login() {
        var response = await fetch('/api/account/login', {
        	method:'POST',
			body: JSON.stringify({
				email,
				password
			}),
			headers: {
				'Content-Type': 'application/json'
			},
        });
		apiKey = response.headers.get('Jinya-Api-Token');
    }
</script>

<main>
	<label for="email">Email: </label>
    <input type="text" id="email" bind:value={email}>
	<label for="password">Password: </label>
    <input type="password" id="password" bind:value={password}>
    <button on:click={login}>Login</button>
    {apiKey}
</main>

<style>
    main {
        text-align: center;
        padding: 1em;
        max-width: 240px;
        margin: 0 auto;
    }
    
    h1 {
        color: #ff3e00;
        text-transform: uppercase;
        font-size: 4em;
        font-weight: 100;
    }
    
    @media (min-width: 640px) {
        main {
            max-width: none;
        }
    }
</style>