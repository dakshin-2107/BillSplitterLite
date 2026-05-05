import { useState, useRef, type ChangeEvent } from 'react';
import { FieldLabel, FieldLegend, FieldSet, Field } from '@/components/ui/field';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { DatePicker } from '@/components/ui/datepicker';
import type { BillFormData } from '../../../common/interfaces';
import './BillForm.css';

interface BillFormProps {
    onBillAdded: (bill: BillFormData) => void;
}

const BillForm = ({ onBillAdded }: BillFormProps) => {
    const [location, setLocation] = useState('');
    const [date, setDate] = useState('');
    const [image, setImage] = useState<File | null>(null);
    const [imagePreview, setImagePreview] = useState<string | null>(null);
    const fileInputRef = useRef<HTMLInputElement>(null);

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
    };

    const handleSubmit = () => {
        if (location && date && image) {
            onBillAdded({ location, date, image });
            reset();
        }
    };

    return (
        <div className='bill-form'>
            <form>
                <FieldSet>
                    <FieldLegend>Bill Details</FieldLegend>
                    <Field>
                        <FieldLabel htmlFor="bill-location">Location</FieldLabel>
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
                        <FieldLabel htmlFor="bill-date">Date</FieldLabel>
                        <DatePicker required value={date} onDateChange={setDate} />
                    </Field>
                    <Field className="image-field">
                        <FieldLabel htmlFor="bill-image">Bill image</FieldLabel>
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
                </FieldSet>
            </form>
        </div>
    );
};

export default BillForm;
