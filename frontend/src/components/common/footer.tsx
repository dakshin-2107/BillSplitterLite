import './footer.css';

const Footer = () => {
    const aboutUri = import.meta.env.VITE_ABOUT_ME_URI as string | undefined;

    return (
        <footer className="app-footer">
            {aboutUri && (
                <a href={aboutUri} className="footer-link" target="_blank" rel="noreferrer">
                    About me
                </a>
            )}
        </footer>
    );
};

export default Footer;
