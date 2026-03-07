import React, { useState, useRef, type ChangeEvent, useEffect } from 'react';
import './BillForm.css';
import { FieldLabel, FieldLegend, FieldSet, Field } from '@/components/ui/field';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { DatePicker } from '@/components/ui/datepicker';

interface Bill {
    location: string;
    date: string;
    image: File;
}

interface BillFormProps {
    onAddBill: (bill: Bill) => void;
}

const BillForm: React.FC<BillFormProps> = ({ onAddBill }) => {
    const [location, setLocation] = useState('');
    const [date, setDate] = useState('');
    const [image, setImage] = useState<File | null>(null);
    const [imagePreview, setImagePreview] = useState<string | null>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);

    const [showResetButton, setShowResetButton] = useState(false);
    const [showAddBillButton, setShowAddBillButton] = useState(false);

    const handleImageChange = (e: ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            setImage(file);
            const reader = new FileReader();
            reader.onloadend = () => {
                setImagePreview(reader.result as string);
            };
            reader.readAsDataURL(file);
        }
    };

    const handleReset = () => {
        setLocation('');
        setDate('');
        setImage(null);
        setImagePreview(null);
        if (fileInputRef.current) fileInputRef.current.value = '';
    };

    const handleAddBill = () => {
        if (location && date && image) {
            onAddBill({ location, date, image });
            handleReset();
        }
    };

    useEffect(() => {
        setShowResetButton((date || location || image) ? false : true);
        setShowAddBillButton((date && location && image) ? false : true);
    }, [date, location, image]);

    return (
        <div className='bill-form'>
            <form>
                <FieldSet>
                    <FieldLegend>Bill Details</FieldLegend>
                    <Field>
                        <FieldLabel htmlFor="bill-location">
                            Location
                        </FieldLabel>
                        <Input
                            id='bill-location'
                            type='text'
                            placeholder="Dunder Mifflin"
                            required
                            value={location}
                            onChange={(e) => setLocation(e.target.value)}
                        />
                    </Field>
                    <Field>
                        <FieldLabel htmlFor="bill-date">
                            Date
                        </FieldLabel>
                        <DatePicker
                            required
                            value={date}
                            onDateChange={setDate}
                        />
                    </Field>
                    <Field>
                        <FieldLabel htmlFor="bill-image">
                            Bill image
                        </FieldLabel>
                        {imagePreview ? (
                            <img src={imagePreview} alt="Bill Preview" className="image-preview" />
                        ) : (
                            <>
                                <label htmlFor="image-input" className="image-preview flex items-center justify-center">
                                    <span>Click to upload the bill</span>
                                </label>
                                <Input
                                    required
                                    id='image-input'
                                    className="image-input"
                                    type="file"
                                    ref={fileInputRef}
                                    accept="image/*"
                                    onChange={handleImageChange}
                                />
                            </>
                        )}
                    </Field>
                    <div className="form-actions">
                        <Button
                            variant="destructive"
                            onClick={handleReset}
                            disabled={showResetButton}
                        >
                            Reset
                        </Button>
                        <Button
                            onClick={handleAddBill}
                            disabled={showAddBillButton}
                        >
                            Add bill
                        </Button>
                    </div>
                </FieldSet>
            </form>
        </div >
    );

}

export default BillForm;
