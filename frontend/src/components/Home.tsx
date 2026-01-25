import React, { useState, useEffect } from 'react';
import './components.css';
import SplitForm from './split-form/SplitForm';
import SplitGrid from './split-grid/SplitGrid';
import { type SplitData, type ApiResponse } from '../common/interfaces';
import { URLProvider } from '../common/urlProvider';

const Home: React.FC = () => {
    const [splitData, setSplitData] = useState<SplitData | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(false);

    // checks for existing session
    useEffect(() => {
        const checkExistingSession = async () => {

            setIsLoading(true);
            try {
                const apiUrl = URLProvider.getHomeUrl();
                if (!apiUrl) {
                    console.error("Backend URL not found");
                    return;
                }

                const response = await fetch(apiUrl, {
                    method: 'POST',
                    credentials: 'include', // FYI: this is required for cookies to be sent
                });

                if (!response.ok) {
                    throw new Error(`Server error: ${response.statusText}`);
                }

                const result: ApiResponse = await response.json();

                if (result.success && result.split) {
                    setSplitData(result.split);
                } else {
                    console.error("Failed to restore session:", result.message);
                }
            } catch (err) {
                console.error("Error restoring session:", err);
            } finally {
                setIsLoading(false);
            }
        };

        checkExistingSession();
    }, []);

    // Handle successful form submission
    const handleFormSuccess = (data: SplitData) => {
        setSplitData(data);
        setError(null);
    };

    // Handle form submission error
    const handleFormError = (errorMessage: string) => {
        setError(errorMessage);
        setSplitData(null);
    };

    // If we have split data, show the SplitGrid
    if (splitData) {
        return (
            <div className="home-container">
                <SplitGrid splitData={splitData} setSplitData={setSplitData} />
            </div>
        );
    }

    // Show the form by default
    return (
        <div className="home-container">
            {isLoading && (
                <div className="loading-overlay">
                    <div className="loading-spinner">Restoring session...</div>
                </div>
            )}
            {error && (
                <div className="error-banner">
                    <p>{error}</p>
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

