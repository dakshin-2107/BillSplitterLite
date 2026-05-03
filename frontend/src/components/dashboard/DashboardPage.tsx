import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Dashboard from './Dashboard';
import { type SplitData, type ApiResponse } from '@utils/interfaces';
import { URLProvider } from '@utils/urlProvider';
import '@comp/components.css';

const DashboardPage = () => {
    const navigate = useNavigate();
    const [splitData, setSplitData] = useState<SplitData | null>(null);
    const [isLoadingSession, setIsLoadingSession] = useState(false);

    useEffect(() => {
        const checkExistingSession = async () => {
            setIsLoadingSession(true);
            try {
                const apiUrl = URLProvider.getHomeUrl();
                if (!apiUrl) return;

                const response = await fetch(apiUrl, {
                    method: 'POST',
                    credentials: 'include',
                });

                if (!response.ok) throw new Error(response.statusText);

                const result: ApiResponse = await response.json();
                if (result.success && result.split) {
                    setSplitData(result.split);
                }
            } catch {
                // silent — user starts fresh
            } finally {
                setIsLoadingSession(false);
            }
        };

        checkExistingSession();
    }, []);

    const navigateToSplitter = (data: SplitData) => {
        navigate(`/splitter/${data.splitId}`, { state: { splitData: data } });
    };

    return (
        <div className="home-container">
            <Dashboard
                onStartNewBill={() => navigate('/splitter')}
                activeSession={splitData}
                onResumeSession={navigateToSplitter}
                isLoadingSession={isLoadingSession}
            />
        </div>
    );
};

export default DashboardPage;
