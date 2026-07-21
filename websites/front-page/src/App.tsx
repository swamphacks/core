import "./App.css";
import "./Landing.css";
import "./fireflies.sass";
import swamphacksIcon from "./assets/shxii-icon-transparent.png";
import Modal from "react-modal";
import { useEffect, useState } from "react";
import About from "./About";
import Sponsors from "./Sponsors";
import Tracks from "./Tracks";
// import StudentOrgs from "./StudentOrgs";
import Team from "./Team";
import FAQ from "./Faq";
import Insta from "./assets/insta.svg";
import Discord from "./assets/discord.svg";
// import Github from "./assets/github.svg";

Modal.setAppElement("#root");

function App() {
  const [showBadge, setShowBadge] = useState(true);

  useEffect(() => {
    const handleScroll = () => {
      setShowBadge(window.scrollY < 600);
    };

    handleScroll();
    window.addEventListener("scroll", handleScroll, { passive: true });

    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  return (
    <>
      <a
        id="mlh-trust-badge"
        style={{
          display: "block",
          maxWidth: "100px",
          minWidth: "60px",
          position: "fixed",
          right: "20px",
          top: 0,
          width: "10%",
          zIndex: 10000,
          opacity: showBadge ? 1 : 0,
          visibility: showBadge ? "visible" : "hidden",
          transition: "opacity 0.2s ease, visibility 0.2s ease",
          pointerEvents: showBadge ? "auto" : "none",
        }}
        href="https://mlh.io/na?utm_source=na-hackathon&utm_medium=TrustBadge&utm_campaign=2026-season&utm_content=white"
        target="_blank"
        rel="noopener noreferrer"
      >
        <img
          src="https://logged-assets.s3.amazonaws.com/trust-badge/2027/mlh-trust-badge-2027-white.svg"
          alt="Major League Hacking 2026 Hackathon Season"
          style={{ width: "100%" }}
        />
      </a>

      <Landing />
      <About />
      <Tracks />
      <Sponsors />
      {/* <StudentOrgs /> */}
      <Team />
      <FAQ />

      <footer>
        <p>Made with ❤️ by the SwampHacks Team © 2026</p>
        <a
          href="https://github.com/MLH/mlh-policies/blob/main/code-of-conduct.md"
          target="_blank"
          className="code-of-conduct"
        >
          MLH Code of Conduct
        </a>
      </footer>
    </>
  );
}

function Landing() {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  return (
    <div className="container">
      <div className="hamburger-menu" onClick={() => setIsMobileMenuOpen(true)}>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="30"
          height="30"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          className="lucide lucide-menu-icon lucide-menu"
        >
          <path d="M4 5h16" />
          <path d="M4 12h16" />
          <path d="M4 19h16" />
        </svg>
      </div>
      <div
        className={`hamburger-menu-content-container ${isMobileMenuOpen ? "open" : ""}`}
      >
        <button
          className="menu-close-btn nes-btn is-error"
          onClick={() => setIsMobileMenuOpen(false)}
        >
          X
        </button>
        <div className="hamburger-menu-content">
          <div>
            <img src={swamphacksIcon} className="menu-icon-image" />
          </div>
          <a
            className="hamburger-menu-link nes-btn is-primary"
            href="#about"
            onClick={() => setIsMobileMenuOpen(false)}
          >
            About
          </a>
          <a
            className="hamburger-menu-link nes-btn is-primary"
            href="#tracks"
            onClick={() => setIsMobileMenuOpen(false)}
          >
            Tracks
          </a>
          <a
            className="hamburger-menu-link nes-btn is-primary"
            href="#sponsors"
            onClick={() => setIsMobileMenuOpen(false)}
          >
            Sponsors
          </a>
          <a
            className="hamburger-menu-link nes-btn is-primary"
            href="#team"
            onClick={() => setIsMobileMenuOpen(false)}
          >
            Team
          </a>
          <a
            className="hamburger-menu-link nes-btn is-primary"
            href="#faq"
            onClick={() => setIsMobileMenuOpen(false)}
          >
            FAQ
          </a>
        </div>
      </div>
      <div className="menu">
        <div>
          <img src={swamphacksIcon} className="menu-icon-image" />
        </div>
        <a href="#about">About</a>
        <a href="#tracks">Tracks</a>
        <a href="#sponsors">Sponsors</a>
        <a href="#team">Team</a>
        <a href="#faq">FAQ</a>
      </div>
      {Array.from({ length: 15 }).map((_, i) => (
        <div className="firefly" key={i}></div>
      ))}
      <div className="background">
        <div className="layer layer-bushes-back"></div>
        <div className="layer layer-trees"></div>
        <div className="layer layer-trees-back1"></div>
        <div className="layer layer-trees-back2"></div>
        <div className="layer layer-trees-back3"></div>
        <div className="layer layer-foreground"></div>
        <div className="layer layer-vegetation"></div>
        <div className="layer layer-canopy"></div>
        <div className="layer layer-canopy-back"></div>
      </div>
      <main>
        <header className="hero">
          <div>
            <img src={swamphacksIcon} className="icon-image" />
          </div>
          <h1 className="title">SwampHacks XII</h1>
          <div className="subheader">
            <p className="date">Oct 16 - 18, 2026</p>
            <p className="location">Reitz Union</p>
          </div>
        </header>

        <div className="menu-container">
          <a
            className="nes-btn is-primary register-button"
            href="https://app.swamphacks.com/application"
            target="_blank"
          >
            Apply
          </a>

          <a
            className="nes-btn is-primary sponsor-button"
            href="mailto:sponsors@swamphacks.com"
            target="_blank"
          >
            Sponsor SH XII
          </a>
        </div>

        <div className="submenu-container">
          <a
            className="submenu-link"
            href="https://swamphack.notion.site/3973b41de22f80b788ced816145e0a2d"
            target="_blank"
          >
            Mentor Application
          </a>
        </div>

        <section className="icon-list">
          <a href="https://www.instagram.com/ufswamphacks/" target="_blank">
            <img src={Insta} />
          </a>

          <a href="https://discord.gg/YBHrXPJ8mR" target="_blank">
            <img src={Discord} />
          </a>
        </section>
      </main>
    </div>
  );
}

export default App;
