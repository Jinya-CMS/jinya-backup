import {pageBody, resetContent} from "./utils.js";
import {html} from "../lib/js/jinya-html.js";

export async function displayImportExport() {
    resetContent();
    pageBody.innerHTML = html`
        <div class="cosmo-page-body__content">
            <h3>Export your current database</h3>
            <div>
                <a href="/api/export" target="_blank" class="cosmo-button" id="exportDatabase">Export</a>
            </div>
            <h3>Import backup</h3>
            <div class="cosmo-input__group">
                <label for="filePicker" class="cosmo-label">Pick backup file</label>
                <input id="filePicker" type="file" class="cosmo-input">
            </div>
            <div class="cosmo-button__container">
                <button class="cosmo-button" id="importDatabase">Import</button>
            </div>
        </div>`;

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