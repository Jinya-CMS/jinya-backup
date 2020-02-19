<script>
    import router from 'page';
    import JobList from './Jobs/JobList.svelte';
    import EditJob from './Jobs/EditJob.svelte';
    import AddJob from './Jobs/AddJob.svelte';
    import Steps from './Steps/Steps.svelte';
    import EditStep from './Steps/EditStep.svelte';
    import AddStep from './Steps/AddStep.svelte';

    let page;
    let params;

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
</script>

<main>
    <h1>Jinya Backup</h1>
    <svelte:component this={page} params={params}/>
</main>

<style>
    :root {
        color: #333;
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
</style>