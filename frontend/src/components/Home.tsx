import React, { useState, useEffect } from 'react';
import './Home.css';
import BillForm from './bill-form/BillForm';
import SplitGrid from './split-grid/SplitGrid';
import { type BillData, type ApiResponse } from '../common/interfaces';
import { URLProvider } from '../common/urlProvider';

const Home: React.FC = () => {
    const [billData, setBillData] = useState<BillData | null>(null);
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

                if (result.success && result.bill) {
                    setBillData(result.bill);
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
    const handleFormSuccess = (billdata: BillData) => {
        setBillData(billdata);
        setError(null);
    };

    // Handle form submission error
    const handleFormError = (errorMessage: string) => {
        setError(errorMessage);
        setBillData(null);
    };

    // // Handle going back to form from SplitGrid
    // const handleBackToForm = () => {
    //     setBillData(null);
    //     setError(null);
    // };

    // If we have bill data, show the SplitGrid
    if (billData) {
        return (
            <div className="home-container">
                {/* <button onClick={handleBackToForm} className="btn-back">
                    ← Back to Form
                </button> */}
                <SplitGrid billData={billData} setBillData={setBillData} />
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
            <BillForm
                onSubmitSuccess={handleFormSuccess}
                onSubmitError={handleFormError}
            />
        </div>
    );
};

export default Home;
