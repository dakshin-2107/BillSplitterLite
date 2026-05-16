import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import './PeoplePanel.css';

const getInitials = (name: string): string => {
    const words = name.trim().split(/\s+/);
    if (words.length === 1) return words[0].slice(0, 2).toUpperCase();
    return (words[0][0] + words[1][0]).toUpperCase();
};

interface PeoplePanelProps {
    people: string[];
    setPeople: (people: string[]) => void;
}

const PeoplePanel = ({ people, setPeople }: PeoplePanelProps) => {
    const [name, setName] = useState('');
    const [error, setError] = useState<string | null>(null);

    const validateName = (value: string) => {
        const trimmed = value.trim();
        if (trimmed === '') return false;
        if (people.some(p => p.toLowerCase() === trimmed.toLowerCase())) {
            setError('Name already added');
            return false;
        }
        setError(null);
        return true;
    };

    const handleAdd = () => {
        if (validateName(name)) {
            setPeople([...people, name.trim()]);
            setName('');
        }
    };

    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter') handleAdd();
    };

    const handleRemove = (index: number) => {
        setPeople(people.filter((_, i) => i !== index));
    };

    return (
        <div className="people-panel">
            <div className="panel-header">
                <div>
                    <span className="step-label">Step 2</span>
                    <h2 className="panel-title">Add the split participants</h2>
                </div>
                {people.length > 0 && (
                    <span className="panel-count-badge">{people.length} added</span>
                )}
            </div>

            <div className="people-list-container">
                {people.length > 0 ? (
                    people.map((person, index) => (
                        <div key={index} className="person-item">
                            <div className="person-item__left">
                                <div className="person-avatar">{getInitials(person)}</div>
                                <span className="person-name">{person}</span>
                            </div>
                            <Button
                                variant="destructive"
                                size="sm"
                                onClick={() => handleRemove(index)}
                                aria-label={`Remove ${person}`}
                            >
                                ✕
                            </Button>
                        </div>
                    ))
                ) : (
                    <p className="people-panel__empty">No participants yet</p>
                )}
            </div>

            <div className="people-add-section">
                {error && <p className="people-panel__error">{error}</p>}
                <div className="people-input-group">
                    <Input
                        type="text"
                        value={name}
                        placeholder="Add participant name..."
                        onChange={(e) => {
                            setName(e.target.value);
                            if (error) setError(null);
                        }}
                        onKeyDown={handleKeyDown}
                    />
                    <Button onClick={handleAdd} disabled={!name.trim()}>Add</Button>
                </div>
            </div>
        </div>
    );
};

export default PeoplePanel;
