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
import Github from "./assets/github.svg";

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
  return (
    <div className="container">
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
        </div>

        <div className="submenu-container">
          <a
            className="submenu-link"
            href="https://swamphack.notion.site/3973b41de22f80b788ced816145e0a2d"
            target="_blank"
          >
            Mentor Application
          </a>
          {/* <p className="submenu-container-title">Support SwampHacks</p>
          <div className="submenu-links-container">
            <a className="submenu-link">Volunteer</a>
            <span>&#9679;</span>
            <a className="submenu-link">Mentor</a>
          </div> */}
        </div>

        <section className="icon-list">
          <a href="https://www.instagram.com/ufswamphacks/" target="_blank">
            <img src={Insta} />
            {/* <i className="nes-icon instagram"></i> */}
          </a>

          <a href="https://discord.gg/YBHrXPJ8mR" target="_blank">
            <img src={Discord} />
            {/* <i className="nes-icon discord"></i> */}
          </a>

          <a href="https://github.com/swamphacks/core" target="_blank">
            <img src={Github} />
            {/* <i className="nes-icon github"></i> */}
          </a>
        </section>

        {/* <img
          src={gatorGif}
          alt="Pixel-art alligator walking in front of the sign"
          className="gator"
        /> */}
      </main>
    </div>
  );
}

// interface InformationModalProps {
//   closeModal: () => void;
// }

// function InformationModal({ closeModal }: InformationModalProps) {
//   const [currentTab, setCurrentTab] = useState("about");

//   const handleSetTab = (tab: string) => {
//     setCurrentTab(tab);
//   };

//   const renderTab = () => {
//     switch (currentTab) {
//       case "about":
//         return <About />;
//       case "tracks":
//         return <Tracks />;
//       case "team":
//         return <Team />;
//       case "faq":
//         return <FAQ />;
//     }
//   };

//   return (
//     <div
//       className={`info-modal-container ${currentTab === "team" ? "modal-large" : ""}`}
//     >
//       <button onClick={closeModal} className="modal-close-btn nes-btn is-error">
//         X
//       </button>
//       <div className="tabs">
//         <button
//           onClick={() => handleSetTab("about")}
//           className="nes-btn"
//           style={{
//             opacity: currentTab === "about" ? "100%" : "50%",
//           }}
//         >
//           About
//         </button>
//         <button
//           onClick={() => handleSetTab("tracks")}
//           className="nes-btn"
//           style={{
//             opacity: currentTab === "tracks" ? "100%" : "50%",
//           }}
//         >
//           Tracks
//         </button>
//         <button
//           onClick={() => handleSetTab("faq")}
//           className="nes-btn"
//           style={{
//             opacity: currentTab === "faq" ? "100%" : "50%",
//           }}
//         >
//           FAQ
//         </button>
//         <button
//           onClick={() => handleSetTab("team")}
//           className="nes-btn"
//           style={{
//             opacity: currentTab === "team" ? "100%" : "50%",
//           }}
//         >
//           Team
//         </button>
//       </div>
//       {renderTab()}
//     </div>
//   );
// }

export default App;
