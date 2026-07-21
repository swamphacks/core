import { useEffect, useState } from "react";
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
import Arrow from "./assets/arrow.png";
import Camera from "./assets/Camera.png";

// Test
const images = [Pic1, Pic2, Pic3, Pic4, Pic5, Pic6];

const customModalStyles = {
  overlay: {
    backgroundColor: "rgba(0, 0, 0, 0.35)",
    zIndex: "100",
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
    zIndex: "100",
  },
};

export default function About() {
  const [activeIndex, setActiveIndex] = useState(0);
  const [selectedIndex, setSelectedIndex] = useState<number | null>(null);
  const isMobile = useIsMobile();
  const itemWidth = isMobile ? 183 : 210;

  const handlePrevious = () => {
    if (activeIndex == 0) {
      return;
    }

    setActiveIndex(
      (currentIndex) => (currentIndex - 1 + images.length) % images.length,
    );
  };

  const handleNext = () => {
    if (isMobile && activeIndex == images.length / 2 + 1) {
      return;
    }
    if (!isMobile && activeIndex == images.length / 2) {
      return;
    }
    setActiveIndex((currentIndex) => (currentIndex + 1) % images.length);
  };

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
          <button
            type="button"
            className="arrow-button"
            onClick={handlePrevious}
            aria-label="Show previous images"
          >
            <img className="arrow" src={Arrow} alt="Previous" />
          </button>

          <div className="carousel-window">
            <div
              className="carousel-track"
              style={{
                transform: `translateX(-${activeIndex * itemWidth}px)`,
              }}
            >
              {images.map((image, index) => (
                <img
                  onClick={() => setSelectedIndex(index)}
                  key={`${image}-${index}`}
                  className="hackathon-picture"
                  src={image}
                  alt={`SwampHacks XI moment ${index + 1}`}
                />
              ))}
            </div>
          </div>

          <button
            type="button"
            className="arrow-button"
            onClick={handleNext}
            aria-label="Show next images"
          >
            <img className="arrow right" src={Arrow} alt="Next" />
          </button>
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
            className="modal-close-btn nes-btn is-error"
          >
            X
          </button>
          <img className="image-modal" src={images[selectedIndex!]} />
        </div>
      </Modal>
    </div>
  );
}

function useIsMobile(breakpoint: number = 768): boolean {
  // Initialize state; default to false if window is not available (SSR safe)
  const [isMobile, setIsMobile] = useState<boolean>(false);

  useEffect(() => {
    // Return early if code runs on the server side
    if (typeof window === "undefined") return;

    const mediaQuery = window.matchMedia(`(max-width: ${breakpoint}px)`);

    // Set initial value
    setIsMobile(mediaQuery.matches);

    // Listener function to catch screen size updates
    const handleMediaQueryChange = (event: MediaQueryListEvent) => {
      setIsMobile(event.matches);
    };

    // Listen for changes
    mediaQuery.addEventListener("change", handleMediaQueryChange);

    // Clean up event listener when component unmounts
    return () =>
      mediaQuery.removeEventListener("change", handleMediaQueryChange);
  }, [breakpoint]);

  return isMobile;
}
