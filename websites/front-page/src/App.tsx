import "./App.css";
import "./Landing.css";
import "./fireflies.sass";
import swamphacksIcon from "./assets/shxii-icon-transparent.png";
import Modal from "react-modal";
import {
  useEffect,
  useRef,
  useState,
  type PointerEvent,
} from "react";
import About from "./About";
import Sponsors from "./Sponsors";
import Tracks from "./Tracks";
import StudentOrgs from "./StudentOrgs";
import Team from "./Team";
import FAQ from "./Faq";
import Insta from "./assets/insta.svg";
import Discord from "./assets/discord.svg";
// import Github from "./assets/github.svg";
import CountdownCorner from "./Countdown";
import "./Buttons.css";

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
          transform: showBadge
            ? "translateY(0) scaleY(1)"
            : "translateY(-8px) scaleY(0)",
          transformOrigin: "top center",
          transition:
            "transform 0.6s cubic-bezier(0.65, 0, 0.35, 1), opacity 0.3s ease, visibility 0.6s",
          willChange: "transform, opacity",
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

      <CountdownCorner
        targetDate={new Date("2026-10-16")}
        label="Swamphacks Begins In:"
      />

      <Landing />
      <ButterflyTrail />
      <About />
      <Tracks />
      <Sponsors />
      <StudentOrgs />
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

function ButterflyTrail() {
  const trailRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    let animationFrame = 0;

    const horizontalPositions = [0.12, 0.31, 0.68, 0.84, 0.49];
    const verticalPositions = [0.12, 0.34, 0.55, 0.72, 0.88];

    const updateRegion = () => {
      const trail = trailRef.current;
      const about = document.getElementById("about");
      const sponsors = document.getElementById("sponsors");

      if (!trail || !about || !sponsors) return;

      const top = about.offsetTop;
      const height = Math.max(sponsors.offsetTop - top, 1);

      trail.style.top = `${top}px`;
      trail.style.height = `${height}px`;
    };

    const animate = (time: number) => {
      const trail = trailRef.current;

      if (trail) {
        const butterflies =
          trail.querySelectorAll<HTMLElement>(".butterfly");
        const width = trail.clientWidth;
        const height = trail.clientHeight;

        butterflies.forEach((butterfly, index) => {
          const horizontalPosition =
            width * horizontalPositions[index] +
            Math.sin(time / (980 + index * 130) + index * 1.8) * 34 +
            Math.cos(time / (570 + index * 80) + index) * 8;

          const verticalPosition =
            height * verticalPositions[index] +
            Math.sin(time / (1250 + index * 140) + index) * 42 +
            Math.cos(time / (760 + index * 90) + index * 1.4) * 15;

          const angle =
            Math.sin(time / 700 + index * 1.4) * 8;

          butterfly.style.transform =
            `translate3d(${horizontalPosition}px, ` +
            `${verticalPosition}px, 0) rotate(${angle}deg)`;
        });
      }

      animationFrame = window.requestAnimationFrame(animate);
    };

    updateRegion();
    window.addEventListener("resize", updateRegion);
    animationFrame = window.requestAnimationFrame(animate);

    return () => {
      window.removeEventListener("resize", updateRegion);
      window.cancelAnimationFrame(animationFrame);
    };
  }, []);

  return (
    <div
      ref={trailRef}
      className="butterfly-trail"
      aria-hidden="true"
    >
      {Array.from({ length: 5 }).map((_, index) => (
        <span className="butterfly" key={index}>
          <span className="butterfly-body" />
        </span>
      ))}
    </div>
  );
}

function Landing() {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  const handleParallax = (event: PointerEvent<HTMLDivElement>) => {
    if (event.pointerType === "touch") return;

    const bounds = event.currentTarget.getBoundingClientRect();
    const x = (event.clientX - bounds.left) / bounds.width - 0.5;
    const y = (event.clientY - bounds.top) / bounds.height - 0.5;

    event.currentTarget.style.setProperty("--scene-x", `${x * 26}px`);
    event.currentTarget.style.setProperty("--scene-y", `${y * 16}px`);
  };

  const resetParallax = (event: PointerEvent<HTMLDivElement>) => {
    event.currentTarget.style.setProperty("--scene-x", "0px");
    event.currentTarget.style.setProperty("--scene-y", "0px");
  };

  return (
    <div
      className="container"
      onPointerMove={handleParallax}
      onPointerLeave={resetParallax}
    >
      <div className="hamburger-menu" onClick={() => setIsMobileMenuOpen(true)}>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="30"
          height="30"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
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
          className="menu-close-btn pixel-button pixel-button--danger"
          onClick={() => setIsMobileMenuOpen(false)}
        >
          X
        </button>
        <div className="hamburger-menu-content">
          <div>
            <img src={swamphacksIcon} className="menu-icon-image" />
          </div>
          <a
            className="hamburger-menu-link pixel-button"
            href="#about"
            onClick={() => setIsMobileMenuOpen(false)}
          >
            About
          </a>
          <a
            className="hamburger-menu-link pixel-button"
            href="#tracks"
            onClick={() => setIsMobileMenuOpen(false)}
          >
            Tracks
          </a>
          <a
            className="hamburger-menu-link pixel-button"
            href="#sponsors"
            onClick={() => setIsMobileMenuOpen(false)}
          >
            Sponsors
          </a>
          <a
            className="hamburger-menu-link pixel-button"
            href="#team"
            onClick={() => setIsMobileMenuOpen(false)}
          >
            Team
          </a>
          <a
            className="hamburger-menu-link pixel-button"
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
          <h1 className="title" tabIndex={0}>SwampHacks XII</h1>
          <div className="subheader">
            <p className="date">Oct 16 - 18, 2026</p>
            <p className="location">Reitz Union</p>
          </div>
        </header>

        <div className="menu-container">
          <a
            className="pixel-button register-button"
            href="https://app.swamphacks.com/application"
            target="_blank"
          >
            Apply
          </a>

          <a
            className="pixel-button sponsor-button"
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

          <a
            className="submenu-link"
            href="https://swamphack.notion.site/54a3b41de22f8324afa9814483091664"
            target="_blank"
          >
            Judges Interest Form
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
