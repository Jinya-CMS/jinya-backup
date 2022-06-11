import {pageBody, resetContent} from "./utils.js";

export async function displayImportExport() {
    resetContent();
    const template = Handlebars.compile(document.getElementById('importExportTemplate').innerHTML);
    pageBody.innerHTML = template();

    document.getElementById('importDatabase').addEventListener('click', () => {
        const files = document.getElementById('filePicker').files;
        if (files.length !== 0) {
            const file = files[0];
            const reader = new FileReader();
            reader.addEventListener('load', async (e) => {
                const content = e.target.result;
                const result = await fetch('/api/import', {
                    body: content,
                    method: 'POST',
                });
                if (result.status !== 204) {
                    alert(await result.text());
                } else {
                    alert('Imported backup successfully');
                }
            });
            reader.readAsText(file);
        }
    });
}