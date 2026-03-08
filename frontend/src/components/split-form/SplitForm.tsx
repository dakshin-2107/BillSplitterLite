import './SplitForm.css';
import { useState, useEffect } from 'react';
import type { BillFormData, ApiResponse, SplitData } from '@utils/interfaces';
import BillCarousel from '@split-form/bill-carousel/BillCarousel';
import BillForm from '@split-form/bill-form/BillForm';
import PeoplePanel from '@split-form/people-panel/PeoplePanel';
import { Button } from '@ui/button';
import { Spinner } from '@ui/spinner';
import { URLProvider } from '@utils/urlProvider';

interface SplitFormProps {
    onSubmitSuccess: (data: SplitData) => void;
    onSubmitError: (message: string) => void;
}

const SplitForm = ({ onSubmitSuccess, onSubmitError }: SplitFormProps) => {
    const [newBill, setNewBill] = useState<BillFormData | null>(null);
    const [editBill, setEditBill] = useState<BillFormData | null>(null);
    const [bills, setBills] = useState<BillFormData[]>([]);
    const [people, setPeople] = useState<string[]>([]);
    const [isSubmitting, setIsSubmitting] = useState(false);

    useEffect(() => {
        if (!newBill) return;
        setBills(prev => [...prev, newBill]);
        setNewBill(null);
    }, [newBill]);

    const handleStartSession = async () => {
        setIsSubmitting(true);
        try {
            const formData = new FormData();
            bills.forEach(bill => {
                formData.append(import.meta.env.VITE_FORM_IMAGES, bill.image);
                formData.append(import.meta.env.VITE_FORM_DATES, bill.date);
                formData.append(import.meta.env.VITE_FORM_LOCATIONS, bill.location);
            });
            people.forEach(person => {
                formData.append(import.meta.env.VITE_FORM_PEOPLE, person);
            });

            const res = await fetch(URLProvider.getHomeUrl(), {
                method: 'POST',
                credentials: 'include',
                body: formData,
            });

            const json: ApiResponse = await res.json();
            if (json.success && json.split) {
                onSubmitSuccess(json.split);
            } else {
                onSubmitError(json.message);
            }
        } catch (e) {
            onSubmitError(e instanceof Error ? e.message : 'Request failed');
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleResetAll = () => {
        setBills([]);
        setPeople([]);
        setNewBill(null);
    };

    return (
        <div className="split-form-container">
            <BillCarousel bills={bills} setBills={setBills} setEditBill={setEditBill} />
            <div className="main-content-grid">
                <BillForm setNewBill={setNewBill} />
                <div className="people-panel-column">
                    <div className="people-panel-actions">
                        {isSubmitting
                            ? <>
                                <Spinner className="size-10 m-auto mb-5 mt-5" />
                            </>
                            : <>
                                <Button type="button" variant="destructive" onClick={handleResetAll} disabled={isSubmitting}>
                                    Reset all
                                </Button>
                                <Button
                                    type="button"
                                    onClick={handleStartSession}
                                    disabled={!bills.length || !people.length}
                                >
                                    Start session
                                </Button>
                            </>
                        }
                    </div>
                    <PeoplePanel people={people} setPeople={setPeople} />
                </div>
            </div>
        </div>
    );
};

export default SplitForm;
