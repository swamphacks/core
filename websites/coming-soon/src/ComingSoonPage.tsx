import { useState } from "react";
import backgroundImage from "./assets/background.png";
import gatorGif from "./assets/gator.gif";
import signupLabel from "./assets/sign up for updates.png";
import sponsorLabel from "./assets/sponsor sh xii.png";
import swamphacksIcon from "./assets/shxii-icon-transparent.png";
import "./ComingSoonPage.css";
import Modal from "react-modal";

const customStyles = {
  overlay: {
    backgroundColor: "rgba(0, 0, 0, 0.35)",
  },
  content: {
    overflow: "hidden",
    maxHeight: "310px",
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

export function ComingSoonPage() {
  const [modalIsOpen, setIsOpen] = useState(false);

  function openModal() {
    setIsOpen(true);
  }

  function closeModal() {
    setIsOpen(false);
  }

  const handleSignupForUpdates = () => {
    openModal();
  };

  return (
    <main
      data-theme="japanese-style"
      className="coming-soon"
      style={{ backgroundImage: `url(${backgroundImage})` }}
    >
      <section className="coming-soon__content">
        <div>
          <img src={swamphacksIcon} className="icon-image" />
        </div>
        <header className="coming-soon__hero">
          <h1 className="coming-soon__title">SwampHacks XII</h1>

          <div className="coming-soon__badge" aria-label="Coming soon">
            <span>Coming Soon!</span>
          </div>
        </header>

        <div className="coming-soon__board-wrap">
          <div className="coming-soon__actions">
            <div
              onClick={handleSignupForUpdates}
              className="subscribe-button coming-soon__button coming-soon__button--top"
            >
              <img
                src={signupLabel}
                alt="Subscribe button"
                aria-hidden="true"
              />
            </div>

            <div>
              <a
                href="mailto:sponsors@swamphacks.com"
                className="sponsor-button coming-soon__button coming-soon__button--bottom"
                aria-label={"Sponsor SwampHacks XII"}
              >
                <img
                  src={sponsorLabel}
                  alt="Sponsor button"
                  aria-hidden="true"
                />
              </a>
            </div>
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
          className="coming-soon__gator"
        />
      </section>

      <Modal
        isOpen={modalIsOpen}
        onRequestClose={closeModal}
        style={customStyles}
        contentLabel="Example Modal"
      >
        <div className="subscribe-container">
          <button
            onClick={closeModal}
            className="modal-close-btn nes-btn is-error"
          >
            X
          </button>
          <p className="subscribe-title">Stay up to date!</p>
          <div className="subscribe-input-container">
            <input placeholder="Enter your email" className="nes-input" />
            <button type="button" className="nes-btn">
              Subscribe
            </button>
          </div>
        </div>
      </Modal>
    </main>
  );
}
