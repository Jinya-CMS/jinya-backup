<script>
    import {getKey, removeKey} from "./authenticationStorage";
    import router from 'page';
    import LoginPage from './LoginPage.svelte';
    import JobList from './Jobs/JobList.svelte';
    import EditJob from './Jobs/EditJob.svelte';
    import AddJob from './Jobs/AddJob.svelte';
    import Steps from './Steps/Steps.svelte';
    import EditStep from './Steps/EditStep.svelte';
    import AddStep from './Steps/AddStep.svelte';

    let page;
    let params;

    router('/login', () => page = LoginPage);
    router('/', () => page = JobList);
    router('/job/new', () => page = AddJob);
    router('/job/:id', (ctx, next) => {
        params = ctx.params;
        next();
    }, () => page = EditJob);
    router('/job/:jobId/steps', (ctx, next) => {
        params = ctx.params;
        next();
    }, () => page = Steps);
    router('/job/:jobId/steps/add/:position', (ctx, next) => {
        params = ctx.params;
        next();
    }, () => page = AddStep);
    router('/job/:jobId/steps/:stepId/edit', (ctx, next) => {
        params = ctx.params;
        next();
    }, () => page = EditStep);

    router.start();

    if (!getKey()) {
        router('/login');
    }

    function logout() {
        removeKey();
        router('/login');
    }
</script>

<main>
    <header>
        <h1>Jinya Backup</h1>
        <button on:click={logout}>Logout</button>
    </header>

    <svelte:component this={page} params={params}/>
</main>

<style>
    :root {
        color: #333;
    }
    
    header {
        display: flex;
        justify-content: space-between;
        align-items: baseline;
    }

    h1 {
        margin-top: 0;
    }

    main {
        margin-right: 20%;
        margin-left: 20%;
        padding: 1rem;
        background: #fafafa;
        border-radius: 5px;
        margin-top: 2rem;
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
</style>