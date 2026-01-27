import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import './components.css'; // Reuse container styles
import { URLProvider } from '../common/urlProvider';

const JoinSession: React.FC = () => {
    const { billId } = useParams<{ billId: string }>();
    const navigate = useNavigate();
    const [status, setStatus] = useState<'joining' | 'error' | 'success'>('joining');
    const [message, setMessage] = useState('Joining bill session...');

    useEffect(() => {
        const joinSession = async () => {
            if (!billId) {
                setStatus('error');
                setMessage('Invalid bill ID in URL');
                return;
            }

            try {
                const response = await fetch(URLProvider.getJoinUrl(billId), {
                    method: 'GET',
                    credentials: 'include'
                });

                if (response.ok) {
                    setStatus('success');
                    setMessage('Successfully joined! Redirecting...');
                    // Redirect to home after a short delay
                    setTimeout(() => {
                        navigate('/');
                    }, 1500);
                } else {
                    const result = await response.json();
                    setStatus('error');
                    setMessage(result.message || 'Failed to join session');
                }
            } catch (err) {
                console.error("Error joining session:", err);
                setStatus('error');
                setMessage('Network error while joining session');
            }
        };

        joinSession();
    }, []);

    return (
        <div className="home-container">
            <div className={`status-card ${status}`}>
                <h2>{status === 'error' ? 'Oops!' : 'One moment...'}</h2>
                <p>{message}</p>
                {status === 'error' && (
                    <button onClick={() => navigate('/')} className="btn-primary" style={{ marginTop: '20px' }}>
                        Join session
                    </button>
                )}
            </div>
        </div>
    );
};

export default JoinSession;
