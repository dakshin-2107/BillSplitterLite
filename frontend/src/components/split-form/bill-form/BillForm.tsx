import { useState, useRef, useEffect, type ChangeEvent } from 'react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { DatePicker } from '@/components/ui/datepicker';
import type { BillFormData } from '../../../common/interfaces';
import './BillForm.css';

interface BillFormProps {
    onBillAdded: (bill: BillFormData) => void;
}

const STORAGE_FORM = 'splitzy_form';

const BillForm = ({ onBillAdded }: BillFormProps) => {
    const [location, setLocation] = useState('');
    const [date, setDate] = useState('');
    const [image, setImage] = useState<File | null>(null);
    const [imagePreview, setImagePreview] = useState<string | null>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);

    useEffect(() => {
        const stored = sessionStorage.getItem(STORAGE_FORM);
        if (stored) {
            const { location: loc, date: d } = JSON.parse(stored);
            if (loc) setLocation(loc);
            if (d) setDate(d);
        }
    }, []);

    useEffect(() => {
        sessionStorage.setItem(STORAGE_FORM, JSON.stringify({ location, date }));
    }, [location, date]);

    const handleImageChange = (e: ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;
        if (imagePreview) URL.revokeObjectURL(imagePreview);
        setImage(file);
        setImagePreview(URL.createObjectURL(file));
    };

    const reset = () => {
        setLocation('');
        setDate('');
        setImage(null);
        if (imagePreview) URL.revokeObjectURL(imagePreview);
        setImagePreview(null);
        if (fileInputRef.current) fileInputRef.current.value = '';
        sessionStorage.removeItem(STORAGE_FORM);
    };

    const handleSubmit = () => {
        if (location && date && image) {
            onBillAdded({ location, date, image });
            reset();
        }
    };

    return (
        <div className='bill-form'>
            <div className="bill-form__header">
                <div>
                    <span className="step-label">Step 1</span>
                    <h2 className="bill-form__title">Add the bills you want to split</h2>
                </div>
            </div>

            <form className="bill-form__body">
                <div className="bill-form__fields">
                    <div className="bill-form__field">
                        <label className="bill-form__label" htmlFor="bill-location">Location</label>
                        <Input
                            id='bill-location'
                            type='text'
                            placeholder="Dunder Mifflin"
                            required
                            value={location}
                            onChange={(e) => setLocation(e.target.value)}
                        />
                    </div>
                    <div className="bill-form__field">
                        <label className="bill-form__label" htmlFor="bill-date">Date</label>
                        <DatePicker required value={date} onDateChange={setDate} />
                    </div>
                </div>

                <div className="bill-form__image-field">
                    <label className="bill-form__label">Bill Image</label>
                    {imagePreview ? (
                        <img
                            src={imagePreview}
                            alt="Bill Preview"
                            className="image-preview"
                            onClick={() => fileInputRef.current?.click()}
                        />
                    ) : (
                        <label htmlFor="image-input" className="image-upload-zone">
                            <span className="image-upload-icon">↑</span>
                            <span className="image-upload-text">Click to upload bill image</span>
                            <span className="image-upload-hint">PNG, JPG, WEBP supported</span>
                        </label>
                    )}
                    <Input
                        required
                        id='image-input'
                        className="image-input"
                        type="file"
                        ref={fileInputRef}
                        accept="image/*"
                        onChange={handleImageChange}
                    />
                </div>

                <div className="form-actions">
                    <Button
                        type="button"
                        variant="destructive"
                        onClick={reset}
                        disabled={!(location || date || image)}
                    >
                        Reset
                    </Button>
                    <Button
                        type="button"
                        onClick={handleSubmit}
                        disabled={!(location && date && image)}
                    >
                        Add bill
                    </Button>
                </div>
            </form>
        </div>
    );
};

export default BillForm;
