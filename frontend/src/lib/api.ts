import { browser } from '$app/environment';

// Default to localhost for development
const devApiUrl = 'http://127.0.0.1:8081';
const devWsUrl = 'ws://127.0.0.1:8081/ws';

// In production, we might want to use relative URLs if served from the same domain
// or use environment variables defined at build time.
export const API_URL = browser ? (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1' ? devApiUrl : `http://${window.location.hostname}:8081`) : devApiUrl;
export const WS_URL = browser ? (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1' ? devWsUrl : `ws://${window.location.hostname}:8081/ws`) : devWsUrl;
