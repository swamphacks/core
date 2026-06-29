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
  const [email, setEmail] = useState("");
  const [subscribeSuccessful, setSubscribeSuccessful] = useState(false);
  const [modalIsOpen, setIsOpen] = useState(false);

  function openModal() {
    setIsOpen(true);
  }

  function closeModal() {
    setIsOpen(false);
    setSubscribeSuccessful(false);
  }

  const handleSignupForUpdates = () => {
    openModal();
  };

  const handleSubscribe = async () => {
    if (email.length === 0 || email.length > 200) return;

    let API_URL_PREFIX: string;

    if (import.meta.env.MODE === "development") {
      API_URL_PREFIX = "http://localhost:8080";
    } else {
      API_URL_PREFIX = "https://api.swamphacks.com";
    }

    const res = await fetch(`${API_URL_PREFIX}/hackathon/interest`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        email,
        source: "coming-soon-page",
      }),
    });

    if (res.ok) {
      setSubscribeSuccessful(true);
    }
  };

  return (
    <main
      data-theme="japanese-style"
      className="coming-soon"
      style={{ backgroundImage: `url(${backgroundImage})` }}
    >
      <a
        id="mlh-trust-badge"
        style={{
          display: "block",
          maxWidth: "100px",
          minWidth: "60px",
          position: "fixed",
          right: "50px",
          top: 0,
          width: "10%",
          zIndex: 10000,
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
      <section className="coming-soon__content">
        <div>
          <img src={swamphacksIcon} className="icon-image" />
        </div>
        <header className="coming-soon__hero">
          <h1 className="coming-soon__title">
            SwampHacks XII
            <span className="coming-soon__badge">Coming Soon!</span>
          </h1>
          <h1 className="coming-soon__date">Oct 16 - 18, 2026</h1>
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

        <div className="codeofconduct">
          <a
            target="_blank"
            href="https://github.com/MLH/mlh-policies/blob/main/code-of-conduct.md"
          >
            MLH Code Of Conduct
          </a>
        </div>

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
            <input
              placeholder="Enter your email"
              className="nes-input"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
            <button onClick={handleSubscribe} type="button" className="nes-btn">
              Subscribe
            </button>
          </div>
          {subscribeSuccessful && (
            <p className="success-message">Subscribed successfully!</p>
          )}
        </div>
      </Modal>
    </main>
  );
}
