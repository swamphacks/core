import "./App.css";
import swamphacksIcon from "./assets/shxii-icon-transparent.png";
import gatorGif from "./assets/gator.gif";
import Modal from "react-modal";
import { useState } from "react";
import Team from "./Team";
import { Tracks } from "./Tracks";
import FAQ from "./Faq";
import About from "./About";
import Sponsors from "./Sponsors";
import Partners from "./Partners";

const customModalStyles = {
  overlay: {
    backgroundColor: "rgba(0, 0, 0, 0.35)",
  },
  content: {
    overflow: "hidden",
    borderColor: "#231909",
    borderWidth: "5px",
    backgroundColor: "#7D573C",
    top: "50%",
    left: "50%",
    right: "auto",
    bottom: "auto",
    marginRight: "-50%",
    transform: "translate(-50%, -50%)",
  },
};

Modal.setAppElement("#root");

function App() {
  const [modalIsOpen, setIsOpen] = useState(false);
  const [sponsorModalIsOpen, setSponsorModalIsOpen] = useState(false);
  const [partnersModalIsOpen, setPartnersModalIsOpen] = useState(false);

  function closeModal() {
    setIsOpen(false);
    setSponsorModalIsOpen(false);
    setPartnersModalIsOpen(false);
  }

  return (
    <div className="container">
      <main>
        <header className="hero">
          <div>
            <img src={swamphacksIcon} className="icon-image" />
          </div>
          <h1 className="title">
            SwampHacks XII
            <span className="badge">Apply now!</span>
          </h1>

          <div className="subheader">
            <p className="date">Oct 16 - 18, 2026</p>
            <p className="location">Reitz Union</p>
          </div>
        </header>

        <div className="menu-container">
          <button className="nes-btn is-primary">Register</button>
          <button
            className="nes-btn is-primary"
            onClick={() => setIsOpen(true)}
          >
            About
          </button>
          <div className="community-buttons">
            <button
              className="nes-btn is-primary"
              onClick={() => setSponsorModalIsOpen(true)}
            >
              Sponsors
            </button>
            <button
              className="nes-btn is-primary"
              onClick={() => setPartnersModalIsOpen(true)}
            >
              Partners
            </button>
          </div>
        </div>

        <div className="submenu-container">
          <p className="submenu-container-title">Support SwampHacks</p>
          <div className="submenu-links-container">
            <a className="submenu-link">Judge</a>
            <span>&#9679;</span>
            <a className="submenu-link">Volunteer</a>
            <span>&#9679;</span>
            <a className="submenu-link">Mentor</a>
          </div>
        </div>

        <section className="icon-list">
          <a href="https://www.instagram.com/ufswamphacks/" target="_blank">
            <i className="nes-icon instagram"></i>
          </a>

          <a href="https://discord.gg/YBHrXPJ8mR" target="_blank">
            <i className="nes-icon discord"></i>
          </a>

          <a href="https://github.com/swamphacks/core" target="_blank">
            <i className="nes-icon github"></i>
          </a>
        </section>

        <img
          src={gatorGif}
          alt="Pixel-art alligator walking in front of the sign"
          className="gator"
        />
      </main>

      <footer>
        <p>Made with ❤️ by the SwampHacks Team © 2026</p>
      </footer>

      <Modal
        isOpen={modalIsOpen || sponsorModalIsOpen || partnersModalIsOpen}
        onRequestClose={closeModal}
        style={customModalStyles}
      >
        {sponsorModalIsOpen ? (
          <Sponsors closeModal={closeModal} />
        ) : partnersModalIsOpen ? (
          <Partners closeModal={closeModal} />
        ) : (
          <InformationModal closeModal={closeModal} />
        )}
      </Modal>
    </div>
  );
}

interface InformationModalProps {
  closeModal: () => void;
}

function InformationModal({ closeModal }: InformationModalProps) {
  const [currentTab, setCurrentTab] = useState("about");

  const handleSetTab = (tab: string) => {
    setCurrentTab(tab);
  };

  const renderTab = () => {
    switch (currentTab) {
      case "about":
        return <About />;
      case "tracks":
        return <Tracks />;
      case "team":
        return <Team />;
      case "faq":
        return <FAQ />;
    }
  };

  return (
    <div
      className={`info-modal-container ${currentTab === "team" ? "modal-large" : ""}`}
    >
      <button onClick={closeModal} className="modal-close-btn nes-btn is-error">
        X
      </button>
      <div className="tabs">
        <button
          onClick={() => handleSetTab("about")}
          className="nes-btn"
          style={{
            opacity: currentTab === "about" ? "100%" : "50%",
          }}
        >
          About
        </button>
        <button
          onClick={() => handleSetTab("tracks")}
          className="nes-btn"
          style={{
            opacity: currentTab === "tracks" ? "100%" : "50%",
          }}
        >
          Tracks
        </button>
        <button
          onClick={() => handleSetTab("faq")}
          className="nes-btn"
          style={{
            opacity: currentTab === "faq" ? "100%" : "50%",
          }}
        >
          FAQ
        </button>
        <button
          onClick={() => handleSetTab("team")}
          className="nes-btn"
          style={{
            opacity: currentTab === "team" ? "100%" : "50%",
          }}
        >
          Team
        </button>
      </div>
      {renderTab()}
    </div>
  );
}

export default App;
