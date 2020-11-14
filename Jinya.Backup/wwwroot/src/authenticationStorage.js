export function saveKey(key) {
    localStorage.setItem('/jinya/backup/key', key);
}

export function getKey() {
    return localStorage.getItem('/jinya/backup/key');
}
export function removeKey() {
    localStorage.removeItem('/jinya/backup/key');
}