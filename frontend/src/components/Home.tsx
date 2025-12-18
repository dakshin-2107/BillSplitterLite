import React, { useState } from 'react';
import './Home.css';
import BillForm from './bill-form/BillForm';
import SplitGrid from './split-grid/SplitGrid';

interface Item {
    id: string;
    name: string;
    price: string;
    takers: Record<string, unknown>;
}

interface BillData {
    date: string;
    location: string;
    total: string;
    items: Record<string, Item>;
    participants: Record<string, string>;
}

const Home: React.FC = () => {
    const [billData, setBillData] = useState<BillData | null>(null);
    const [error, setError] = useState<string | null>(null);

    // Handle successful form submission
    const handleFormSuccess = (data: BillData) => {
        setBillData(data);
        setError(null);
    };

    // Handle form submission error
    const handleFormError = (errorMessage: string) => {
        setError(errorMessage);
        setBillData(null);
    };

    // Handle going back to form from SplitGrid
    const handleBackToForm = () => {
        setBillData(null);
        setError(null);
    };

    // If we have bill data, show the SplitGrid
    if (billData) {
        return (
            <div className="home-container">
                <button onClick={handleBackToForm} className="btn-back">
                    ← Back to Form
                </button>
                <SplitGrid data={billData} />
            </div>
        );
    }

    // Show the form by default
    return (
        <div className="home-container">
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
