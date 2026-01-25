import React, { useState } from 'react';
import './SplitForm.css';
import type { ApiResponse, SplitData, BillFormData, SplitFormData } from '../../common/interfaces';
import { URLProvider } from '../../common/urlProvider';
import BillCarousel from './bill-carousel/BillCarousel';
import BillForm from './bill-form/BillForm';
import PeoplePanel from './people-panel/PeoplePanel';

interface SplitFormProps {
    onSubmitSuccess: (splitData: SplitData) => void;
    onSubmitError: (error: string) => void;
}

const SplitForm: React.FC<SplitFormProps> = ({ onSubmitSuccess, onSubmitError }) => {

    const [splitFormData, setSplitFormData] = useState<SplitFormData>({
        billData: [],
        peopleList: [],
    });

    const [isSubmitting, setIsSubmitting] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // bill form stuff
    const handleAddBill = (bill: BillFormData) => {
        setSplitFormData((prev: SplitFormData) => ({
            ...prev,
            billData: [...prev.billData, bill]
        }));
    };

    // bill carousel stuff 
    const handleEditBill = (index: number) => {
        console.log('Edit bill at index:', index);
    };

    const handleDeleteBill = (index: number) => {
        setSplitFormData((prev: SplitFormData) => ({
            ...prev,
            billData: prev.billData.filter((_, i) => i !== index)
        }));
    };

    const handleResetAll = () => {
        setSplitFormData((prev: SplitFormData) => ({
            ...prev,
            billData: []
        }));
    };


    // people panel stuff
    const handleSetPeople = (people: string[]) => {
        setSplitFormData((prev: SplitFormData) => ({
            ...prev,
            peopleList: people
        }));
    };

    // split form stuff
    const handleStartSession = async () => {
        setError(null);
        if (splitFormData.billData.length === 0) {
            setError('Please add at least one bill');
            return;
        }
        if (splitFormData.peopleList.length === 0) {
            setError('Please add at least one person');
            return;
        }

        setIsSubmitting(true);

        try {

            const submitData = new FormData();
            splitFormData.peopleList.forEach((person: string) => submitData.append(import.meta.env.VITE_FORM_PEOPLE, person));
            splitFormData.billData.forEach((bill: BillFormData) => {
                submitData.append(import.meta.env.VITE_FORM_IMAGES, bill.image);
                submitData.append(import.meta.env.VITE_FORM_DATES, bill.date);
                submitData.append(import.meta.env.VITE_FORM_LOCATIONS, bill.location);
            });

            console.log('Submit data:', Array.from(submitData.entries()));

            const apiUrl = URLProvider.getHomeUrl();
            const response = await fetch(apiUrl, {
                method: 'POST',
                body: submitData,
                credentials: 'include',
            });

            if (!response.ok) {
                throw new Error(`Server error: ${response.statusText}`);
            }

            const result: ApiResponse = await response.json();
            if (result.success && result.split) {
                onSubmitSuccess(result.split);
            } else {
                throw new Error(result.message || 'Failed to process session');
            }

        } catch (err) {
            const errorMessage = err instanceof Error ? err.message : 'Failed to start session';
            setError(errorMessage);
            onSubmitError(errorMessage);
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div className="split-form-container">
            <BillCarousel
                splitFormData={splitFormData}
                onClickDeleteBill={handleDeleteBill}
                onClickEditBill={handleEditBill}
                onStartSession={handleStartSession}
                onResetAll={handleResetAll}
                isSubmitting={isSubmitting}
            />

            <div className="main-content-grid">
                <BillForm onAddBill={handleAddBill} />
                <PeoplePanel
                    people={splitFormData.peopleList}
                    setPeople={handleSetPeople}
                />
            </div>

            {error && <p style={{ color: '#ff4444', textAlign: 'center', marginTop: '1rem' }}>{error}</p>}
        </div>
    );
};

export default SplitForm;
