const getBaseUrl = () => import.meta.env.VITE_BASE_URL || '';
const getHomeEndpoint = () => import.meta.env.VITE_BASE_URL_HOME_ENDPOINT || 'home';
const getCloseEndpoint = () => import.meta.env.VITE_CLOSE_ENDPOINT || 'close';
const getWsBaseUrl = () => import.meta.env.VITE_WS_BASE_URL || '';
const getWsActionsEndpoint = () => import.meta.env.VITE_WS_BASE_URL_ACTIONS_ENDPOINT || 'actions';

export const URLProvider = {
    getHomeUrl: () => `${getBaseUrl()}/${getHomeEndpoint()}`,

    getJoinUrl: (billId: string) => `${getBaseUrl()}/join/${billId}`,

    getGenerateUrl: () => `${getBaseUrl()}/generate`,

    getCloseUrl: () => `${getBaseUrl()}/${getCloseEndpoint()}`,

    getPingUrl: () => `${import.meta.env.VITE_BACKEND_BASE_URL || getBaseUrl()}/ping`,

    getSocketUrl: () => `${getWsBaseUrl()}/${getWsActionsEndpoint()}`,

    // Helper to get site origin for sharing
    getSiteOrigin: () => window.location.origin
};
