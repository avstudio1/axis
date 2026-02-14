import { useState, useEffect, useCallback } from 'react';
import { fetchJson } from '../utils/fetchJson';

const STATUS_CYCLE = ['Pending', 'Execute'];

export function useRegistry({ addLog, onRegistryChange } = {}) {
    const [mode, setMode] = useState('MANUAL');
    const [registry, setRegistry] = useState([]);
    const [user, setUser] = useState(null);
    const [connected, setConnected] = useState(false);
    const [secondsRemaining, setSecondsRemaining] = useState(null);

    const syncMode = useCallback(async (newMode) => {
        setMode(newMode);
        try {
            await fetchJson(`/api/mode?set=${newMode}`);
        } catch (err) {
            addLog?.('error', `Failed to sync mode ${newMode}`);
        }
    }, [addLog]);

    const fetchRegistry = useCallback(async () => {
        try {
            const data = await fetchJson('/api/registry');
            const list = Array.isArray(data) ? data : [];
            const filtered = list.filter(item => item.type === 'keep');
            setRegistry(filtered);
            onRegistryChange?.(filtered);
            addLog?.('success', 'Manual registry refresh.');
        } catch (err) {
            addLog?.('error', 'Failed to retrieve registry.');
        }
    }, [addLog, onRegistryChange]);

    const fetchDetail = useCallback(async (item) => {
        if (!item || !item.id) {
            throw new Error('Missing item identifier.');
        }
        let url = '';
        switch (item.type) {
            case 'keep':
                url = `/api/notes/detail?id=${encodeURIComponent(item.id)}`;
                break;
            case 'doc':
                url = `/api/docs?id=${encodeURIComponent(item.id)}`;
                break;
            case 'sheet':
                url = `/api/sheets?id=${encodeURIComponent(item.id)}`;
                break;
            default:
                throw new Error(`Unknown item type: ${item.type}`);
        }
        return fetchJson(url);
    }, []);

    const deleteItem = useCallback(async (item) => {
        if (!item || !item.id) return;
        let url = '';
        switch (item.type) {
            case 'keep':
                url = `/api/notes/delete?id=${encodeURIComponent(item.id)}`;
                break;
            case 'doc':
                url = `/api/docs/delete?id=${encodeURIComponent(item.id)}`;
                break;
            case 'sheet':
                url = `/api/sheets/delete?id=${encodeURIComponent(item.id)}`;
                break;
            default:
                addLog?.('error', `Unknown item type for deletion: ${item.type}`);
                return;
        }
        try {
            const res = await fetch(url, { method: 'DELETE' });
            if (!res.ok) throw new Error('Purge request failed');
            addLog?.('success', `Object purged (${item.type}): ${item.id}`);
        } catch (err) {
            addLog?.('error', `Purge failed for ${item.type}: ${item.id}`);
        }
    }, [addLog]);

    const updateStatus = useCallback(async (item, status) => {
        if (!item || !item.id) return;
        try {
            await fetch(`/api/status?id=${encodeURIComponent(item.id)}&status=${status}`, { method: 'POST' });
        } catch (err) {
            addLog?.('error', 'Failed to save status');
            throw err;
        }
    }, [addLog]);

    const nextStatus = useCallback((current, direction) => {
        const idx = STATUS_CYCLE.indexOf(current);
        const safeIdx = idx === -1 ? 0 : idx;
        if (direction === 'forward') {
            return STATUS_CYCLE[(safeIdx + 1) % STATUS_CYCLE.length];
        }
        return STATUS_CYCLE[(safeIdx - 1 + STATUS_CYCLE.length) % STATUS_CYCLE.length];
    }, []);

    useEffect(() => {
        const init = async () => {
            try {
                const userData = await fetchJson('/api/user');
                setUser(userData);
            } catch (e) {
                /* swallow init user errors */
            }

            try {
                const modeData = await fetchJson('/api/mode');
                if (modeData?.mode) {
                    setMode(modeData.mode);
                    addLog?.('system', `State asserted: ${modeData.mode}`);
                }
            } catch (e) {
                /* swallow init mode errors */
            }
        };
        init();
    }, [addLog]);

    useEffect(() => {
        const es = new EventSource('/api/events');
        es.onopen = () => { setConnected(true); addLog?.('success', 'Uplink established (SSE).'); };
        es.onmessage = (e) => {
            try {
                const data = JSON.parse(e.data);
                const list = Array.isArray(data) ? data : [];
                const filtered = list.filter(item => item.type === 'keep');
                setRegistry(filtered);
                onRegistryChange?.(filtered);
                setSecondsRemaining(60);
            } catch (err) { console.error('Stream parse error', err); }
        };

        es.addEventListener('tick', (e) => {
            try {
                const data = JSON.parse(e.data);
                if (data.seconds_remaining !== undefined) {
                    setSecondsRemaining(data.seconds_remaining);
                }
            } catch (err) { console.error('Tick parse error', err); }
        });

        es.addEventListener('status', (e) => {
            try {
                const data = JSON.parse(e.data);
                if (data.status && data.title) {
                    const logType = data.status === 'Execute' ? 'execute' : 'warning';
                    addLog?.(logType, `Status → ${data.status}: ${data.title}`);
                }
            } catch (err) { console.error('Status event parse error', err); }
        });

        es.onerror = () => setConnected(false);
        return () => { es.close(); setConnected(false); };
    }, [addLog, onRegistryChange]);

    return {
        mode,
        registry,
        setRegistry,
        user,
        connected,
        secondsRemaining,
        syncMode,
        fetchRegistry,
        fetchDetail,
        deleteItem,
        updateStatus,
        nextStatus,
    };
}
