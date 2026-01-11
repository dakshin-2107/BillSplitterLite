import React, { useEffect, useState } from 'react';
import { URLProvider } from '../../common/urlProvider';

const Ping: React.FC = () => {
    const [message, setMessage] = useState<string>('Pinging backend...');
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const pingBackend = async () => {
            try {
                const response = await fetch(URLProvider.getPingUrl(), {
                    method: 'GET',
                    credentials: 'include'
                });

                if (!response.ok) {
                    throw new Error(`Ping failed with status: ${response.status}`);
                }

                const data = await response.json();
                setMessage(data);
            } catch (err) {
                console.error("Error pinging backend:", err);
                setError(err instanceof Error ? err.message : 'Failed to ping backend');
            } finally {
                setLoading(false);
            }
        };

        pingBackend();
    }, []);

    if (loading) {
        return (
            <div style={{ padding: '2rem', textAlign: 'center' }}>
                <h2>Pinging...</h2>
                <div className="loading-spinner"></div>
            </div>
        );
    }

    if (error) {
        return (
            <div style={{ padding: '2rem', textAlign: 'center', color: '#ef4444' }}>
                <h2>Error</h2>
                <p>{error}</p>
            </div>
        );
    }

    return (
        <div style={{ padding: '2rem', textAlign: 'center' }}>
            <h2>Backend response:</h2>
            <div style={{
                background: 'rgba(255, 255, 255, 0.05)',
                padding: '1rem',
                borderRadius: '10px',
                display: 'inline-block',
                marginTop: '1rem',
                fontSize: '1.2rem',
                border: '1px solid rgba(255, 255, 255, 0.1)'
            }}>
                {message}
            </div>
        </div>
    );
};

export default Ping;
