import "./App.css";
import swamphacksIcon from "./assets/shxii-icon-transparent.png";
import gatorGif from "./assets/gator.gif";
import Modal from "react-modal";
import { useState } from "react";

const customStyles = {
  overlay: {
    backgroundColor: "rgba(0, 0, 0, 0.35)",
  },
  content: {
    overflow: "hidden",
    // maxHeight: "310px",
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

  function openModal() {
    setIsOpen(true);
  }

  function closeModal() {
    setIsOpen(false);
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
          <h1 className="subtitle">Oct 16 - 18, 2026 | Reitz Union</h1>
        </header>

        <div className="menu-container">
          <button className="nes-btn is-primary">Register</button>
          <button className="nes-btn is-primary" onClick={openModal}>
            About
          </button>
          <div className="community-buttons">
            <button className="nes-btn is-primary">Sponsors</button>
            <button className="nes-btn is-primary">Partners</button>
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
        isOpen={modalIsOpen}
        onRequestClose={closeModal}
        style={customStyles}
      >
        <InformationModal closeModal={closeModal} />
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
    <div className="info-modal-container">
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
          onClick={() => handleSetTab("team")}
          className="nes-btn"
          style={{
            opacity: currentTab === "team" ? "100%" : "50%",
          }}
        >
          Team
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
      </div>
      {renderTab()}
    </div>
  );
}

function About() {
  return (
    <div>
      <p className="tab-title">About SwampHacks</p>
      <div className="tab-container">
        <p>
          SwampHacks is the University of Florida’s flagship hackathon, bringing
          together 400+ students each year for 36 hours of creativity,
          collaboration, and innovation. Recognized for excellence by UF’s
          Herbert Wertheim College of Engineering, SwampHacks offers hands-on
          workshops, mentorship, and community-building activities that help
          hackers grow their skills and bring their ideas to life. Whether
          you’re a first-time hacker or a seasoned coder, SwampHacks is the
          place to build, connect, and inspire.
        </p>
      </div>
    </div>
  );
}

function Tracks() {
  const tracks = [
    {
      name: "Overall Prize",
      description:
        "All projects are considered for the Overall Prize. The General Track is for any innovative projects that do not necessarily fit under any of the other tracks.",
    },
    {
      name: "Education, Accessibility & Social Impact",
      description:
        "All projects are considered for the Overall Prize. The General Track is for any innovative projects that do not necessarily fit under any of the other tracks.",
    },
    {
      name: "Sustainability",
      description:
        "All projects are considered for the Overall Prize. The General Track is for any innovative projects that do not necessarily fit under any of the other tracks.",
    },
    {
      name: "Game Design",
      description:
        "All projects are considered for the Overall Prize. The General Track is for any innovative projects that do not necessarily fit under any of the other tracks.",
    },
    {
      name: "The Design of Everyday Life (Human-Centered Design)",
      description:
        "All projects are considered for the Overall Prize. The General Track is for any innovative projects that do not necessarily fit under any of the other tracks.",
    },
  ];
  const [openIndex, setOpenIndex] = useState<number | null>(null);

  const toggle = (index: number) => {
    setOpenIndex(openIndex === index ? null : index);
  };

  // TODO: when click on "team", resize the modal to be taller to fit more content

  return (
    <div>
      <p className="tab-title">Tracks</p>
      <div className="tab-container track-container">
        {tracks.map((track, idx) => (
          <div
            key={track.name}
            className={`track ${openIndex === idx ? "open" : ""}`}
          >
            <button
              className={`track-header ${openIndex === idx ? "open" : ""}`}
              onClick={() => toggle(idx)}
              aria-expanded={openIndex === idx}
            >
              {track.name}
            </button>

            {openIndex === idx && (
              <div className="track-description">
                <p>{track.description}</p>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}

function Team() {
  return <div>Team</div>;
}

function FAQ() {
  return <div>FAQ</div>;
}

export default App;
