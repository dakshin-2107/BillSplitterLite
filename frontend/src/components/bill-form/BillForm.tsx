import React, { useState, type FormEvent, type ChangeEvent, useContext, createContext, useMemo } from 'react';
import './BillForm.css';
import type { BillData, FormData, ApiResponse } from '../../common/interfaces';

interface BillFormProps {
    onSubmitSuccess: (billData: BillData) => void;
    onSubmitError: (error: string) => void;
}


const BillForm: React.FC<BillFormProps> = ({ onSubmitSuccess, onSubmitError }) => {

    const [formData, setFormData] = useState<FormData>({
        place: '',
        dateTime: '',
        names: [''],
        image: null,
    });

    const [currentName, setCurrentName] = useState('');
    const [imagePreview, setImagePreview] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Handle text input changes
    const handleInputChange = (e: ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value,
        }));
    };

    // Add a person to the list
    const addPerson = () => {
        if (currentName.trim()) {
            setFormData(prev => ({
                ...prev,
                names: [...prev.names.filter(n => n.trim() !== ''), currentName.trim(), ''].filter((n, i, arr) => i < arr.length - 1 || n === ''),
            }));
            setCurrentName('');
        }
    };

    // Handle Enter key in person input
    const handlePersonKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            addPerson();
        }
    };

    // Remove a person from the list
    const removePerson = (index: number) => {
        setFormData(prev => ({
            ...prev,
            names: prev.names.filter((_, i) => i !== index),
        }));
    };

    // Handle image upload
    const handleImageChange = (e: ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            // Validate file type
            if (!file.type.startsWith('image/')) {
                setError('Please select a valid image file');
                return;
            }

            // Validate file size (max 5MB)
            if (file.size > 5 * 1024 * 1024) {
                setError('Image size should be less than 5MB');
                return;
            }

            setFormData(prev => ({
                ...prev,
                image: file,
            }));

            // Create preview
            const reader = new FileReader();
            reader.onloadend = () => {
                setImagePreview(reader.result as string);
            };
            reader.readAsDataURL(file);
            setError(null);
        }
    };

    // Validate form
    const validateForm = (): boolean => {
        if (!formData.place.trim()) {
            setError('Please enter a location');
            return false;
        }

        if (!formData.dateTime) {
            setError('Please select a date');
            return false;
        }

        const validPeople = formData.names.filter(person => person.trim() !== '');
        if (validPeople.length === 0) {
            setError('Please add at least one person');
            return false;
        }

        if (!formData.image) {
            setError('Please upload an image');
            return false;
        }

        return true;
    };

    // Handle form submission
    const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setError(null);

        if (!validateForm()) {
            return;
        }

        setIsSubmitting(true);

        try {
            // Create FormData for multipart/form-data submission
            const submitData = new FormData();
            submitData.append('place', formData.place);
            submitData.append('dateTime', formData.dateTime);

            // Filter out empty names and add each person
            const validPeople = formData.names.filter(person => person.trim() !== '');
            validPeople.forEach((person) => {
                submitData.append('names', person);
            });

            if (formData.image) {
                submitData.append('image', formData.image);
            }

            // Get API URL from environment variable
            const apiUrl = import.meta.env.VITE_BACKEND_URL;
            //const apiUrl = "http://localhost:8080/home";
            if (!apiUrl) {
                throw new Error('API URL not configured. Please check your .env file.');
            }

            // Log what we're sending to the endpoint
            console.log('📤 Sending form data to:', apiUrl);
            console.log('📦 Form data contents:', submitData);

            // Submit form data
            const response = await fetch(apiUrl, {
                method: 'POST',
                body: submitData,
                credentials: 'include',
            });

            if (!response.ok) {
                throw new Error(`Server error: ${response.statusText}`);
            }

            const result: ApiResponse = await response.json();
            console.log('📥 Server response:', result);

            if (result.success && result.bill) {
                // Reset form on success
                setFormData({
                    place: '',
                    dateTime: '',
                    names: [''],
                    image: null,
                });
                setCurrentName('');
                setImagePreview(null);
                setError(null);

                // Call success callback
                onSubmitSuccess(result.bill);
            } else {
                throw new Error(result.message || 'Failed to process bill');
            }

        } catch (err) {
            const errorMessage = err instanceof Error ? err.message : 'Failed to submit form';
            setError(errorMessage);
            onSubmitError(errorMessage);
        } finally {
            setIsSubmitting(false);
        }
    };

    const validPeople = formData.names.filter(person => person.trim() !== '');

    return (
        <div className="bill-form-container">
            <div className="form-wrapper">
                <form onSubmit={handleSubmit} className="event-form">
                    {/* Left Column */}
                    <div className="form-left">
                        {/* Location Input */}
                        <div className="form-field">
                            <label htmlFor="place" className="form-label">
                                Location :
                            </label>
                            <input
                                type="text"
                                id="place"
                                name="place"
                                value={formData.place}
                                onChange={handleInputChange}
                                className="form-input"
                                placeholder="Dunder mifflin"
                                disabled={isSubmitting}
                            />
                        </div>

                        {/* Date Input */}
                        <div className="form-field">
                            <label htmlFor="dateTime" className="form-label">
                                Date :
                            </label>
                            <input
                                type="date"
                                id="dateTime"
                                name="dateTime"
                                value={formData.dateTime}
                                onChange={handleInputChange}
                                className="form-input"
                                disabled={isSubmitting}
                            />
                        </div>

                        {/* People Section */}
                        <div className="people-section">
                            <label className="form-label">People involved :</label>

                            <div className="people-input-group">
                                <input
                                    type="text"
                                    value={currentName}
                                    onChange={(e) => setCurrentName(e.target.value)}
                                    onKeyPress={handlePersonKeyPress}
                                    className="form-input"
                                    placeholder="Enter a name"
                                    disabled={isSubmitting}
                                />
                                <button
                                    type="button"
                                    onClick={addPerson}
                                    className="btn-add"
                                    disabled={isSubmitting || !currentName.trim()}
                                >
                                    Add
                                </button>
                            </div>

                            {/* People List */}
                            {validPeople.length > 0 && (
                                <div className="people-list">
                                    {validPeople.map((person, index) => (
                                        <div key={index} className="person-item">
                                            <span className="person-name">{person}</span>
                                            <button
                                                type="button"
                                                onClick={() => removePerson(formData.names.indexOf(person))}
                                                className="btn-remove-person"
                                                disabled={isSubmitting}
                                                aria-label="Remove person"
                                            >
                                                −
                                            </button>
                                        </div>
                                    ))}
                                </div>
                            )}
                        </div>
                    </div>

                    {/* Right Column */}
                    <div className="form-right">
                        <div className="upload-section">
                            <label className="form-label upload-label">Upload bill image :</label>

                            <input
                                type="file"
                                id="image"
                                accept="image/*"
                                onChange={handleImageChange}
                                className="file-input"
                                disabled={isSubmitting}
                            />

                            <label htmlFor="image" className="upload-area">
                                {imagePreview ? (
                                    <div className="image-preview">
                                        <img src={imagePreview} alt="Bill preview" />
                                        <div className="overlay">
                                            <span>Click to change image</span>
                                        </div>
                                    </div>
                                ) : (
                                    <div className="upload-placeholder">
                                        <div className="plus-icon">+</div>
                                        <p className="upload-instruction">
                                            Click to upload<br />the image of a bill
                                        </p>
                                    </div>
                                )}
                            </label>

                            {/* Submit Button */}
                            <button
                                type="submit"
                                className="btn-submit"
                                disabled={isSubmitting}
                            >
                                {isSubmitting ? 'Creating session...' : 'Create new session'}
                            </button>

                            <p className="submit-hint">
                                Click this button to submit the form and create a new split up session
                            </p>
                        </div>
                    </div>

                    {/* Error Message */}
                    {error && (
                        <div className="error-message">
                            {error}
                        </div>
                    )}
                </form>
            </div>
        </div>
    );
};

export default BillForm;
