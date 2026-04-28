import { useState, useEffect } from 'react';
import { toast } from 'sonner';
import SplitForm from '@split-form/SplitForm';
import SplitGrid from './split-grid/SplitGrid';
import { type SplitData, type ApiResponse } from '@utils/interfaces';
import { URLProvider } from '@utils/urlProvider';
import { Toaster } from '@ui/sonner';
import './components.css';

const Home = () => {
    const [splitData, setSplitData] = useState<SplitData | null>(null);
    const [isLoading, setIsLoading] = useState(false);

    useEffect(() => {
        const checkExistingSession = async () => {
            setIsLoading(true);
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
                setIsLoading(false);
            }
        };

        checkExistingSession();
    }, []);

    const handleFormSuccess = (data: SplitData) => {
        setSplitData(data);
        toast.success('Session started successfully');
    };

    const handleFormError = (errorMessage: string) => {
        toast.error(errorMessage);
    };

    if (splitData) {
        return (
            <div className="home-container">
                <SplitGrid splitData={splitData} setSplitData={setSplitData} />
            </div>
        );
    }

    return (
        <div className="home-container">
            <Toaster />
            {isLoading && (
                <div className="loading-overlay">
                    <div className="loading-spinner">Restoring session...</div>
                </div>
            )}
            <SplitForm
                onSubmitSuccess={handleFormSuccess}
                onSubmitError={handleFormError}
            />
        </div>
    );
};

export default Home;
