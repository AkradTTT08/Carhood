import { browser } from '$app/environment';

// Default to localhost for development
const devApiUrl = 'http://127.0.0.1:8081';
const devWsUrl = 'ws://127.0.0.1:8081/ws';

// In production, we might want to use relative URLs if served from the same domain
// or use environment variables defined at build time.
export const API_URL = browser ? (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1' ? devApiUrl : `http://${window.location.hostname}:8081`) : devApiUrl;
export const WS_URL = browser ? (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1' ? devWsUrl : `ws://${window.location.hostname}:8081/ws`) : devWsUrl;

/**
 * Resolves an image URL by:
 * 1. Replacing legacy 'localhost:8081' or '127.0.0.1:8081' with the current dynamic API_URL
 * 2. Prepending API_URL to relative paths starting with '/uploads'
 */
export function resolveImageUrl(url: string | null | undefined): string {
    if (!url) return '';

    // Handle legacy hardcoded localhost/127.0.0.1 URLs
    if (url.includes('localhost:8081') || url.includes('127.0.0.1:8081')) {
        const path = url.split('/uploads/')[1];
        return `${API_URL}/uploads/${path}`;
    }

    // Handle relative paths from the new backend logic
    if (url.startsWith('/uploads/')) {
        return `${API_URL}${url}`;
    }

    return url;
}
