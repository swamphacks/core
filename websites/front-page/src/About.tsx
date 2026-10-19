import { useState } from "react";
import "./About.css";
import Modal from "react-modal";
import Sign from "./assets/Sign.png";
import Sign2 from "./assets/Sign3.png";
import Pic1 from "./assets/random/pic1.jpg";
import Pic2 from "./assets/random/pic2.jpg";
import Pic3 from "./assets/random/pic3.jpg";
import Pic4 from "./assets/random/pic4.jpg";
import Pic5 from "./assets/random/pic5.jpg";
import Pic6 from "./assets/random/pic6.jpg";
import Camera from "./assets/Camera.png";

// Test
const images = [Pic1, Pic2, Pic3, Pic4, Pic5, Pic6];

const customModalStyles = {
  overlay: {
    backgroundColor: "rgba(3, 7, 2, 0.62)",
    zIndex: "10000",
  },
  content: {
    top: "50%",
    left: "50%",
    right: "auto",
    bottom: "auto",
    width: "min(calc(100% - 32px), 540px)",
    maxHeight: "82vh",
    marginRight: "-50%",
    padding: "28px",
    overflow: "auto",
    color: "#f4e4cf",
    backgroundColor: "#4d3222",
    border: "3px solid #b98762",
    borderRadius: "4px",
    boxShadow:
      "8px 8px 0 #1d130d, inset 0 0 0 2px #2a1b12",
    outline: "none",
    transform: "translate(-50%, -50%)",
    zIndex: "10001",
  },
};

export default function About() {
  const [selectedIndex, setSelectedIndex] = useState<number | null>(null);

  return (
    <div id="about" className="about-container">
      <div className="about-header">
        <h1>About Us</h1>
        <div className="about-summary-container">
          <p className="summary">
            SwampHacks is the University of Florida’s flagship hackathon,
            bringing together 400+ students each year for 36 hours of
            creativity, collaboration, and innovation. Recognized for excellence
            by UF’s Herbert Wertheim College of Engineering, SwampHacks offers
            hands-on workshops, mentorship, and community-building activities
            that help hackers grow their skills and bring their ideas to life.
            Whether you’re a first-time hacker or a seasoned coder, SwampHacks
            is the place to build, connect, and inspire.
          </p>
        </div>
        <p className="sh-xi">A look back at SwampHacks XI</p>
        <div className="pictures-container">
          <img className="camera" src={Camera} alt="" />
          <div className="carousel-window">
            <div className="carousel-track">
              {[0, 1].map((copy) => (
                <div
                  className="carousel-set"
                  key={copy}
                  aria-hidden={copy === 1}
                >
                  {images.map((image, index) => (
                    <img
                      onClick={() => setSelectedIndex(index)}
                      key={`${copy}-${image}-${index}`}
                      className="hackathon-picture"
                      src={image}
                      alt={
                        copy === 0
                          ? `SwampHacks XI moment ${index + 1}`
                          : ""
                      }
                    />
                  ))}
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      <div className="stats">
        <img className="stats-sign" src={Sign} alt="" />
        <span
          className="stats-badge"
          style={{
            transform: "rotate(2deg)",
          }}
        >
          <img className="stats-background" src={Sign2} alt="" />
          <span className="is-primary stats-text">400+ hackers</span>
        </span>
        <span className="stats-badge">
          <img className="stats-background" src={Sign2} alt="" />
          <span className="stats-text">10+ workshops</span>
        </span>
        <span
          className="stats-badge"
          style={{
            transform: "rotate(-2deg)",
          }}
        >
          <img className="stats-background" src={Sign2} alt="" />
          <span className="stats-text">100+ projects</span>
        </span>
        <span className="stats-badge">
          <img className="stats-background" src={Sign2} alt="" />
          <span className="stats-text">$10k+ in prizes</span>
        </span>
      </div>

      <Modal
        isOpen={selectedIndex !== null}
        onRequestClose={() => setSelectedIndex(null)}
        style={customModalStyles}
      >
        <div className="tracks-modal-container">
          <button
            onClick={() => setSelectedIndex(null)}
            className="modal-close-btn pixel-button pixel-button--danger"
          >
            X
          </button>
          <img className="image-modal" src={images[selectedIndex!]} />
        </div>
      </Modal>
    </div>
  );
}

