// Formats a YYYY-MM-DD date string as "19th Apr, 2026"
export const formatDate = (dateStr: string): string => {
    const [year, month, day] = dateStr.split('-').map(Number);
    const suffix =
        day === 1 || day === 21 || day === 31 ? 'st' :
        day === 2 || day === 22 ? 'nd' :
        day === 3 || day === 23 ? 'rd' : 'th';
    const monthName = new Date(year, month - 1, 1).toLocaleString('en-US', { month: 'short' });
    return `${day}${suffix} ${monthName}, ${year}`;
};
