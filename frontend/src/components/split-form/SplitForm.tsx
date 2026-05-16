import { useState, useEffect } from 'react';
import type { BillFormData, ApiResponse, SplitData } from '@utils/interfaces';
import BillCarousel from '@split-form/bill-carousel/BillCarousel';
import BillForm from '@split-form/bill-form/BillForm';
import PeoplePanel from '@split-form/people-panel/PeoplePanel';
import { Button } from '@ui/button';
import { Spinner } from '@ui/spinner';
import { URLProvider } from '@utils/urlProvider';
import './SplitForm.css';

const STORAGE_BILLS = 'splitzy_bills';
const STORAGE_PEOPLE = 'splitzy_people';

interface StoredBill {
    location: string;
    date: string;
    imageDataUrl: string;
    imageName: string;
    imageType: string;
}

const billToStored = (bill: BillFormData): Promise<StoredBill> =>
    new Promise((resolve) => {
        const reader = new FileReader();
        reader.onload = (e) => resolve({
            location: bill.location,
            date: bill.date,
            imageDataUrl: e.target!.result as string,
            imageName: bill.image.name,
            imageType: bill.image.type,
        });
        reader.readAsDataURL(bill.image);
    });

const storedToBill = async (stored: StoredBill): Promise<BillFormData> => {
    const res = await fetch(stored.imageDataUrl);
    const blob = await res.blob();
    const image = new File([blob], stored.imageName, { type: stored.imageType });
    return { location: stored.location, date: stored.date, image };
};

interface SplitFormProps {
    onSubmitSuccess: (data: SplitData) => void;
    onSubmitError: (message: string) => void;
}

const SplitForm = ({ onSubmitSuccess, onSubmitError }: SplitFormProps) => {
    const [bills, setBills] = useState<BillFormData[]>([]);
    const [people, setPeople] = useState<string[]>([]);
    const [isSubmitting, setIsSubmitting] = useState(false);

    useEffect(() => {
        const storedPeople = sessionStorage.getItem(STORAGE_PEOPLE);
        if (storedPeople) setPeople(JSON.parse(storedPeople));

        const storedBills = sessionStorage.getItem(STORAGE_BILLS);
        if (storedBills) {
            const parsed: StoredBill[] = JSON.parse(storedBills);
            Promise.all(parsed.map(storedToBill)).then(setBills);
        }
    }, []);

    useEffect(() => {
        sessionStorage.setItem(STORAGE_PEOPLE, JSON.stringify(people));
    }, [people]);

    useEffect(() => {
        Promise.all(bills.map(billToStored)).then(stored => {
            sessionStorage.setItem(STORAGE_BILLS, JSON.stringify(stored));
        });
    }, [bills]);

    const handleBillAdded = (bill: BillFormData) => {
        setBills(prev => [...prev, bill]);
    };

    const handleStartSession = async () => {
        setIsSubmitting(true);
        try {
            const formData = new FormData();
            bills.forEach(bill => {
                formData.append('images', bill.image);
                formData.append('dates', bill.date);
                formData.append('locations', bill.location);
            });
            people.forEach(person => {
                formData.append('people', person);
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
        sessionStorage.removeItem(STORAGE_BILLS);
        sessionStorage.removeItem(STORAGE_PEOPLE);
    };

    return (
        <div className="split-form-container">
            <BillCarousel bills={bills} setBills={setBills} />
            <div className="main-content-grid">
                <BillForm onBillAdded={handleBillAdded} />
                <div className="people-panel-column">
                    <div className="people-panel-actions">
                        {isSubmitting
                            ? <Spinner className="size-10 m-auto mb-5 mt-5" />
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
