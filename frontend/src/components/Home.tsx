import { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import SplitForm from '@split-form/SplitForm';
import SplitGrid from './split-grid/SplitGrid';
import Dashboard from './dashboard/Dashboard';
import { type SplitData, type ApiResponse } from '@utils/interfaces';
import { URLProvider } from '@utils/urlProvider';
import useNotify from '../hooks/useNotify';
import './components.css';

type View = 'dashboard' | 'form' | 'grid';

const Home = () => {
    const notify = useNotify();
    const location = useLocation();
    const initialView: View = (location.state as { view?: View } | null)?.view ?? 'dashboard';
    const [view, setView] = useState<View>(initialView);
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

                if (!response.ok) {
                    throw new Error(`Server error: ${response.statusText}`);
                }

                const result: ApiResponse = await response.json();
                if (result.success && result.split) {
                    setSplitData(result.split);
                }
            } catch {
                // session restore failure is silent — user starts fresh
            } finally {
                setIsLoadingSession(false);
            }
        };

        checkExistingSession();
    }, []);

    const handleFormSuccess = (data: SplitData) => {
        setSplitData(data);
        setView('grid');
        notify.success('Session started successfully');
    };

    const handleFormError = (errorMessage: string) => {
        notify.error(errorMessage);
    };

    if (view === 'grid' && splitData) {
        return (
            <div className="home-container">
                <SplitGrid splitData={splitData} setSplitData={setSplitData} />
            </div>
        );
    }

    if (view === 'form') {
        return (
            <div className="home-container">
                <SplitForm
                    onSubmitSuccess={handleFormSuccess}
                    onSubmitError={handleFormError}
                />
            </div>
        );
    }

    return (
        <div className="home-container">
            <Dashboard
                onStartNewBill={() => setView('form')}
                activeSession={splitData}
                onResumeSession={(data) => {
                    setSplitData(data);
                    setView('grid');
                }}
                isLoadingSession={isLoadingSession}
            />
        </div>
    );
};

export default Home;
