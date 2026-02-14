export async function fetchJson(url, options = {}) {
    const res = await fetch(url, options);
    if (!res.ok) {
        const text = await res.text().catch(() => '');
        const error = new Error(`Request failed: ${res.status}`);
        error.status = res.status;
        error.body = text;
        throw error;
    }
    return res.json();
}
